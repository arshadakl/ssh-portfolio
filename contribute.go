package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Intent is a visitor-qualified path shown at the bottom of the contribute
// section. The cursor highlights one; activating it expands a tailored panel.
type Intent struct {
	Key   string
	Label string
	Panel func(r *lipgloss.Renderer) []string
}

var intents = []Intent{
	{Key: "hiring", Label: "hiring a developer", Panel: panelHiring},
	{Key: "product", Label: "product to build", Panel: panelProduct},
	{Key: "improve", Label: "improve my product", Panel: panelImprove},
	{Key: "opportunity", Label: "an opportunity", Panel: panelOpportunity},
}

// buildContribute renders the section dynamically based on intent state.
// cursor = highlighted intent row, active = expanded intent (-1 = none).
func buildContribute(r *lipgloss.Renderer, cursor, active int) []string {
	orange := styleOrange(r).Bold(true)
	dim    := styleDim(r)
	green  := styleGreen(r)
	text   := styleText(r)
	link   := styleLink(r)

	proof := func(s string) string { return dim.Render("    proof → ") + green.Render(s) }

	lines := []string{
		text.Render("I don't just write features. I turn complex product"),
		text.Render("requirements into maintainable, production-ready software."),
		"",
		orange.Render("> How I Can Contribute"),
		"",
		"  " + orange.Render("▸ Build production applications"),
		dim.Render("    React / Next.js apps, APIs, dashboards, internal"),
		dim.Render("    platforms, SaaS products."),
		proof("Triple i Learning Platform:") + " " + link.Render(hyperlink("https://app.eltglobal.in/", "app.eltglobal.in")),
		"",
		"  " + orange.Render("▸ Solve complex frontend problems"),
		dim.Render("    Scheduling engines, dynamic forms, large data tables,"),
		dim.Render("    state-heavy UIs, performance issues."),
		proof("Triple i Admin Portal — O(1) scheduling, runtime Zod"),
		green.Render("    schema switching (see projects)"),
		"",
		"  " + orange.Render("▸ Improve existing products"),
		dim.Render("    Refactoring, performance optimization, component"),
		dim.Render("    architecture, API integration, TypeScript adoption."),
		proof("+35-40% page perf · -87% API payload · -60% integration"),
		green.Render("    overhead (at ELT Global)"),
		"",
		"  " + orange.Render("▸ Build scalable UI systems"),
		dim.Render("    Figma → reusable components → Storybook → production."),
		proof("design system powering 30+ enterprise modules"),
		"",
		"  " + orange.Render("▸ Work across the stack"),
		dim.Render("    Frontend-heavy full-stack: Node.js, NestJS, PostgreSQL,"),
		dim.Render("    REST APIs, Docker."),
		proof("this terminal — Go + SSH + GCP + CI/CD") + "  " + link.Render(hyperlink("https://github.com/arshadakl/ssh-portfolio", "src")),
		"",
		orange.Render("> What I Bring"),
		"",
	}

	bring := [][2]string{
		{"Production experience", "real enterprise software, 30+ modules"},
		{"Performance mindset", "data-heavy interfaces at 16M+ req/mo scale"},
		{"System thinking", "design systems, components, API architecture"},
		{"Ownership", "frontend, backend, deployment — end to end"},
		{"Security mindset", "CERT-In recognized research & disclosure"},
	}
	for _, row := range bring {
		lines = append(lines,
			"  "+text.Render(fmt.Sprintf("%-22s", row[0]))+dim.Render("→ "+row[1]))
	}

	lines = append(lines, "", orange.Render("> What Brings You Here?"), "")

	// 2x2 button grid. Selected = filled pill (same language as the nav bar),
	// expanded = green dot. Labels padded so both columns align.
	intentItem := func(i int) string {
		label := fmt.Sprintf("%-18s", intents[i].Label)
		switch {
		case i == cursor && i == active:
			return stylePill(r).Padding(0, 2).Render(label) + green.Render(" ●")
		case i == cursor:
			return stylePill(r).Padding(0, 2).Render(label)
		case i == active:
			return text.Render("[ "+label+" ]") + green.Render(" ●")
		default:
			return dim.Render("[ " + label + " ]")
		}
	}

	lines = append(lines,
		"  "+intentItem(0)+"  "+intentItem(1),
		"",
		"  "+intentItem(2)+"  "+intentItem(3),
		"",
		dim.Render("  ←/→ choose · enter view · esc close"),
	)

	if active >= 0 && active < len(intents) {
		lines = append(lines, "")
		lines = append(lines, intents[active].Panel(r)...)
	}

	return lines
}

