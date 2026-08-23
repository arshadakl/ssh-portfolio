package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewMode int

const (
	ViewLoading ViewMode = iota
	ViewNormal
)

type tickMsg time.Time

// contactResultMsg is posted back after the async HTTP POST completes.
type contactResultMsg struct {
	ok  bool
	err error
}

type Model struct {
	renderer  *lipgloss.Renderer
	width     int
	height    int
	selected  int
	scrollPos int
	mode      ViewMode
	animFrame int
	sections  []Section

	// personalization
	clientIP     string
	sessionStart time.Time

	// command bar (:)
	cmdMode  bool
	cmdInput string
	cmdMsg   string
	blinkOn  bool

	// contribute section intent flow
	intentCursor int
	intentActive int // -1 = collapsed

	// contact form (in-section, not command-bar wizard)
	contactFormActive bool   // true when form is open and editable
	contactField      int    // 0=name, 1=email, 2=subject, 3=message
	contactSubmitting bool   // true while HTTP POST is in flight
	contactResult     string // plain-text status shown after submit
	contactResultOK   bool   // true when contactResult is a success message
	contactVisited    bool   // true once the visitor has reached the contact section
	contactName       string
	contactEmail      string
	contactSubject    string
	contactMessage    string

	// Assistant section (display transcript is local to this SSH connection)
	chatClient          *portfolioChatClient
	chatFocused         bool
	chatInput           string
	chatPending         bool
	chatClearPending    bool
	chatPendingQuestion string
	chatTurns           []chatTurn
	chatNotice          string
	chatError           string
	chatSpinner         int
}

