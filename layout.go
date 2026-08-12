package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)


// ── splash ──────────────────────────────────────────────────────────────────

func (m Model) renderSplash() string {
	splashStyle := m.renderer.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Width(m.width).Height(m.height).
		Background(colorBg).Foreground(colorOrange)

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
	visible := m.animFrame
	if visible > len(fullText) {
		visible = len(fullText)
	}
	display := fullText[:visible]
	if visible < len(fullText) {
		display += "▊"
	}
	return splashStyle.Render(asciiArt + "\n\n" + display)
}

// ── wide layout ──────────────────────────────────────────────────────────────

func (m Model) renderWideLayout() string {
	sidebarW := m.sidebarWidth()
	contentW := m.width - sidebarW - 1
	bodyH := m.height - 2 - 1 // 2 header rows + 1 footer row
	if m.cmdBarVisible() {
		bodyH--
	}

	header  := m.renderHeader()
	sidebar := m.renderSidebar(sidebarW, bodyH)
	content := m.renderContent(contentW, bodyH)
	body    := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
	footer  := m.renderFooter(m.width, false)

	if m.cmdBarVisible() {
		return lipgloss.JoinVertical(lipgloss.Left, header, body, m.renderCommandBar(m.width), footer)
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m Model) sidebarWidth() int {
	w := int(float64(m.width) * 0.28)
	if w < 26 {
		w = 26
	}
	if w > 32 {
		w = 32
	}
	return w
}

// ── header ───────────────────────────────────────────────────────────────────

func (m Model) renderHeader() string {
	leftText  := styleOrange(m.renderer).Bold(true).Render("arshad@portfolio:~")
	rightText := styleGreen(m.renderer).Render("● online")
	if m.clientIP != "" {
		rightText = styleDim(m.renderer).Render("visitor "+m.clientIP+"  ") + rightText
	}

	leftW  := lipgloss.Width(leftText)
	rightW := lipgloss.Width(rightText)
	centerW := m.width - leftW - rightW - 2
	if centerW < 0 {
		centerW = 0
	}

	const icon = "✦"
	const slot = 3
	numIcons := centerW / slot
	if numIcons < 1 {
		numIcons = 1
	}
	icons := strings.Repeat(icon+"  ", numIcons)
	for lipgloss.Width(icons) > centerW {
		icons = icons[:len(icons)-1]
	}
	centerCell := m.renderer.NewStyle().
		Width(centerW).
		Align(lipgloss.Center).
		Foreground(colorBorder).
		Render(icons)

	line := m.renderer.NewStyle().
		Background(colorBg).
		Width(m.width).
		Padding(0, 1).
		Render(leftText + centerCell + rightText)

	sep := m.renderer.NewStyle().
		Foreground(colorBorder).
		Width(m.width).
		Render(strings.Repeat("─", m.width))

	return lipgloss.JoinVertical(lipgloss.Left, line, sep)
}

// ── sidebar ──────────────────────────────────────────────────────────────────

func (m Model) renderSidebar(w, h int) string {
	var sb strings.Builder

	navLabel := styleOrange(m.renderer).Bold(true).Render("Navigation")
	sb.WriteString(navLabel + "\n\n")

	for i, sec := range m.sections {
		label := sec.Label
		if sec.Key == "contact" {
			label = "💬 " + sec.Label
		}
		if i == m.selected {
			row := stylePill(m.renderer).Padding(0, 1).Render(fmt.Sprintf("[%d] %s", i+1, label))
			sb.WriteString(row + "\n")
		} else {
			row := styleDim(m.renderer).Render(fmt.Sprintf(" [%d] %s", i+1, label))
			sb.WriteString(row + "\n")
		}
	}

	divider := m.renderer.NewStyle().Foreground(colorBorder).Render(strings.Repeat("─", w-2))
	sb.WriteString("\n" + divider + "\n\n")

	quickLabel := styleOrange(m.renderer).Bold(true).Render("Quick Info")
	sb.WriteString(quickLabel + "\n\n")

	infoRows := []struct{ label, value, color string }{
		{"Location  ", "India", ""},
		{"Experience", "2+ years", ""},
		{"Focus     ", "Full Stack", ""},
		{"Status    ", "Open to work", "green"},
	}
	for _, row := range infoRows {
		label := styleDim(m.renderer).Render(row.label + " : ")
		var val string
		if row.color == "green" {
			val = styleGreen(m.renderer).Bold(true).Render(row.value)
		} else {
			val = styleText(m.renderer).Render(row.value)
		}
		sb.WriteString(label + val + "\n")
	}

	sidebarStyle := m.renderer.NewStyle().
		Width(w).
		Height(h).
		Background(colorBg).
		Padding(0, 1).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorBorder)

	return sidebarStyle.Render(sb.String())
}

