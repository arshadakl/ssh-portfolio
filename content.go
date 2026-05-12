package main

import "github.com/charmbracelet/lipgloss"

type Section struct {
	Key   string
	Label string
	Icon  string
	Lines []string
}

// hyperlink wraps text in an OSC 8 terminal hyperlink.
// Works in iTerm2, Kitty, WezTerm, GNOME Terminal, and modern terminals over SSH.
func hyperlink(url, text string) string {
	return "\x1b]8;;" + url + "\x1b\\" + text + "\x1b]8;;\x1b\\"
}

func buildSections(r *lipgloss.Renderer) []Section {
	return []Section{
		{Key: "whoami", Label: "whoami", Icon: "◈", Lines: buildWhoami(r)},
		{Key: "experience", Label: "experience", Icon: "▸", Lines: buildExperience(r)},
		{Key: "projects", Label: "projects", Icon: "⬡", Lines: buildProjects(r)},
		{Key: "recognition", Label: "recognition", Icon: "⚡", Lines: buildRecognition(r)},
		{Key: "skills", Label: "skills", Icon: "⬢", Lines: buildSkills(r)},
		{Key: "contact", Label: "contact", Icon: "✉", Lines: buildContact(r)},
	}
}

func buildWhoami(r *lipgloss.Renderer) []string {
	orange := styleOrange(r).Bold(true)
	dim    := styleDim(r)
	green  := styleGreen(r)
	text   := styleText(r)

	return []string{
		orange.Render("Arshad A."),
		text.Render("Full-stack engineer + security researcher."),
		dim.Render("2+ years in production engineering."),
		"",
		orange.Render("> What I've Built"),
		"  • Shipped two products at " + green.Render("ELT Global") + " – student-facing LMS and",
		"    admin operations portal.",
		"  • Serving " + green.Render("100,000+") + " monthly users, " + green.Render("10,000+") + " DAU,",
		"    " + green.Render("16M+ API") + " requests/month.",
		"",
		orange.Render("> Tech Stack"),
		"  • " + green.Render("NestJS") + " / " + green.Render("Express.js") + " backends + " + green.Render("Next.js") + " frontends",
		"    in a monorepo. " + green.Render("TypeScript") + " throughout.",
		"",
		orange.Render("> Security Background"),
		"  • " + r.NewStyle().Foreground(colorOrange).Render("CERT-In Hall of Fame") + " (Government of India).",
		"  • Bug bounty focus: web app vulnerabilities, API security,",
		"    and business logic flaws.",
		"",
		orange.Render("> Currently"),
		"  Exploring new opportunities. Open to full-time roles in",
		"  product engineering or security-adjacent engineering.",
	}
}

func buildExperience(r *lipgloss.Renderer) []string {
	orange := styleOrange(r).Bold(true)
	dim    := styleDim(r)
	green  := styleGreen(r)
	text   := styleText(r)

	return []string{
		orange.Render("ELT Global Pvt Ltd") + text.Render(" — Software Engineer"),
		dim.Render("Bangalore, India  |  Aug 2024 – Present"),
		"",
		"Full-stack on EdTech platform at " + green.Render("10k+ DAU") + " — two products: a",
		"student-facing LMS and admin portal. Monorepo:",
		green.Render("NestJS") + " / " + green.Render("Express.js") + " backends + " + green.Render("Next.js") + " frontends.",
		"",
		orange.Render("> Key Work"),
		"  • SDUI-driven interfaces — backend controls screen rendering",
		"    without client redeploys",
		"  • Built automation tooling (Google Apps Script workflows +",
		"    Docker crash monitor with Slack alerting) — eliminated",
		"    " + green.Render("3-4 hrs/day") + " of repetitive overhead",
		"  • Type-safe OpenAPI Swagger + Codegen pipeline — auto-generates",
		"    all API schemas and TypeScript types from a single command,",
		"    reducing frontend-backend integration overhead by " + green.Render("~60%"),
		"  • Security hardening across auth flows, token rotation, RBAC,",
		"    input validation, and request signing",
		"  • Improved page performance " + green.Render("35-40%") + " via code-splitting,",
		"    TanStack Query cache-first patterns, request deduplication,",
		"    API batching",
		"  • Cut key analytics API payload by " + green.Render("~87%") + " through client-side",
		"    derivation",
		"",
		orange.Render("Brototype") + text.Render(" — Full-Stack Engineering Intern"),
		dim.Render("Calicut  |  2023 – 2024"),
		"",
		"  • Built freelancer marketplace with Socket.IO real-time collab,",
		"    Stripe payments, AWS EC2 deployment, ranking algorithm",
		"  • Developed full e-commerce platform: coupon engine, Razorpay,",
		"    session-based auth, admin panel — deployed end-to-end",
	}
}

