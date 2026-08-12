package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// commandResult is the outcome of executing a command-bar input.
type commandResult struct {
	message      string // feedback shown in the command bar ("" = hide bar)
	jump         int    // section index to navigate to (-1 = stay)
	activate     int    // intent index to expand in contribute (-1 = none)
	startContact bool   // launch the contact form wizard
	quit         bool
}

// execCommand parses and executes a command string. Commands are deliberate
// utilities: navigate, surface a copyable link, or qualify a visitor.
func execCommand(r *lipgloss.Renderer, input string) commandResult {
	res := commandResult{jump: -1, activate: -1}

	cmd := strings.ToLower(strings.TrimSpace(input))
	if cmd == "" {
		return res
	}

	dim  := styleDim(r)
	link := styleLink(r)

	linkMsg := func(label, url, display string) string {
		return dim.Render(label+" → ") + link.Render(hyperlink(url, display))
	}

	// section navigation by name or number
	sections := map[string]int{
		"home": 0, "whoami": 1, "experience": 2, "projects": 3,
		"contribute": 4, "recognition": 5, "skills": 6, "contact": 7,
		"1": 0, "2": 1, "3": 2, "4": 3, "5": 4, "6": 5, "7": 6, "8": 7,
	}
	if idx, ok := sections[cmd]; ok {
		res.jump = idx
		return res
	}

	switch cmd {
	case "help", "?":
		res.message = dim.Render("resume · email · github · linkedin · blog · hire · message · <section> · quit")

	case "resume", "cv":
		res.message = linkMsg("resume", "https://arshadakl.in/docs/arshad_2026.pdf", "arshadakl.in/docs/arshad_2026.pdf")

	case "email", "mail":
		res.message = linkMsg("email", "mailto:arshadayanikkal@gmail.com", "arshadayanikkal@gmail.com")

	case "github", "gh":
		res.message = linkMsg("github", "https://github.com/arshadakl", "github.com/arshadakl")

	case "linkedin", "li":
		res.message = linkMsg("linkedin", "https://linkedin.com/in/arshad-akl", "linkedin.com/in/arshad-akl")

	case "blog":
		res.message = linkMsg("blog", "https://blog.arshadakl.in", "blog.arshadakl.in")

	case "leetcode", "lc":
		res.message = linkMsg("leetcode", "https://leetcode.com/u/arshadakl/", "leetcode.com/u/arshadakl")

	case "hire":
		// conversion shortcut: contribute section, "hiring a developer" expanded
		res.jump = 4
		res.activate = 0

	case "message", "msg":
		res.startContact = true

	case "clear", "cls":
		res.message = ""

	case "quit", "exit", "q":
		res.quit = true

	default:
		res.message = styleDim(r).Render("command not found: " + cmd + " — try 'help'")
	}

	return res
}