// ── content ───────────────────────────────────────────────────────────────────

// renderDivider builds a labeled section divider: ──┤ label ├────────
func (m Model) renderDivider(label string, w int) string {
	pre  := "──┤ "
	post := " ├"
	fill := w - lipgloss.Width(pre+label+post)
	if fill < 0 {
		fill = 0
	}
	return styleDim(m.renderer).Render(pre) +
		styleOrange(m.renderer).Bold(true).Render(label) +
		styleDim(m.renderer).Render(post+strings.Repeat("─", fill))
}

func (m Model) renderContent(w, h int) string {
	sec    := m.sections[m.selected]
	innerW := w - 4

	header := m.renderDivider(sec.Label, innerW) + "\n\n"
	const headerLines = 2

	lines     := m.visibleLines()
	available := h - headerLines - 1
	if available < 1 {
		available = 1
	}

	maxScroll := len(lines) - available
	if maxScroll < 0 {
		maxScroll = 0
	}

	start := m.scrollPos
	if start > maxScroll {
		start = maxScroll
	}
	end := start + available
	if end > len(lines) {
		end = len(lines)
	}

	visible := lines[start:end]
	for len(visible) < available {
		visible = append(visible, "")
	}

	pageInfo := fmt.Sprintf("%d/%d", m.selected+1, len(m.sections))
	if maxScroll > 0 {
		pct := int(float64(start) / float64(maxScroll) * 100)
		pageInfo = fmt.Sprintf("%d%%  %s", pct, pageInfo)
	}
	pageIndicator := m.renderer.NewStyle().
		Foreground(colorDim).
		Width(innerW).
		Align(lipgloss.Right).
		Render(pageInfo)

	body := header + strings.Join(visible, "\n") + "\n" + pageIndicator

	return m.renderer.NewStyle().
		Width(w).
		Height(h).
		Background(colorBg).
		Padding(0, 2).
		Render(body)
}

// ── footer ────────────────────────────────────────────────────────────────────

func (m Model) renderFooter(w int, narrow bool) string {
	var keys string
	if narrow {
		keys = styleDim(m.renderer).Render("←→") + styleText(m.renderer).Render(" navigate") +
			"  " + styleOrange(m.renderer).Render("w/s") + styleText(m.renderer).Render(" scroll") +
			"  " + styleDim(m.renderer).Render(":") + styleText(m.renderer).Render(" cmd") +
			"  " + styleDim(m.renderer).Render("q") + styleText(m.renderer).Render(" quit")
	} else {
		keys = styleDim(m.renderer).Render("j/k ↑↓") + styleText(m.renderer).Render(" navigate") +
			"  " + styleOrange(m.renderer).Render("w/s") + styleText(m.renderer).Render(" scroll") +
			"  " + styleDim(m.renderer).Render("1-8") + styleText(m.renderer).Render(" jump") +
			"  " + styleDim(m.renderer).Render(":") + styleText(m.renderer).Render(" cmd") +
			"  " + styleDim(m.renderer).Render("q") + styleText(m.renderer).Render(" quit")
	}

	elapsed := time.Since(m.sessionStart).Round(time.Second)
	timer := fmt.Sprintf("%02d:%02d", int(elapsed.Minutes()), int(elapsed.Seconds())%60)
	var tip string
	if narrow {
		tip = styleDim(m.renderer).Render(timer)
	} else {
		tip = styleDim(m.renderer).Render("session " + timer)
	}

	keysW := lipgloss.Width(keys)
	tipW  := lipgloss.Width(tip)
	gap   := w - keysW - tipW - 4
	if gap < 1 {
		gap = 1
	}

	return m.renderer.NewStyle().
		Background(colorStatusBg).
		Foreground(colorDim).
		Width(w).
		Padding(0, 1).
		Render(keys + strings.Repeat(" ", gap) + tip)
}

