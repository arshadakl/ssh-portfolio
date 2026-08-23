package main

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

var markdownLinkPattern = regexp.MustCompile(`\[([^\]]+)\]\((https?://[^\s)]+)\)`)

func sanitizeTerminal(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	var out strings.Builder
	state := byte(0) // 0 normal, 1 ESC, 2 CSI, 3 OSC, 4 OSC after ESC
	for _, r := range value {
		switch state {
		case 1:
			switch r {
			case '[':
				state = 2
			case ']':
				state = 3
			default:
				state = 0
			}
			continue
		case 2:
			if r >= 0x40 && r <= 0x7e {
				state = 0
			}
			continue
		case 3:
			if r == '\a' {
				state = 0
			} else if r == 0x1b {
				state = 4
			}
			continue
		case 4:
			if r == '\\' {
				state = 0
			} else {
				state = 3
			}
			continue
		}

		if r == 0x1b {
			state = 1
			continue
		}
		if r == '\n' || r == '\t' || (!unicode.IsControl(r) && !(r >= 0x7f && r <= 0x9f)) {
			out.WriteRune(r)
		}
	}
	return out.String()
}

func plainMarkdown(value string) string {
	value = markdownLinkPattern.ReplaceAllString(value, "$1 ($2)")
	replacer := strings.NewReplacer("**", "", "__", "", "~~", "", "`", "")
	return strings.TrimSpace(replacer.Replace(value))
}