func NewModel(r *lipgloss.Renderer, clientIP string) Model {
	return Model{
		renderer:     r,
		mode:         ViewLoading,
		sections:     buildSections(r),
		clientIP:     clientIP,
		sessionStart: time.Now(),
		intentActive: -1,
		chatClient:   newPortfolioChatClientFromEnv(clientIP),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) isNarrow() bool { return m.width < 70 }

func (m Model) contentInnerWidth() int {
	if m.isNarrow() {
		if width := m.width - 4; width > 0 {
			return width
		}
		return 1
	}
	width := m.width - m.sidebarWidth() - 5
	if width < 1 {
		return 1
	}
	return width
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if m.mode == ViewLoading {
			m.animFrame++
			fullText := "Accessing Arshad's system..."
			if m.animFrame > len(fullText)+5 {
				m.mode = ViewNormal
			}
			return m, tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
		// 1s heartbeat drives the footer session timer and cursor blink
		m.blinkOn = !m.blinkOn
		m.chatSpinner = (m.chatSpinner + 1) % 4
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case contactResultMsg:
		m.contactSubmitting = false
		m.contactFormActive = false
		m.contactField = 0
		m.resetContact()
		if msg.ok {
			m.contactResult = "Message sent. I'll respond usually within 24h."
			m.contactResultOK = true
		} else {
			m.contactResult = "Failed to send. Try emailing directly: arshadayanikkal@gmail.com"
			m.contactResultOK = false
		}
		return m, nil

	case chatResultMsg:
		m.chatPending = false
		m.chatPendingQuestion = ""
		if msg.err != nil {
			m.chatInput = msg.question
			m.chatError = chatErrorMessage(msg.err)
			m.chatNotice = ""
		} else {
			m.chatTurns = append(m.chatTurns, chatTurn{Question: msg.question, Response: msg.response})
			m.chatInput = ""
			m.chatError = ""
			m.chatNotice = ""
		}
		m.scrollToBottom()
		return m, nil

	case chatClearResultMsg:
		m.chatClearPending = false
		if msg.err != nil {
			m.chatError = "Visible chat cleared, but the server session could not be cleared; it expires automatically."
		} else {
			m.chatNotice = "Conversation cleared."
			m.chatError = ""
		}
		m.scrollToBottom()
		return m, nil

	case tea.MouseMsg:
		if m.mode != ViewLoading && !m.cmdMode && !m.contactFormActive && msg.Action == tea.MouseActionPress {
			switch msg.Button {
			case tea.MouseButtonWheelUp:
				m.scrollUp()
			case tea.MouseButtonWheelDown:
				m.scrollDown()
			}
		}

	case tea.KeyMsg:
		if m.mode == ViewLoading {
			return m, nil
		}

		// command bar captures all input while open
		if m.cmdMode {
			return m.handleCommandKey(msg)
		}

		// The assistant owns normal typing while its prompt is focused.
		if m.selected == m.chatIndex() {
			if m.chatFocused {
				return m.handleChatKey(msg)
			}
			if msg.Type == tea.KeyEnter {
				m.chatFocused = true
				m.chatError = ""
				m.scrollToBottom()
				return m, nil
			}
		}

		// contact section: in-section form
		if m.selected == m.contactIndex() {
			if m.contactFormActive {
				return m.handleContactFormKey(msg)
			}
			if msg.Type == tea.KeyEnter {
				m.contactFormActive = true
				m.contactField = 0
				m.contactResult = ""
				m.scrollPos = 0
				return m, nil
			}
		}

		// contribute section: intent flow owns left/right/tab/enter/esc
		if m.selected == m.contributeIndex() {
			switch msg.Type {
			case tea.KeyLeft:
				m.intentPrev()
				return m, nil
			case tea.KeyRight, tea.KeyTab:
				m.intentNext()
				return m, nil
			case tea.KeyEnter:
				m.intentToggle()
				return m, nil
			case tea.KeyEscape:
				m.intentActive = -1
				return m, nil
			}
		}

		// number keys 1-9 jump to sections
		if s := msg.String(); len(s) == 1 && s[0] >= '1' && s[0] <= '9' {
			m.jumpTo(int(s[0] - '1'))
			return m, nil
		}

		// Arrow keys navigate sections (type-based for SSH compat)
		switch msg.Type {
		case tea.KeyUp:
			m.navigatePrev()
			return m, nil
		case tea.KeyDown:
			m.navigateNext()
			return m, nil
		case tea.KeyLeft:
			if m.isNarrow() {
				m.navigatePrev()
			}
			return m, nil
		case tea.KeyRight:
			if m.isNarrow() {
				m.navigateNext()
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case ":":
			m.cmdMode = true
			m.cmdInput = ""
			m.cmdMsg = ""

		case "j", "down":
			m.navigateNext()
		case "k", "up":
			m.navigatePrev()
		case "right":
			if m.isNarrow() {
				m.navigateNext()
			}
		case "left":
			if m.isNarrow() {
				m.navigatePrev()
			}

		case "w":
			m.scrollUp()
		case "s":
			m.scrollDown()
		case "pgup", "ctrl+u":
			for i := 0; i < 10; i++ {
				m.scrollUp()
			}
		case "pgdown", "ctrl+d":
			for i := 0; i < 10; i++ {
				m.scrollDown()
			}

		case "enter":
			m.scrollPos = 0
		}
	}
	return m, nil
}

func (m Model) handleChatKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		m.chatFocused = false
		return m, nil
	case tea.KeyUp:
		m.scrollUp()
		return m, nil
	case tea.KeyDown:
		m.scrollDown()
		return m, nil
	case tea.KeyPgUp:
		for i := 0; i < 10; i++ {
			m.scrollUp()
		}
		return m, nil
	case tea.KeyPgDown:
		for i := 0; i < 10; i++ {
			m.scrollDown()
		}
		return m, nil
	case tea.KeyBackspace:
		runes := []rune(m.chatInput)
		if len(runes) > 0 && !m.chatPending && !m.chatClearPending {
			m.chatInput = string(runes[:len(runes)-1])
		}
		return m, nil
	case tea.KeyEnter:
		return m.submitChatInput()
	}

	switch msg.String() {
	case "ctrl+c":
		if !m.chatPending && !m.chatClearPending {
			if m.chatInput == "" {
				m.chatFocused = false
			} else {
				m.chatInput = ""
			}
		}
		return m, nil
	case "ctrl+u":
		for i := 0; i < 10; i++ {
			m.scrollUp()
		}
		return m, nil
	case "ctrl+d":
		for i := 0; i < 10; i++ {
			m.scrollDown()
		}
		return m, nil
	}

	if (msg.Type == tea.KeySpace || msg.Type == tea.KeyRunes) && !m.chatPending && !m.chatClearPending {
		remaining := chatQuestionLimit - len([]rune(m.chatInput))
		if remaining > 0 {
			runes := msg.Runes
			if len(runes) > remaining {
				runes = runes[:remaining]
			}
			m.chatInput += string(runes)
			m.chatError = ""
			m.chatNotice = ""
			m.scrollToBottom()
		}
	}
	return m, nil
}

func (m Model) submitChatInput() (tea.Model, tea.Cmd) {
	if m.chatPending || m.chatClearPending {
		return m, nil
	}
	question := strings.TrimSpace(m.chatInput)
	if question == "" {
		return m, nil
	}
	if strings.HasPrefix(question, "/") {
		return m.handleChatCommand(question)
	}
	if m.chatClient == nil || m.chatClient.configErr != nil {
		reason := "chat client is unavailable"
		if m.chatClient != nil && m.chatClient.configErr != nil {
			reason = m.chatClient.configErr.Error()
		}
		m.chatError = "AI is unavailable: " + reason
		return m, nil
	}
	m.chatPending = true
	m.chatPendingQuestion = question
	m.chatInput = ""
	m.chatError = ""
	m.chatNotice = ""
	m.scrollToBottom()
	return m, sendChatCmd(m.chatClient, question, "professional")
}

