package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	chatRequestTimeout = 25 * time.Second
	chatResponseLimit  = 1 << 20
	chatQuestionLimit  = 500
)

type chatRenderBlock struct {
	Type     string          `json:"type"`
	Markdown string          `json:"markdown,omitempty"`
	Text     string          `json:"text,omitempty"`
	Value    json.RawMessage `json:"value,omitempty"`
	Level    int             `json:"level,omitempty"`
	Style    string          `json:"style,omitempty"`
	Items    []string        `json:"items,omitempty"`
	Columns  []string        `json:"columns,omitempty"`
	Rows     [][]string      `json:"rows,omitempty"`
	Language *string         `json:"language,omitempty"`
	Code     string          `json:"code,omitempty"`
	Tone     string          `json:"tone,omitempty"`
	Title    *string         `json:"title,omitempty"`
}

type chatRenderDocument struct {
	Version int               `json:"version"`
	Blocks  []chatRenderBlock `json:"blocks"`
}

type chatResponse struct {
	Answer     string             `json:"answer"`
	Document   chatRenderDocument `json:"document"`
	Sources    []string           `json:"sources"`
	AnswerMode string             `json:"answerMode"`
	FollowUp   bool               `json:"followUp"`
	Format     string             `json:"format"`
}

type chatTurn struct {
	Question string
	Response chatResponse
}

type chatResultMsg struct {
	question string
	response chatResponse
	err      error
}

type chatClearResultMsg struct{ err error }

type chatHTTPError struct {
	Status     int
	Message    string
	RetryAfter string
}

func (e *chatHTTPError) Error() string { return e.Message }

type portfolioChatClient struct {
	baseURL   *url.URL
	token     string
	visitorID string
	http      *http.Client
	configErr error
}

func newPortfolioChatClientFromEnv(clientIP string) *portfolioChatClient {
	// Re-check local configuration for every SSH session. This also supports a
	// .env file created after the server process started.
	_ = loadEnvironment()
	return newPortfolioChatClient(
		os.Getenv("PORTFOLIO_RAG_URL"),
		os.Getenv("PORTFOLIO_RAG_SSH_TOKEN"),
		clientIP,
	)
}

func newPortfolioChatClient(baseURL, token, clientIP string) *portfolioChatClient {
	client := &portfolioChatClient{token: strings.TrimSpace(token)}
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		client.configErr = errors.New("PORTFOLIO_RAG_URL is not configured")
		return client
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && !(parsed.Scheme == "http" && parsed.Hostname() == "127.0.0.1")) {
		client.configErr = errors.New("PORTFOLIO_RAG_URL must use HTTPS")
		return client
	}
	if client.token == "" {
		client.configErr = errors.New("PORTFOLIO_RAG_SSH_TOKEN is not configured")
		return client
	}

	jar, _ := cookiejar.New(nil)
	client.baseURL = parsed
	client.visitorID = visitorIdentity(client.token, clientIP)
	client.http = &http.Client{
		Timeout: chatRequestTimeout,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.Scheme != parsed.Scheme || req.URL.Host != parsed.Host {
				return errors.New("refusing cross-origin chat redirect")
			}
			if len(via) >= 5 {
				return errors.New("too many chat redirects")
			}
			return nil
		},
	}
	return client
}

func visitorIdentity(token, clientIP string) string {
	mac := hmac.New(sha256.New, []byte(token))
	mac.Write([]byte(strings.TrimSpace(strings.ToLower(clientIP))))
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *portfolioChatClient) endpoint(path string) (string, error) {
	if c == nil {
		return "", errors.New("chat client is unavailable")
	}
	if c.configErr != nil {
		return "", c.configErr
	}
	target := *c.baseURL
	target.Path = strings.TrimRight(target.Path, "/") + path
	target.RawQuery = ""
	target.Fragment = ""
	return target.String(), nil
}

func (c *portfolioChatClient) request(method, path string, body io.Reader) (*http.Request, error) {
	endpoint, err := c.endpoint(path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-Portfolio-Visitor-ID", c.visitorID)
	return req, nil
}

func (c *portfolioChatClient) ask(question, mode string) (chatResponse, error) {
	payload, err := json.Marshal(struct {
		Question string `json:"question"`
		Mode     string `json:"mode"`
	}{Question: question, Mode: mode})
	if err != nil {
		return chatResponse{}, err
	}
	req, err := c.request(http.MethodPost, "/api/chat", bytes.NewReader(payload))
	if err != nil {
		return chatResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return chatResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, chatResponseLimit+1))
	if err != nil {
		return chatResponse{}, err
	}
	if len(body) > chatResponseLimit {
		return chatResponse{}, errors.New("assistant response was too large")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var problem struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &problem)
		if problem.Error == "" {
			problem.Error = http.StatusText(resp.StatusCode)
		}
		return chatResponse{}, &chatHTTPError{
			Status: resp.StatusCode, Message: sanitizeTerminal(problem.Error),
			RetryAfter: resp.Header.Get("Retry-After"),
		}
	}

	var result chatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return chatResponse{}, errors.New("assistant returned an invalid response")
	}
	if result.Document.Version != 1 && strings.TrimSpace(result.Answer) == "" {
		return chatResponse{}, errors.New("assistant returned an empty response")
	}
	return result, nil
}

func (c *portfolioChatClient) clear() error {
	if c != nil {
		defer c.resetCookies()
	}
	req, err := c.request(http.MethodDelete, "/api/chat/session", nil)
	if err != nil {
		return err
	}
	resp, requestErr := c.http.Do(req)
	if resp != nil {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
	}
	if requestErr != nil {
		return requestErr
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("session clear returned %d", resp.StatusCode)
	}
	return nil
}

func (c *portfolioChatClient) resetCookies() {
	if c == nil || c.http == nil {
		return
	}
	jar, _ := cookiejar.New(nil)
	c.http.Jar = jar
}

func sendChatCmd(client *portfolioChatClient, question, mode string) tea.Cmd {
	return func() tea.Msg {
		response, err := client.ask(question, mode)
		return chatResultMsg{question: question, response: response, err: err}
	}
}

func clearChatCmd(client *portfolioChatClient) tea.Cmd {
	return func() tea.Msg { return chatClearResultMsg{err: client.clear()} }
}

func chatErrorMessage(err error) string {
	var apiErr *chatHTTPError
	if errors.As(err, &apiErr) {
		switch apiErr.Status {
		case http.StatusUnauthorized, http.StatusForbidden:
			return "AI connection is not configured correctly (authentication failed)."
		case http.StatusTooManyRequests:
			if apiErr.RetryAfter != "" {
				return "Rate limit reached. Try again in " + sanitizeTerminal(apiErr.RetryAfter) + " seconds."
			}
			return "Rate limit reached. Please wait a moment and try again."
		default:
			return apiErr.Message
		}
	}
	if errors.Is(err, http.ErrHandlerTimeout) || strings.Contains(strings.ToLower(err.Error()), "timeout") || strings.Contains(strings.ToLower(err.Error()), "deadline exceeded") {
		return "The assistant timed out. Your question was restored; try again."
	}
	return "Could not reach the assistant. Your question was restored; try again."
}