// ── command bar ──────────────────────────────────────────────────────────────

func (m Model) cmdBarVisible() bool {
	return m.cmdMode || m.cmdMsg != ""
}

func (m Model) renderCommandBar(w int) string {
	var content string

	if m.cmdMode {
		cursor := " "
		if m.blinkOn {
			cursor = styleOrange(m.renderer).Render("\u258A")
		}
		content = styleOrange(m.renderer).Bold(true).Render(":") +
			styleText(m.renderer).Render(m.cmdInput) + cursor
	} else {
		content = m.cmdMsg
	}

	return m.renderer.NewStyle().
		Background(colorStatusBg).
		Width(w).
		Padding(0, 1).
		Render(content)
}

// ── narrow layout ─────────────────────────────────────────────────────────────

func (m Model) renderNarrowLayout() string {
	header  := m.renderNarrowHeader()
	navbar  := m.renderNarrowNavbar()
	content := m.renderNarrowContent()
	footer  := m.renderFooter(m.width, true)

	if m.cmdBarVisible() {
		return lipgloss.JoinVertical(lipgloss.Left, header, navbar, content, m.renderCommandBar(m.width), footer)
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, navbar, content, footer)
}

func (m Model) renderNarrowHeader() string {
	left  := styleOrange(m.renderer).Bold(true).Render("arshad@portfolio:~")
	right := styleGreen(m.renderer).Render("● online")
	if m.clientIP != "" && m.width >= 62 {
		right = styleDim(m.renderer).Render("visitor "+m.clientIP+"  ") + right
	}
	gap   := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 1 {
		gap = 1
	}
	return m.renderer.NewStyle().
		Background(colorBg).
		Width(m.width).
		Padding(0, 1).
		Render(left + strings.Repeat(" ", gap) + right)
}

func (m Model) renderNarrowNavbar() string {
	var sb strings.Builder
	for i, sec := range m.sections {
		var item string
		if i == m.selected {
			label := sec.Label
			if sec.Key == "contact" {
				label = "💬 " + sec.Label
			}
			item = stylePill(m.renderer).Padding(0, 1).Render(label)
		} else {
			// icon-only keeps 8 sections on one row in narrow terminals
			item = styleDim(m.renderer).Padding(0, 1).Render(sec.Icon)
		}
		sb.WriteString(item)
		if i < len(m.sections)-1 {
			sb.WriteString(" ")
		}
	}

	return m.renderer.NewStyle().
		Background(colorBg).
		Width(m.width).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorBorder).
		Render(sb.String())
}

func (m Model) renderNarrowContent() string {
	h := m.height - 4
	if m.cmdBarVisible() {
		h--
	}
	if h < 1 {
		h = 1
	}

	sec    := m.sections[m.selected]
	innerW := m.width - 4

	header := m.renderDivider(sec.Label, innerW) + "\n\n"
	const headerLines = 2

	lines     := m.visibleLines()
	available := h - headerLines - 1
	if available < 1 {
		available = 1
	}

	maxScroll := len(lines) - available
	if maxScroll < 0 {
		maxScroll = 0
	}
	start := m.scrollPos
	if start > maxScroll {
		start = maxScroll
	}
	end := start + available
	if end > len(lines) {
		end = len(lines)
	}

	visible := lines[start:end]
	for len(visible) < available {
		visible = append(visible, "")
	}

	pageInfo := fmt.Sprintf("%d/%d", m.selected+1, len(m.sections))
	if maxScroll > 0 {
		pct := int(float64(start) / float64(maxScroll) * 100)
		pageInfo = fmt.Sprintf("%d%%  %s", pct, pageInfo)
	}
	pageIndicator := m.renderer.NewStyle().
		Foreground(colorDim).
		Width(innerW).
		Align(lipgloss.Right).
		Render(pageInfo)

	body := header + strings.Join(visible, "\n") + "\n" + pageIndicator

	return m.renderer.NewStyle().
		Width(m.width).
		Height(h).
		Background(colorBg).
		Padding(0, 2).
		Render(body)
}