func (m Model) handleChatCommand(input string) (tea.Model, tea.Cmd) {
	parts := strings.Fields(strings.ToLower(input))
	m.chatInput = ""
	m.chatError = ""
	switch parts[0] {
	case "/clear":
		m.chatTurns = nil
		m.chatPendingQuestion = ""
		m.chatNotice = ""
		m.chatClearPending = true
		m.scrollPos = 0
		return m, clearChatCmd(m.chatClient)
	default:
		m.chatError = "Unknown AI command. Available command: /clear."
	}
	m.scrollToBottom()
	return m, nil
}

// handleCommandKey processes keystrokes while the command bar is open.

// -- contact wizard -----------------------------------------------------------

// handleContactFormKey processes keystrokes while the in-section contact form
// is open. Tab/Enter advance fields; Enter on the message field submits.
func (m Model) handleContactFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.contactSubmitting {
		return m, nil // HTTP in flight -- block input
	}

	switch msg.String() {
	case "esc":
		m.contactFormActive = false
		m.contactField = 0
		m.contactResult = ""
		m.resetContact()
		return m, nil

	case "enter", "tab":
		if m.contactField < 3 {
			m.contactField++
			return m, nil
		}
		// last field (message) + enter => submit
		if strings.TrimSpace(m.contactName) == "" || strings.TrimSpace(m.contactMessage) == "" {
			m.contactResult = "name and message are required"
			m.contactResultOK = false
			return m, nil
		}
		m.contactSubmitting = true
		return m, sendContactCmd(
			m.contactName, m.contactEmail,
			m.contactSubject, m.contactMessage,
		)

	case "shift+tab", "up":
		if m.contactField > 0 {
			m.contactField--
		}
		return m, nil

	case "down":
		if m.contactField < 3 {
			m.contactField++
		}
		return m, nil
	}

	field := m.contactFieldPtr()
	if m.contactResult != "" && !m.contactResultOK {
		m.contactResult = "" // clear validation hint once the user edits
	}
	switch {
	case msg.Type == tea.KeyBackspace:
		if len(*field) > 0 {
			*field = (*field)[:len(*field)-1]
		}
	case msg.Type == tea.KeySpace, msg.Type == tea.KeyRunes:
		if len(msg.Runes) > 0 {
			*field += string(msg.Runes)
		}
	}
	return m, nil
}

// contactFieldPtr returns a pointer to the field currently being edited.
func (m *Model) contactFieldPtr() *string {
	switch m.contactField {
	case 0:
		return &m.contactName
	case 1:
		return &m.contactEmail
	case 2:
		return &m.contactSubject
	default:
		return &m.contactMessage
	}
}

func (m *Model) resetContact() {
	m.contactName = ""
	m.contactEmail = ""
	m.contactSubject = ""
	m.contactMessage = ""
}

// sendContactCmd fires an async HTTP POST to the contact API.
func sendContactCmd(name, email, subject, message string) tea.Cmd {
	const endpoint = "https://minecraft.arshadakl.in/api/contact"

	return func() tea.Msg {
		// Tag the subject so emails are traceable back to this portfolio.
		if strings.TrimSpace(subject) != "" {
			subject = "[from ssh portfolio] " + subject
		} else {
			subject = "[from ssh portfolio]"
		}
		body := map[string]string{
			"name": name, "email": email,
			"subject": subject, "message": message,
		}
		payload, err := json.Marshal(body)
		if err != nil {
			return contactResultMsg{err: err}
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Post(endpoint, "application/json", bytes.NewReader(payload))
		if err != nil {
			return contactResultMsg{err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return contactResultMsg{err: fmt.Errorf("server returned %d", resp.StatusCode)}
		}
		return contactResultMsg{ok: true}
	}
}

// handleCommandKey processes keystrokes while the command bar is open.
func (m Model) handleCommandKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		m.cmdMode = false
		m.cmdInput = ""
		return m, nil

	case tea.KeyEnter:
		res := execCommand(m.renderer, m.cmdInput)
		m.cmdMode = false
		m.cmdInput = ""
		m.cmdMsg = res.message
		if res.jump >= 0 {
			m.jumpTo(res.jump)
		}
		if res.activate >= 0 {
			m.jumpTo(m.contributeIndex())
			m.intentCursor = res.activate
			m.intentActive = res.activate
			m.scrollToBottom()
		}
		if res.quit {
			return m, tea.Quit
		}
		if res.startContact {
			m.jumpTo(m.contactIndex())
			m.contactFormActive = true
			m.contactField = 0
			m.contactResult = ""
			m.scrollPos = 0
		}
		return m, nil

	case tea.KeyBackspace:
		if len(m.cmdInput) > 0 {
			m.cmdInput = m.cmdInput[:len(m.cmdInput)-1]
		}
		return m, nil

	case tea.KeySpace, tea.KeyRunes:
		m.cmdInput += string(msg.Runes)
		return m, nil
	}
	return m, nil
}

