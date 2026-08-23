package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func testRenderer() *lipgloss.Renderer { return lipgloss.NewRenderer(io.Discard) }

func TestPortfolioChatConfigurationComesOnlyFromEnvironment(t *testing.T) {
	client := newPortfolioChatClient("", "", "203.0.113.8")
	if client.configErr == nil || client.configErr.Error() != "PORTFOLIO_RAG_URL is not configured" {
		t.Fatalf("missing URL error = %v", client.configErr)
	}

	t.Setenv("PORTFOLIO_RAG_URL", "https://rag.example.com")
	t.Setenv("PORTFOLIO_RAG_SSH_TOKEN", "environment-secret")
	client = newPortfolioChatClientFromEnv("203.0.113.8")
	if client.configErr != nil || client.baseURL.String() != "https://rag.example.com" || client.token != "environment-secret" {
		t.Fatalf("environment configuration was not used: %#v", client)
	}
}

func TestPortfolioChatClientSendsOnlyCurrentTurnAndKeepsItsCookie(t *testing.T) {
	var mu sync.Mutex
	var cookies []string
	var bodies []map[string]any
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" || r.Method != http.MethodPost {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("authorization = %q", got)
		}
		if got := r.Header.Get("X-Portfolio-Visitor-ID"); got != visitorIdentity("test-token", "203.0.113.8") {
			t.Errorf("visitor identity = %q", got)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		mu.Lock()
		bodies = append(bodies, body)
		cookies = append(cookies, r.Header.Get("Cookie"))
		mu.Unlock()
		http.SetCookie(w, &http.Cookie{Name: "__Host-portfolio_chat", Value: "session-one", Path: "/", Secure: true, HttpOnly: true})
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"answer":"ok","document":{"version":1,"blocks":[{"type":"text","text":"ok"}]},"sources":[],"format":"markdown"}`)
	}))
	defer server.Close()

	client := newPortfolioChatClient(server.URL, "test-token", "203.0.113.8")
	client.http.Transport = server.Client().Transport
	if _, err := client.ask("first question", "professional"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ask("follow up", "chaos"); err != nil {
		t.Fatal(err)
	}

	if len(bodies) != 2 || bodies[0]["question"] != "first question" || bodies[1]["question"] != "follow up" {
		t.Fatalf("unexpected bodies: %#v", bodies)
	}
	for _, body := range bodies {
		if len(body) != 2 || body["history"] != nil || body["state"] != nil {
			t.Fatalf("request leaked client state: %#v", body)
		}
	}
	if cookies[0] != "" || !strings.Contains(cookies[1], "__Host-portfolio_chat=session-one") {
		t.Fatalf("cookie continuity failed: %#v", cookies)
	}
}

func TestPortfolioChatClientsHaveIsolatedCookieJarsAndClear(t *testing.T) {
	var mu sync.Mutex
	seen := map[string][]string{}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		visitor := r.Header.Get("X-Portfolio-Visitor-ID")
		mu.Lock()
		seen[visitor] = append(seen[visitor], r.Header.Get("Cookie"))
		mu.Unlock()
		if r.Method == http.MethodDelete {
			w.Header().Set("Set-Cookie", "__Host-portfolio_chat=; Path=/; Max-Age=0; Secure; HttpOnly")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "__Host-portfolio_chat", Value: visitor[:8], Path: "/", Secure: true, HttpOnly: true})
		io.WriteString(w, `{"answer":"ok","document":{"version":1,"blocks":[{"type":"text","text":"ok"}]}}`)
	}))
	defer server.Close()

	first := newPortfolioChatClient(server.URL, "token", "198.51.100.1")
	second := newPortfolioChatClient(server.URL, "token", "198.51.100.2")
	first.http.Transport = server.Client().Transport
	second.http.Transport = server.Client().Transport
	if _, err := first.ask("one", "professional"); err != nil {
		t.Fatal(err)
	}
	if _, err := second.ask("two", "professional"); err != nil {
		t.Fatal(err)
	}
	if _, err := first.ask("follow up", "professional"); err != nil {
		t.Fatal(err)
	}
	if err := first.clear(); err != nil {
		t.Fatal(err)
	}
	if _, err := first.ask("fresh", "professional"); err != nil {
		t.Fatal(err)
	}

	firstSeen := seen[first.visitorID]
	secondSeen := seen[second.visitorID]
	if len(firstSeen) != 4 || firstSeen[0] != "" || firstSeen[1] == "" || firstSeen[3] != "" {
		t.Fatalf("first cookie lifecycle: %#v", firstSeen)
	}
	if len(secondSeen) != 1 || secondSeen[0] != "" {
		t.Fatalf("second client inherited a cookie: %#v", secondSeen)
	}
}

func TestPortfolioChatClientErrorsAreBoundedAndActionable(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "7")
		w.WriteHeader(http.StatusTooManyRequests)
		io.WriteString(w, `{"error":"slow down"}`)
	}))
	defer server.Close()
	client := newPortfolioChatClient(server.URL, "token", "127.0.0.1")
	client.http.Transport = server.Client().Transport
	_, err := client.ask("question", "professional")
	if got := chatErrorMessage(err); got != "Rate limit reached. Try again in 7 seconds." {
		t.Fatalf("message = %q", got)
	}

	authServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error":"Unauthorized SSH client"}`)
	}))
	defer authServer.Close()
	authClient := newPortfolioChatClient(authServer.URL, "token", "127.0.0.1")
	authClient.http.Transport = authServer.Client().Transport
	_, err = authClient.ask("question", "professional")
	if got := chatErrorMessage(err); !strings.Contains(got, "authentication failed") {
		t.Fatalf("auth message = %q", got)
	}

	invalidServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{not-json`)
	}))
	defer invalidServer.Close()
	invalidClient := newPortfolioChatClient(invalidServer.URL, "token", "127.0.0.1")
	invalidClient.http.Transport = invalidServer.Client().Transport
	_, err = invalidClient.ask("question", "professional")
	if err == nil || !strings.Contains(err.Error(), "invalid response") {
		t.Fatalf("invalid response error = %v", err)
	}

	timeoutServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
	}))
	defer timeoutServer.Close()
	timeoutClient := newPortfolioChatClient(timeoutServer.URL, "token", "127.0.0.1")
	timeoutClient.http.Transport = timeoutServer.Client().Transport
	timeoutClient.http.Timeout = 10 * time.Millisecond
	_, err = timeoutClient.ask("question", "professional")
	if !strings.Contains(chatErrorMessage(err), "timed out") {
		t.Fatalf("timeout message = %q", chatErrorMessage(err))
	}
}

func TestTerminalSanitizationAndTypedRendering(t *testing.T) {
	malicious := "safe\x1b]0;owned\x07 text\x1b[31m red\x1b[0m\x00"
	if got := sanitizeTerminal(malicious); got != "safe text red" {
		t.Fatalf("sanitized = %q", got)
	}
	title := "Note"
	language := "go"
	response := chatResponse{Document: chatRenderDocument{Version: 1, Blocks: []chatRenderBlock{
		{Type: "markdown", Markdown: "**Markdown** response"},
		{Type: "heading", Text: "Heading"},
		{Type: "text", Text: malicious},
		{Type: "list", Style: "ordered", Items: []string{"one", "two"}},
		{Type: "table", Columns: []string{"A", "B"}, Rows: [][]string{{"1", "2"}}},
		{Type: "code", Language: &language, Code: "fmt.Println(\"ok\")"},
		{Type: "json", Value: json.RawMessage(`{"ok":true}`)},
		{Type: "callout", Tone: "info", Title: &title, Text: "Details"},
		{Type: "future", Text: "must be ignored"},
	}}}
	m := Model{renderer: testRenderer()}
	joined := strings.Join(renderResponseLines(m, response, 28), "\n")
	for _, expected := range []string{"Markdown response", "Heading", "safe text red", "1. one", "A | B", "fmt.Println", `"ok"`, "Note", "Details"} {
		if !strings.Contains(joined, expected) {
			t.Errorf("missing %q in %q", expected, joined)
		}
	}
	if strings.Contains(joined, "owned") || strings.Contains(joined, "must be ignored") {
		t.Fatalf("unsafe or unknown content rendered: %q", joined)
	}
	for _, line := range renderResponseLines(m, chatResponse{Document: chatRenderDocument{Version: 1, Blocks: []chatRenderBlock{{Type: "text", Text: strings.Repeat("wide ", 20)}}}}, 16) {
		if lipgloss.Width(line) > 16 {
			t.Fatalf("line was not width-aware: width=%d line=%q", lipgloss.Width(line), line)
		}
	}

	fallback := renderResponseLines(m, chatResponse{Answer: "**legacy answer**"}, 20)
	if !strings.Contains(strings.Join(fallback, "\n"), "legacy answer") {
		t.Fatalf("legacy fallback missing: %#v", fallback)
	}
}

func TestAskAIInputFocusCommandsAndLimits(t *testing.T) {
	m := NewModel(testRenderer(), "203.0.113.9")
	m.width, m.height, m.mode = 100, 35, ViewNormal
	m.jumpTo(m.chatIndex())
	if !m.chatFocused {
		t.Fatal("assistant did not auto-focus")
	}
	section := m.sections[m.chatIndex()]
	if section.Label != "ask my assistant" || section.Icon != "🤖" {
		t.Fatalf("assistant navigation identity = label %q icon %q", section.Label, section.Icon)
	}

	longInput := []rune(strings.Repeat("界", chatQuestionLimit+10))
	updated, _ := m.handleChatKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: longInput})
	m = updated.(Model)
	if got := len([]rune(m.chatInput)); got != chatQuestionLimit {
		t.Fatalf("input runes = %d", got)
	}
	updated, _ = m.handleChatKey(tea.KeyMsg{Type: tea.KeyBackspace})
	m = updated.(Model)
	if got := len([]rune(m.chatInput)); got != chatQuestionLimit-1 {
		t.Fatalf("backspace runes = %d", got)
	}

	m.chatInput = "/help"
	updated, cmd := m.submitChatInput()
	m = updated.(Model)
	if cmd != nil || !strings.Contains(m.chatError, "Available command: /clear") {
		t.Fatalf("removed help command was still accepted: error=%q cmd=%v", m.chatError, cmd)
	}
	m.chatInput = "/mode chaos"
	updated, cmd = m.submitChatInput()
	m = updated.(Model)
	if cmd != nil || !strings.Contains(m.chatError, "Available command: /clear") {
		t.Fatalf("removed mode command was still accepted: error=%q cmd=%v", m.chatError, cmd)
	}
	m.chatPending = true
	m.chatInput = "blocked duplicate"
	updated, cmd = m.submitChatInput()
	if cmd != nil || !updated.(Model).chatPending {
		t.Fatal("duplicate submission was not blocked")
	}

	m.chatPending = false
	updated, _ = m.handleChatKey(tea.KeyMsg{Type: tea.KeyEscape})
	m = updated.(Model)
	if m.chatFocused {
		t.Fatal("escape did not restore portfolio navigation")
	}
	m.jumpTo(8)
	if m.sections[m.selected].Key != "contact" {
		t.Fatalf("section 9 = %q", m.sections[m.selected].Key)
	}
}

func TestAssistantTranscriptScrollsWhilePromptIsFocused(t *testing.T) {
	m := NewModel(testRenderer(), "203.0.113.9")
	m.width, m.height, m.mode = 100, 18, ViewNormal
	m.jumpTo(m.chatIndex())
	for index := 0; index < 10; index++ {
		m.chatTurns = append(m.chatTurns, chatTurn{
			Question: "Tell me about another project.",
			Response: chatResponse{Answer: "This is a detailed assistant response that occupies transcript space."},
		})
	}
	m.scrollToBottom()
	bottom := m.scrollPos
	if bottom == 0 {
		t.Fatal("test transcript did not overflow the viewport")
	}

	updated, _ := m.handleChatKey(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.scrollPos != bottom-1 {
		t.Fatalf("focused up key did not scroll: got %d, want %d", m.scrollPos, bottom-1)
	}
	updated, _ = m.handleChatKey(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.scrollPos != bottom {
		t.Fatalf("focused down key did not scroll: got %d, want %d", m.scrollPos, bottom)
	}

	// Old bottom-following used a huge sentinel that made wheel-up appear stuck.
	m.scrollPos = 1 << 30
	m.scrollUp()
	if m.scrollPos != bottom-1 {
		t.Fatalf("out-of-range scroll was not clamped before moving up: got %d, want %d", m.scrollPos, bottom-1)
	}
}

func TestAskAIViewHidesResponseSources(t *testing.T) {
	m := NewModel(testRenderer(), "203.0.113.9")
	m.chatTurns = []chatTurn{{
		Question: "What did you build?",
		Response: chatResponse{
			Answer:  "A portfolio.",
			Sources: []string{"private-source-label"},
		},
	}}

	view := strings.Join(buildChat(m, 80), "\n")
	if strings.Contains(view, "Sources:") || strings.Contains(view, "private-source-label") {
		t.Fatalf("chat view exposed response sources: %q", view)
	}
}