func buildProjects(r *lipgloss.Renderer) []string {
	orange := styleOrange(r).Bold(true)
	dim    := styleDim(r)
	green  := styleGreen(r)

	return []string{
		orange.Render("Triple i Admin Portal"),
		dim.Render("Next.js, NestJS, MongoDB, PostgreSQL, TanStack (Query/Table/Form),"),
		dim.Render("Zustand, Zod, Storybook"),
		"",
		orange.Render("> Highlights"),
		"  • SDUI scheduling module — date-keyed lookup map for instant slot",
		"    resolution, pre-computes all month sections at mount, resolves",
		"    full-day and partial-day instructor availability conflicts",
		"  • Fee policy builder with runtime Zod schema switching,",
		"    business-rule validation layer decoupled from schema",
		"    validation, ref-driven derived state — eliminates re-renders",
		"    across complex multi-section forms. Event bus for cross-module",
		"    state propagation.",
		"  • Google Drive-style file management: upload, rename, nested",
		"    folders — TanStack Table, Zustand, Zod",
		"  • Figma-to-code design system with reusable component library,",
		"    validated through " + green.Render("Storybook") + " for visual regression testing",
		"",
		"",
		orange.Render("Triple i Learning Platform"),
		dim.Render("Next.js, NestJS, PostgreSQL, HLS, SSE, SDUI"),
		"",
		orange.Render("> Highlights"),
		"  • SDUI-driven exam module: objective, descriptive, and",
		"    scenario-based questions with dynamic image rendering —",
		"    backend-controlled layouts and scoring without client redeploys",
		"  • NestJS analytics services: score banding, percentile",
		"    computation, multi-cohort comparative reporting — surfaced",
		"    through zoomable client-side charts with rank computation",
		"  • Live-class infrastructure with " + green.Render("HLS streaming") + ",",
		"    role-based screen-control management, SSE-driven real-time",
		"    schedule notifications",
	}
}

func buildRecognition(r *lipgloss.Renderer) []string {
	orange := styleOrange(r).Bold(true)
	dim    := styleDim(r)
	green  := styleGreen(r)
	text   := styleText(r)
	link   := styleLink(r)

	return []string{
		orange.Render("CERT-In Hall of Fame") + text.Render(" — Government of India"),
		"",
		"Discovered and responsibly disclosed a critical vulnerability in a",
		"major Kerala university's official website. Exposed unauthenticated",
		"access to backend systems affecting " + green.Render("200,000+ student records") + " —",
		"including Aadhaar numbers, contact details, and academic data.",
		"",
		"Coordinated disclosure with CERT-In. Vulnerability patched.",
		"Featured in national Malayalam news outlets.",
		"",
		orange.Render("> News"),
		"  • " + dim.Render("onmanorama.com") + " — Oct 2025",
		"    " + link.Render(hyperlink(
			"https://www.onmanorama.com/news/kerala/2025/10/06/kerala-techie-cybersecurity-hall-of-fame-arshad.html",
			"onmanorama.com/kerala-techie-cybersecurity-hall-of-fame",
		)),
		"",
		"  • " + dim.Render("mathrubhumi.com") + " — Oct 2025",
		"    " + link.Render(hyperlink(
			"https://www.mathrubhumi.com/technology/news/kerala-tech-whiz-fixes-university-security-flaw-ueatbr7i",
			"mathrubhumi.com/kerala-tech-whiz-fixes-university-security-flaw",
		)),
	}
}

func buildSkills(r *lipgloss.Renderer) []string {
	orange := styleOrange(r).Bold(true)
	green  := styleGreen(r)
	dim    := styleDim(r)

	return []string{
		orange.Render("> Tech"),
		"  " + green.Render("NestJS") + ", " + green.Render("Express.js") + "  —  " + green.Render("React") + " (" + green.Render("Next.js") + "), " + green.Render("TanStack") + " (Query, Table, Form),",
		"  " + green.Render("Zustand") + ", " + green.Render("Redux") + ", " + green.Render("Zod") + ", Tailwind CSS, " + green.Render("Storybook"),
		"",
		orange.Render("> AI"),
		"  " + green.Render("Claude Code") + ", Codex CLI, GitHub Copilot, OpenCode",
		"  " + dim.Render("MCP-integrated workflows") + " (Figma, Supabase, GitHub)",
		"",
		orange.Render("> Infra"),
		"  " + green.Render("AWS") + " (EC2/S3/R2), " + green.Render("Supabase") + ", " + green.Render("Docker") + ", CI/CD,",
		"  Cloudflare, Sentry, Coolify",
		"",
		orange.Render("> Data & Tooling"),
		"  " + green.Render("PostgreSQL") + ", " + green.Render("MongoDB") + "  —  OpenAPI Swagger + Codegen,",
		"  " + green.Render("Socket.IO") + ", Git",
	}
}

func buildContact(r *lipgloss.Renderer) []string {
	orange := styleOrange(r).Bold(true)
	dim    := styleDim(r)
	link   := styleLink(r)

	return []string{
		orange.Render("Get in touch"),
		"",
		"  " + dim.Render("Email   ") + "  arshadayanikkal@gmail.com",
		"  " + dim.Render("Website ") + "  " + link.Render(hyperlink(
			"https://arshadakl.in",
			"arshadakl.in",
		)),
		"  " + dim.Render("GitHub  ") + "  " + link.Render(hyperlink(
			"https://github.com/arshadakl",
			"github.com/arshadakl",
		)),
		"  " + dim.Render("LinkedIn") + "  " + link.Render(hyperlink(
			"https://linkedin.com/in/arshad-akl",
			"linkedin.com/in/arshad-akl",
		)),
		"",
		styleDim(r).Render("  Open to full-time roles in product engineering or"),
		styleDim(r).Render("  security-adjacent engineering."),
		"",
		r.NewStyle().Foreground(colorGreen).Bold(true).Render("  Status: Open to work"),
	}
}