func panelHeader(r *lipgloss.Renderer, title string) string {
	dashes := 40 - len(title)
	if dashes < 4 {
		dashes = 4
	}
	return stylePurple(r).Render("  ── " + title + " " + strings.Repeat("─", dashes))
}

func panelHiring(r *lipgloss.Renderer) []string {
	dim  := styleDim(r)
	text := styleText(r)
	link := styleLink(r)

	return []string{
		panelHeader(r, "hiring a developer"),
		"",
		text.Render("  Frontend / Full-stack Engineer — 2+ years production."),
		dim.Render("  React · Next.js · TypeScript · Node.js · NestJS"),
		"",
		text.Render("  Interested in: ") + dim.Render("Frontend Engineer · Full-stack"),
		dim.Render("  Engineer · Product Engineer"),
		"",
		"  → " + dim.Render("resume   ") + " " + link.Render(hyperlink("https://arshadakl.in/docs/arshad_2026.pdf", "arshadakl.in/docs/arshad_2026.pdf")),
		"  → " + dim.Render("linkedin ") + " " + link.Render(hyperlink("https://linkedin.com/in/arshad-akl", "linkedin.com/in/arshad-akl")),
		"  → " + dim.Render("email    ") + " " + link.Render(hyperlink("mailto:arshadayanikkal@gmail.com", "arshadayanikkal@gmail.com")),
		"",
		styleDim(r).Render("  or type :message to open the message form"),
	}
}

func panelProduct(r *lipgloss.Renderer) []string {
	dim  := styleDim(r)
	green := styleGreen(r)
	text := styleText(r)
	link := styleLink(r)

	return []string{
		panelHeader(r, "product to build"),
		"",
		text.Render("  0 → production, end to end."),
		"",
		"  • EdTech platform — SDUI exams, HLS live classes, analytics;",
		"    " + green.Render("100k+ monthly users") + ", 99.9% uptime",
		"  • Marketplace — real-time collab, Stripe payments, ranking",
		"",
		dim.Render("  Stack: Next.js · NestJS · PostgreSQL · Docker · AWS/GCP"),
		"",
		"  → " + dim.Render("discuss scope ") + " " + link.Render(hyperlink("mailto:arshadayanikkal@gmail.com", "arshadayanikkal@gmail.com")),
	}
}

func panelImprove(r *lipgloss.Renderer) []string {
	dim   := styleDim(r)
	green := styleGreen(r)
	text  := styleText(r)
	link  := styleLink(r)

	return []string{
		panelHeader(r, "improve my product"),
		"",
		text.Render("  Proof in numbers, at production scale:"),
		"",
		"  • " + green.Render("+35-40%") + " page performance — code-splitting, cache-first",
		"    patterns, request deduplication, API batching",
		"  • " + green.Render("-87%") + " analytics API payload — client-side derivation",
		"  • " + green.Render("-60%") + " frontend-backend integration overhead — OpenAPI",
		"    Swagger + Codegen pipeline",
		"  • Design systems — Figma → Storybook → production",
		"",
		"  → " + dim.Render("email ") + " " + link.Render(hyperlink("mailto:arshadayanikkal@gmail.com", "arshadayanikkal@gmail.com")),
	}
}

func panelOpportunity(r *lipgloss.Renderer) []string {
	dim   := styleDim(r)
	green := styleGreen(r)
	text  := styleText(r)
	link  := styleLink(r)

	return []string{
		panelHeader(r, "an opportunity"),
		"",
		text.Render("  Available: ") + dim.Render("Full-time · Contract · Selected freelance"),
		text.Render("  Currently: ") + dim.Render("Software Engineer @ ELT Global, Bangalore"),
		text.Render("  Response:  ") + dim.Render("usually within 24h"),
		"",
		"  → " + dim.Render("email    ") + " " + link.Render(hyperlink("mailto:arshadayanikkal@gmail.com", "arshadayanikkal@gmail.com")),
		"  → " + dim.Render("linkedin ") + " " + link.Render(hyperlink("https://linkedin.com/in/arshad-akl", "linkedin.com/in/arshad-akl")),
		"",
		green.Bold(true).Render("  Status: Open to work"),
		styleDim(r).Render("  type :message to open the message form"),
	}
}