func (m *Model) navigateNext() {
	if m.selected < len(m.sections)-1 {
		m.selected++
		m.scrollPos = 0
		m.autoOpenContact()
		m.autoFocusChat()
	}
}

func (m *Model) navigatePrev() {
	if m.selected > 0 {
		m.selected--
		m.scrollPos = 0
		m.autoOpenContact()
		m.autoFocusChat()
	}
}

func (m *Model) jumpTo(idx int) {
	if idx >= 0 && idx < len(m.sections) {
		m.selected = idx
		m.scrollPos = 0
		m.mode = ViewNormal
		m.autoOpenContact()
		m.autoFocusChat()
	}
}

func (m *Model) autoFocusChat() {
	m.chatFocused = m.selected == m.chatIndex()
}

// autoOpenContact opens the message form the first time the visitor arrives at
// the contact section so the form is immediately discoverable.
func (m *Model) autoOpenContact() {
	if m.selected == m.contactIndex() && !m.contactVisited {
		m.contactVisited = true
		m.contactFormActive = true
		m.contactField = 0
		m.contactResult = ""
	}
}

func (m *Model) scrollUp() {
	m.clampScroll()
	if m.scrollPos > 0 {
		m.scrollPos--
	}
}

func (m *Model) scrollDown() {
	m.clampScroll()
	if m.scrollPos < m.maxScroll() {
		m.scrollPos++
	}
}

func (m *Model) scrollToBottom() {
	m.scrollPos = m.maxScroll()
}

func (m *Model) clampScroll() {
	maximum := m.maxScroll()
	if m.scrollPos > maximum {
		m.scrollPos = maximum
	}
	if m.scrollPos < 0 {
		m.scrollPos = 0
	}
}

func (m Model) maxScroll() int {
	available := m.height - 6
	if m.isNarrow() {
		available = m.height - 7
	}
	if m.cmdBarVisible() {
		available--
	}
	if available < 1 {
		available = 1
	}
	maximum := len(m.visibleLines()) - available
	if maximum < 0 {
		return 0
	}
	return maximum
}

func (m *Model) intentNext() {
	m.intentCursor = (m.intentCursor + 1) % len(intents)
}

func (m *Model) intentPrev() {
	m.intentCursor = (m.intentCursor + len(intents) - 1) % len(intents)
}

func (m *Model) intentToggle() {
	if m.intentActive == m.intentCursor {
		m.intentActive = -1
	} else {
		m.intentActive = m.intentCursor
		m.scrollToBottom()
	}
}

func (m Model) contributeIndex() int {
	for i, sec := range m.sections {
		if sec.Key == "contribute" {
			return i
		}
	}
	return -1
}

func (m Model) chatIndex() int {
	for i, sec := range m.sections {
		if sec.Key == "ask-ai" {
			return i
		}
	}
	return -1
}

func (m Model) contactIndex() int {
	for i, sec := range m.sections {
		if sec.Key == "contact" {
			return i
		}
	}
	return -1
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	if m.mode == ViewLoading {
		return m.renderSplash()
	}
	if m.isNarrow() {
		return m.renderNarrowLayout()
	}
	return m.renderWideLayout()
}

func (m Model) visibleLines() []string {
	sec := m.sections[m.selected]
	if sec.Key == "ask-ai" {
		return buildChat(m, m.contentInnerWidth())
	}
	if sec.Key == "contribute" {
		return buildContribute(m.renderer, m.intentCursor, m.intentActive)
	}
	if sec.Key == "contact" {
		if m.contactFormActive {
			return buildContactForm(m)
		}
		if m.contactResult != "" {
			var status string
			if m.contactResultOK {
				status = styleGreen(m.renderer).Bold(true).Render("✓ " + m.contactResult)
			} else {
				status = styleDim(m.renderer).Render("✗ " + m.contactResult)
			}
			return append([]string{status, ""}, sec.Lines...)
		}
	}
	return sec.Lines
}
