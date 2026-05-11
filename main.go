package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	menuWhoami     = "whoami"
	menuExperience = "experience"
	menuSecurity   = "security"
	menuProjects   = "projects"
	menuSkills     = "skills"
	menuContact    = "contact"
)

var (
	menuItems = []string{
		menuWhoami,
		menuExperience,
		menuSecurity,
		menuProjects,
		menuSkills,
		menuContact,
	}

	menuIcons = map[string]string{
		menuWhoami:     "◈",
		menuExperience: "▸",
		menuSecurity:   "⚡",
		menuProjects:   "⬡",
		menuSkills:     "⬢",
		menuContact:    "✉",
	}

	darkBg       = lipgloss.Color("0")
	orange       = lipgloss.Color("208")
	lightGray    = lipgloss.Color("250")
	darkGray     = lipgloss.Color("240")
	subtleAccent = lipgloss.Color("66")
)

type tickMsg time.Time

type Model struct {
	width        int
	height       int
	selected     int
	contentMap   map[string]string
	loading      bool
	animFrame    int
	animComplete bool
	scrollPos    int
}

func NewModel() Model {
	return Model{
		selected:   0,
		loading:    true,
		animFrame:  0,
		contentMap: buildContentMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(60*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		if m.loading {
			m.animFrame++
			fullText := "Accessing Arshad's system..."
			if m.animFrame > len(fullText)+5 {
				m.loading = false
				m.animComplete = true
			}
			return m, tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
	case tea.MouseMsg:
		if !m.loading {
			if msg.Action == tea.MouseActionPress {
				switch msg.Button {
				case tea.MouseButtonWheelUp:
					if m.scrollPos > 0 {
						m.scrollPos--
					}
				case tea.MouseButtonWheelDown:
					m.scrollPos++
				}
			}
		}
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left":
			if m.selected > 0 {
				m.selected--
				m.scrollPos = 0
			}
		case "right":
			if m.selected < len(menuItems)-1 {
				m.selected++
				m.scrollPos = 0
			}
		case "up":
			if m.scrollPos > 0 {
				m.scrollPos--
			}
		case "down":
			m.scrollPos++
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	if m.loading {
		return m.renderSplash()
	}

	contentWidth := m.width - 8
	if contentWidth < 40 {
		contentWidth = m.width - 18
	}
	contentBoxHeight := m.getContentBoxHeight()

	currentMenu := menuItems[m.selected]
	fullContent := m.contentMap[currentMenu]
	wrappedLines := strings.Split(wrapText(fullContent, contentWidth-6), "\n")
	totalLines := len(wrappedLines)

	navbar := m.renderNavbar(contentWidth)
	content := m.renderContent(contentWidth, contentBoxHeight)
	footer := m.renderFooter(contentWidth, totalLines)

	bodyContent := lipgloss.JoinVertical(lipgloss.Left, navbar, content, footer)

	horizontalPadding := (m.width - contentWidth) / 2
	if horizontalPadding < 1 {
		horizontalPadding = 1
	}

	totalHeight := contentBoxHeight + 4
	verticalPadding := (m.height - totalHeight) / 2
	if verticalPadding < 0 {
		verticalPadding = 0
	}

	paddedBox := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		PaddingLeft(horizontalPadding).
		PaddingTop(verticalPadding).
		Background(darkBg).
		Render(bodyContent)

	return paddedBox
}

func (m Model) getContentBoxHeight() int {
	maxHeight := m.height - 8
	if maxHeight > 28 {
		maxHeight = 28
	}
	if maxHeight < 6 {
		maxHeight = 6
	}
	return maxHeight
}

func (m Model) renderSplash() string {
	splashStyle := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Width(m.width).
		Height(m.height).
		Background(darkBg).
		Foreground(orange)

	asciiArt := `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣤⠶⠶⠶⠶⢦⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⡾⠛⠁⠀⠀⠀⠀⠀⠀⠈⠙⢷⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⠏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⢷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡾⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡾⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣷⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⠀⣀⣀⣀⣀⣀⣀⠀⠀⠀⠀⠀⠀⠀⠸⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⠀⠀⠀⠀⣠⡴⠞⠛⠉⠉⣩⣍⠉⠉⠛⠳⢦⣄⠀⠀⠀⠀⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⡀⠀⣴⡿⣧⣀⠀⢀⣠⡴⠋⠙⢷⣄⡀⠀⣀⣼⢿⣦⠀⠀⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣧⡾⠋⣷⠈⠉⠉⠉⠉⠀⠀⠀⠀⠉⠉⠋⠉⠁⣼⠙⢷⣼⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣇⠀⢻⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⡟⠀⣸⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣹⣆⠀⢻⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⡟⠀⣰⣏⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⠞⠋⠁⠙⢷⣄⠙⢷⣀⠀⠀⠀⠀⠀⠀⢀⡴⠋⢀⡾⠋⠈⠙⠻⢦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⡾⠋⠀⠀⠀⠀⠀⠀⠹⢦⡀⠙⠳⠶⢤⡤⠶⠞⠋⢀⡴⠟⠀⠀⠀⠀⠀⠀⠙⠻⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⣼⠋⠀⠀⢀⣤⣤⣤⣤⣤⣤⣤⣿⣦⣤⣤⣤⣤⣤⣤⣴⣿⣤⣤⣤⣤⣤⣤⣤⡀⠀⠀⠙⣧⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣸⠏⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⠀⢠⣴⠞⠛⠛⠻⢦⡄⠀⠀⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠸⣇⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢠⡟⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⠀⣿⣿⢶⣄⣠⡶⣦⣿⠀⠀⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀⢻⡄⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣾⠁⠀⠀⠀⠀⠘⣇⠀⠀⠀⠀⠀⠀⠀⢻⣿⠶⠟⠻⠶⢿⡿⠀⠀⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀⠈⣿⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢰⡏⠀⠀⠀⠀⠀⠀⣿⠀⠀⠀⠀⠀⠀⢾⣄⣹⣦⣀⣀⣴⢟⣠⡶⠀⠀⠀⠀⠀⠀⣼⠀⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⣿⠀⠀⠀⠀⠀⠀⠀⠈⠛⠿⣭⣭⡿⠛⠁⠀⠀⠀⠀⠀⠀⠀⣿⠀⠀⠀⠀⠀⠀⠘⣧⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⢿⡀⠀⠀⠀⠀⠀⠀⣀⡴⠞⠋⠙⠳⢦⣀⠀⠀⠀⠀⠀⠀⠀⣿⠀⠀⠀⠀⠀⠀⢰⡏⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠈⢿⣄⣀⠀⠀⢀⣤⣼⣧⣤⣤⣤⣤⣤⣿⣭⣤⣤⣤⣤⣤⣤⣭⣿⣤⣤⣤⣤⣤⣼⣿⣤⣄⠀⠀⣀⣠⡾⠁⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠈⠉⠛⠛⠻⢧⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠼⠟⠛⠛⠉⠁⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⣷⣶⣶⣶⣶⣶⣿⣷⣶⣿⣿⣾⣿⣶⣶⣿⣿⣷⣿⣿⣿⣿⣿⣿⣾⣿⣿⣿⣿⣷⣷⣿⣷⣶⣶⣶⣶⣶⣶⣶⣶⣶⣶⣶⣶⣶⣶⣶⣶⣶
⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣷⣶⣿⣿`

	fullText := "Accessing Arshad's system..."
	visibleChars := m.animFrame
	if visibleChars > len(fullText) {
		visibleChars = len(fullText)
	}

	displayText := fullText[:visibleChars]
	if visibleChars < len(fullText) {
		displayText += "▊"
	}

	content := asciiArt + "\n\n" + displayText

	return splashStyle.Render(content)
}

func (m Model) renderNavbar(navbarWidth int) string {
	navbarContent := strings.Builder{}
	navbarContent.WriteString("Arshad | Portfolio")

	spacing := navbarWidth - 20
	if spacing > 0 {
		navbarContent.WriteString(strings.Repeat(" ", spacing))
	}

	barStyle := lipgloss.NewStyle().
		Width(navbarWidth).
		Padding(0, 1).
		Background(darkBg).
		Foreground(lightGray)

	navbarContent.WriteString("\n")

	menuStr := strings.Builder{}
	for i, item := range menuItems {
		var itemStyle lipgloss.Style
		if i == m.selected {
			itemStyle = lipgloss.NewStyle().
				Foreground(darkBg).
				Background(orange).
				Bold(true).
				Padding(0, 1)
		} else {
			itemStyle = lipgloss.NewStyle().
				Foreground(lightGray).
				Padding(0, 1)
		}

		menuStr.WriteString(itemStyle.Render(menuIcons[item] + " " + item))
		if i < len(menuItems)-1 {
			menuStr.WriteString(" ")
		}
	}

	navbarContent.WriteString(menuStr.String())
	return barStyle.Render(navbarContent.String())
}

func (m Model) renderContent(contentWidth, contentHeight int) string {
	currentMenu := menuItems[m.selected]
	fullContent := m.contentMap[currentMenu]

	headerStyle := lipgloss.NewStyle().Foreground(orange).Bold(true)
	sepStyle := lipgloss.NewStyle().Foreground(darkGray)
	sep := strings.Repeat("─", contentWidth-6)
	header := headerStyle.Render(currentMenu) + "\n" + sepStyle.Render(sep)

	wrappedContent := wrapText(fullContent, contentWidth-6)
	lines := strings.Split(wrappedContent, "\n")

	// 2 header lines + 1 blank separator + 2 padding lines = 5 overhead
	availableLines := contentHeight - 5
	if availableLines < 1 {
		availableLines = 1
	}

	startLine := m.scrollPos
	endLine := startLine + availableLines
	if endLine > len(lines) {
		endLine = len(lines)
	}
	if startLine >= len(lines) {
		startLine = len(lines) - availableLines
		if startLine < 0 {
			startLine = 0
		}
		endLine = len(lines)
	}

	displayLines := lines[startLine:endLine]
	displayContent := header + "\n" + strings.Join(displayLines, "\n")

	contentStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		PaddingLeft(2).
		PaddingRight(2).
		PaddingTop(1).
		PaddingBottom(1).
		Background(darkBg).
		Foreground(lightGray).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(darkGray)

	return contentStyle.Render(displayContent)
}

func (m Model) renderFooter(footerWidth int, totalLines int) string {
	scrollDisplay := fmt.Sprintf("[line %d/%d]", m.scrollPos+1, totalLines)
	footerText := fmt.Sprintf("q quit  ←→ sections  ↑↓ scroll  %s", scrollDisplay)
	padding := footerWidth - len(footerText) - 2
	if padding < 0 {
		padding = 0
	}

	footer := fmt.Sprintf(" %s%s", footerText, strings.Repeat(" ", padding))

	footerStyle := lipgloss.NewStyle().
		Width(footerWidth).
		Background(darkBg).
		Foreground(darkGray)

	return footerStyle.Render(footer)
}

func wrapText(text string, maxWidth int) string {
	if maxWidth < 1 {
		maxWidth = 1
	}

	lines := strings.Split(text, "\n")
	var wrapped []string

	for _, line := range lines {
		if line == "" {
			wrapped = append(wrapped, "")
			continue
		}

		// Skip word-wrapping lines that contain ANSI escape codes (styled content)
		if strings.Contains(line, "\033[") {
			wrapped = append(wrapped, line)
			continue
		}

		words := strings.Fields(line)
		var currentLine strings.Builder
		currentLen := 0

		for _, word := range words {
			wordLen := len(word)

			if currentLen+wordLen+1 > maxWidth && currentLen > 0 {
				wrapped = append(wrapped, currentLine.String())
				currentLine.Reset()
				currentLen = 0
			}

			if currentLen > 0 {
				currentLine.WriteString(" ")
				currentLen++
			}

			currentLine.WriteString(word)
			currentLen += wordLen
		}

		if currentLine.Len() > 0 {
			wrapped = append(wrapped, currentLine.String())
		}
	}

	return strings.Join(wrapped, "\n")
}

func buildContentMap() map[string]string {
	sub := lipgloss.NewStyle().Foreground(subtleAccent)

	return map[string]string{
		menuWhoami: `Arshad. Full-stack engineer + security researcher.
2+ years in production engineering.

Shipped two products at ELT Global — student-facing LMS and admin
operations portal serving 100,000+ monthly users, 10,000+ DAU,
16M+ API requests/month.

Stack: NestJS backends + Next.js frontends in a monorepo. TypeScript
throughout. Focus on performance, correctness, and maintainability.

Security background: CERT-In Hall of Fame (Government of India).
Active on HackerOne and Bugcrowd. Bug bounty focus on web app
vulnerabilities, API security, and business logic flaws.

Currently exploring new opportunities. Open to full-time roles in
product engineering or security-adjacent engineering.`,

		menuExperience: sub.Render("ELT Global Pvt Ltd — Software Engineer") + `
Bangalore, India | Aug 2024 – Present

Full-stack on an EdTech platform at 10k+ DAU — two products: a
student-facing LMS and an admin operations portal. Monorepo:
NestJS backends + Next.js frontends.

- SDUI-driven interfaces — backend controls screen rendering
  without client redeploys
- Built automation tooling (Google Apps Script workflows +
  Docker crash monitor with Slack alerting) — eliminated
  3-4 hrs/day of repetitive overhead
- Type-safe OpenAPI Swagger + Codegen pipeline — auto-generates
  all API schemas and TypeScript types from a single command,
  reducing frontend-backend integration overhead by ~60%
- Security hardening across auth flows, token rotation, RBAC,
  input validation, and request signing
- Improved page performance 35-40% via code-splitting, TanStack
  Query cache-first patterns, request deduplication, API batching
- Cut key analytics API payload by ~87% through client-side
  derivation

` + sub.Render("Brototype — Full-Stack Engineering Intern") + `
Calicut | 2023 – 2024

- Built freelancer marketplace with Socket.IO real-time collab,
  Stripe payments, AWS EC2 deployment, ranking algorithm
- Developed full e-commerce platform: coupon engine, Razorpay,
  session-based auth, admin panel — deployed end-to-end`,

		menuSecurity: sub.Render("CERT-In Hall of Fame — Government of India") + `

Discovered and responsibly disclosed a critical vulnerability in a
major Kerala university's official website. Exposed unauthenticated
access to backend systems affecting 200,000+ student records —
including Aadhaar numbers, contact details, and academic data.

Coordinated disclosure with CERT-In. Vulnerability patched.
Featured in national Malayalam news outlets.

` + sub.Render("Bug Bounty") + `
Active researcher on HackerOne and Bugcrowd.
Focus: Web app vulnerabilities, API security, business logic flaws,
WAF evasion, misconfiguration assessments.

` + sub.Render("Press") + `
- onmanorama.com — Oct 2025
  https://www.onmanorama.com/news/kerala/2025/10/06/kerala-techie-cybersecurity-hall-of-fame-arshad.html
- mathrubhumi.com — Oct 2025
  https://www.mathrubhumi.com/technology/news/kerala-tech-whiz-fixes-university-security-flaw-ueatbr7i`,

		menuProjects: sub.Render("Triple i Admin Portal") + `
Next.js, NestJS, MongoDB, PostgreSQL, TanStack (Query/Table/Form),
Zustand, Zod, Storybook

- SDUI scheduling module — date-keyed lookup map for instant slot
  resolution, pre-computes all month sections at mount, resolves
  full-day and partial-day instructor availability conflicts
- Fee policy builder with runtime Zod schema switching,
  business-rule validation layer decoupled from schema
  validation, ref-driven derived state — eliminates re-renders
  across complex multi-section forms. Event bus for cross-module
  state propagation.
- Google Drive-style file management: upload, rename, nested
  folders — TanStack Table, Zustand, Zod
- Figma-to-code design system with reusable component library,
  validated through Storybook for visual regression testing


` + sub.Render("Triple i Learning Platform") + `
Next.js, NestJS, PostgreSQL, HLS, WebRTC, SSE, SDUI

- SDUI-driven exam module: objective, descriptive, and
  scenario-based questions with dynamic image rendering —
  backend-controlled layouts and scoring without client redeploys
- NestJS analytics services: score banding, percentile
  computation, multi-cohort comparative reporting — surfaced
  through zoomable client-side charts with rank computation
- Live-class infrastructure with HLS streaming and WebRTC,
  role-based screen-control management, SSE-driven real-time
  schedule notifications`,

		menuSkills: `Languages:
TypeScript, JavaScript, Python, SQL, Go

Frontend:
React, Next.js, Three.js, Tailwind CSS, TanStack Query, Zustand, Redux

Backend:
NestJS, Express.js, Node.js

Database:
MongoDB, Redis, PostgreSQL

DevOps:
Docker, CI/CD, AWS, Coolify, Nixpacks

Tools:
Storybook, Zod, Figma, Claude Code`,

		menuContact: `Email: arshadayanikkal@gmail.com
GitHub: https://github.com/arshadakl
LinkedIn: https://www.linkedin.com/in/arshad-akl/`,
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	m := NewModel()
	m.width = pty.Window.Width
	m.height = pty.Window.Height
	return m, []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
}

func main() {
	s, err := wish.NewServer(
		wish.WithAddress("0.0.0.0:22"),
		wish.WithHostKeyPath("/app/keys/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not start server: %v\n", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("SSH portfolio server running on 0.0.0.0:22")

	go func() {
		if err := s.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
	fmt.Println("Server stopped")
}