func wrapTerminalText(value string, width int) []string {
	if width < 8 {
		width = 8
	}
	value = sanitizeTerminal(value)
	var result []string
	for _, paragraph := range strings.Split(value, "\n") {
		if paragraph == "" {
			result = append(result, "")
			continue
		}
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}
		line := ""
		for _, word := range words {
			for runewidth.StringWidth(word) > width {
				if line != "" {
					result = append(result, line)
					line = ""
				}
				part, rest := splitDisplayWidth(word, width)
				result = append(result, part)
				word = rest
			}
			candidate := word
			if line != "" {
				candidate = line + " " + word
			}
			if runewidth.StringWidth(candidate) > width {
				result = append(result, line)
				line = word
			} else {
				line = candidate
			}
		}
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

func splitDisplayWidth(value string, width int) (string, string) {
	var left, right strings.Builder
	used := 0
	for _, r := range value {
		runeWidth := runewidth.RuneWidth(r)
		if used+runeWidth <= width || left.Len() == 0 {
			left.WriteRune(r)
			used += runeWidth
		} else {
			right.WriteRune(r)
		}
	}
	return left.String(), right.String()
}

func renderMarkdownLines(m Model, markdown string, width int) []string {
	var lines []string
	inCode := false
	for _, raw := range strings.Split(sanitizeTerminal(markdown), "\n") {
		trimmed := strings.TrimSpace(raw)
		if strings.HasPrefix(trimmed, "```") {
			inCode = !inCode
			continue
		}
		if trimmed == "" {
			lines = append(lines, "")
			continue
		}
		if inCode {
			for _, line := range wrapTerminalText(raw, width-2) {
				lines = append(lines, styleLink(m.renderer).Render("  "+line))
			}
			continue
		}
		if heading := strings.TrimSpace(strings.TrimLeft(trimmed, "#")); len(trimmed) > len(heading) && strings.HasPrefix(trimmed, "#") {
			for _, line := range wrapTerminalText(plainMarkdown(heading), width) {
				lines = append(lines, styleOrange(m.renderer).Bold(true).Render(line))
			}
			continue
		}
		prefix := ""
		content := trimmed
		if strings.HasPrefix(content, "- ") || strings.HasPrefix(content, "* ") || strings.HasPrefix(content, "+ ") {
			prefix, content = "• ", strings.TrimSpace(content[2:])
		}
		wrapped := wrapTerminalText(plainMarkdown(content), width-runewidth.StringWidth(prefix))
		for index, line := range wrapped {
			if index == 0 {
				lines = append(lines, styleText(m.renderer).Render(prefix+line))
			} else {
				lines = append(lines, styleText(m.renderer).Render(strings.Repeat(" ", runewidth.StringWidth(prefix))+line))
			}
		}
	}
	return trimBlankLines(lines)
}

func renderResponseLines(m Model, response chatResponse, width int) []string {
	var lines []string
	if response.Document.Version == 1 {
		for _, block := range response.Document.Blocks {
			switch block.Type {
			case "markdown":
				lines = append(lines, renderMarkdownLines(m, block.Markdown, width)...)
			case "text":
				for _, line := range wrapTerminalText(block.Text, width) {
					lines = append(lines, styleText(m.renderer).Render(line))
				}
			case "heading":
				for _, line := range wrapTerminalText(block.Text, width) {
					lines = append(lines, styleOrange(m.renderer).Bold(true).Render(line))
				}
			case "list":
				for index, item := range block.Items {
					prefix := "• "
					if block.Style == "ordered" {
						prefix = strconv.Itoa(index+1) + ". "
					}
					wrapped := wrapTerminalText(item, width-runewidth.StringWidth(prefix))
					for lineIndex, line := range wrapped {
						if lineIndex == 0 {
							lines = append(lines, styleText(m.renderer).Render(prefix+line))
						} else {
							lines = append(lines, styleText(m.renderer).Render(strings.Repeat(" ", runewidth.StringWidth(prefix))+line))
						}
					}
				}
			case "table":
				if len(block.Columns) > 0 {
					for _, line := range wrapTerminalText(strings.Join(cleanStrings(block.Columns), " | "), width) {
						lines = append(lines, styleOrange(m.renderer).Bold(true).Render(line))
					}
				}
				for _, row := range block.Rows {
					for _, line := range wrapTerminalText(strings.Join(cleanStrings(row), " | "), width) {
						lines = append(lines, styleText(m.renderer).Render(line))
					}
				}
			case "code":
				for _, line := range strings.Split(sanitizeTerminal(block.Code), "\n") {
					for _, wrapped := range wrapTerminalText(line, width-2) {
						lines = append(lines, styleLink(m.renderer).Render("  "+wrapped))
					}
				}
			case "json":
				var value any
				if len(block.Value) > 0 && json.Unmarshal(block.Value, &value) == nil {
					pretty, _ := json.MarshalIndent(value, "", "  ")
					for _, line := range strings.Split(string(pretty), "\n") {
						for _, wrapped := range wrapTerminalText(line, width-2) {
							lines = append(lines, styleLink(m.renderer).Render("  "+wrapped))
						}
					}
				}
			case "callout":
				if block.Title != nil && strings.TrimSpace(*block.Title) != "" {
					for _, line := range wrapTerminalText(*block.Title, width) {
						lines = append(lines, styleOrange(m.renderer).Bold(true).Render(line))
					}
				}
				for _, line := range wrapTerminalText(block.Text, width-2) {
					lines = append(lines, styleText(m.renderer).Render("│ "+line))
				}
			}
		}
	}
	if len(lines) == 0 && strings.TrimSpace(response.Answer) != "" {
		lines = renderMarkdownLines(m, response.Answer, width)
	}
	return trimBlankLines(lines)
}

func cleanStrings(values []string) []string {
	result := make([]string, len(values))
	for index, value := range values {
		result[index] = sanitizeTerminal(value)
	}
	return result
}

func trimBlankLines(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func buildChat(m Model, width int) []string {
	if width < 12 {
		width = 12
	}
	status := styleGreen(m.renderer).Render("● connected")
	if m.chatClient == nil || m.chatClient.configErr != nil {
		status = styleDim(m.renderer).Render("● unavailable")
	}
	lines := []string{
		styleOrange(m.renderer).Bold(true).Render("✦ ask-arshad") +
			styleDim(m.renderer).Render("  ") + status,
		"",
	}

	if len(m.chatTurns) == 0 && !m.chatPending {
		for _, line := range wrapTerminalText("Ask about Arshad's projects, engineering decisions, experience, or security research.", width) {
			lines = append(lines, styleText(m.renderer).Render(line))
		}
		lines = append(lines, "", styleDim(m.renderer).Render("Try:"))
		for _, suggestion := range []string{
			"What project best represents how Arshad thinks?",
			"What did he contribute at ELT Global?",
			"How does his RAG system work?",
		} {
			wrapped := wrapTerminalText(suggestion, width-2)
			for index, line := range wrapped {
				prefix := "  "
				if index == 0 {
					prefix = "• "
				}
				lines = append(lines, styleDim(m.renderer).Render(prefix+line))
			}
		}
	}

	for _, turn := range m.chatTurns {
		lines = append(lines, styleOrange(m.renderer).Bold(true).Render("You"))
		for _, line := range wrapTerminalText(turn.Question, width-2) {
			lines = append(lines, styleText(m.renderer).Render("  "+line))
		}
		lines = append(lines, "", styleLink(m.renderer).Bold(true).Render("Arshad AI"))
		lines = append(lines, renderResponseLines(m, turn.Response, width-2)...)
		lines = append(lines, "")
	}

	if m.chatPending {
		lines = append(lines, styleOrange(m.renderer).Bold(true).Render("You"))
		for _, line := range wrapTerminalText(m.chatPendingQuestion, width-2) {
			lines = append(lines, styleText(m.renderer).Render("  "+line))
		}
		spinner := []string{"✦", "✧", "✦", "·"}[m.chatSpinner%4]
		lines = append(lines, "", styleLink(m.renderer).Render(spinner+" Thinking…"), "")
	}
	if m.chatNotice != "" {
		for _, line := range wrapTerminalText(m.chatNotice, width) {
			lines = append(lines, styleGreen(m.renderer).Render(line))
		}
		lines = append(lines, "")
	}
	if m.chatError != "" {
		for _, line := range wrapTerminalText("Error: "+m.chatError, width) {
			lines = append(lines, styleOrange(m.renderer).Render(line))
		}
		lines = append(lines, "")
	}

	lines = append(lines, styleDim(m.renderer).Render(strings.Repeat("─", width)))
	if m.chatFocused {
		promptWidth := width - 2
		wrapped := wrapTerminalText(m.chatInput, promptWidth)
		if len(wrapped) == 0 {
			wrapped = []string{""}
		}
		for index, line := range wrapped {
			prefix := "  "
			if index == 0 {
				prefix = "> "
			}
			cursor := ""
			if index == len(wrapped)-1 && m.blinkOn && !m.chatPending && !m.chatClearPending {
				cursor = "▊"
			}
			lines = append(lines, styleOrange(m.renderer).Bold(true).Render(prefix)+styleText(m.renderer).Render(line)+styleOrange(m.renderer).Render(cursor))
		}
	} else {
		lines = append(lines, styleDim(m.renderer).Render("press enter to focus the AI prompt"))
	}
	lines = append(lines, styleDim(m.renderer).Render("enter send · ↑↓/wheel scroll · esc navigate · /clear"))
	return lines
}
