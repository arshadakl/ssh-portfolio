package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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
		{Key: "home", Label: "home", Icon: "⌂", Lines: buildHome(r)},
		{Key: "whoami", Label: "whoami", Icon: "◈", Lines: buildWhoami(r)},
		{Key: "experience", Label: "experience", Icon: "▸", Lines: buildExperience(r)},
		{Key: "projects", Label: "projects", Icon: "⬡", Lines: buildProjects(r)},
		// contribute renders dynamically — see contribute.go
		{Key: "contribute", Label: "contribute", Icon: "→", Lines: nil},
		{Key: "recognition", Label: "recognition", Icon: "⚡", Lines: buildRecognition(r)},
		{Key: "skills", Label: "skills", Icon: "⬢", Lines: buildSkills(r)},
		{Key: "contact", Label: "contact", Icon: "💬", Lines: buildContact(r)},
	}
}

func buildHome(r *lipgloss.Renderer) []string {
	orange := styleOrange(r).Bold(true)
	dim    := styleDim(r)
	green  := styleGreen(r).Bold(true)
	text   := styleText(r)
	pink   := stylePink(r)

	// metric card: rounded border, bold value over dim label
	card := func(value, label string) string {
		return r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Width(20).
			Padding(0, 1).
			Render(green.Render(value) + "\n" + dim.Render(label))
	}

	lines := []string{
		orange.Render("Arshad A") + text.Render(" — Full-stack Engineer & Security Researcher"),
		text.Render("I find security holes before attackers do."),
		"",
	}

	// 2x3 card grid — split joined rows into single lines so scroll math holds
	grid := [][2][2]string{
		{{"100k+", "monthly users"}, {"16M+", "API requests/mo"}},
		{{"200k+", "records secured"}, {"99.9%", "platform uptime"}},
		{{"30+", "enterprise modules"}, {"2+ yrs", "production eng."}},
	}
	for _, pair := range grid {
		row := lipgloss.JoinHorizontal(lipgloss.Top,
			card(pair[0][0], pair[0][1]), " ", card(pair[1][0], pair[1][1]))
		lines = append(lines, strings.Split(row, "\n")...)
	}

	return append(lines,
		"",
		pink.Render("CERT-In Hall of Fame")+text.Render(" (Gov. of India) — national recognition"),
		text.Render("for responsible vulnerability disclosure."),
		"",
		green.Render("Open to: ")+text.Render("Frontend Engineer · Full-stack Engineer roles"),
		"",
		dim.Render("j/k or 1-8 navigate · w/s or wheel scroll"),
		dim.Render("':' commands · ':message' to send a note · 'q' quit"),
	)
}

func buildWhoami(r *lipgloss.Renderer) []string {
	orange := styleOrange(r).Bold(true)
	green  := styleGreen(r)
	text   := styleText(r)
	pink   := stylePink(r)

	return []string{
		orange.Render("Arshad A."),
		text.Render("Full-stack engineer + security researcher."),
		text.Render("I find security holes before attackers do — then build"),
		text.Render("software where those holes don't exist."),
		"",
		orange.Render("> Current Impact"),
		"  • Two production products at " + green.Render("ELT Global") + " (EdTech): student-facing",
		"    LMS + admin operations portal — " + green.Render("10k+ DAU") + ".",
		"  • " + green.Render("100k+ monthly") + " users · " + green.Render("16M+ API") + " requests/mo · " + green.Render("99.9%") + " uptime.",
		"  • " + green.Render("NestJS") + " / " + green.Render("Express.js") + " backends + " + green.Render("Next.js") + " frontends,",
		"    " + green.Render("TypeScript") + " throughout.",
		"",
		orange.Render("> Security Background"),
		"  • " + pink.Render("CERT-In Hall of Fame") + " (Government of India).",
		"  • Protected " + green.Render("200k+ student records") + " via responsible disclosure.",
		"  • Bug bounty focus: web app vulnerabilities, API security,",
		"    and business logic flaws.",
		"",
		orange.Render("> Currently"),
		"  Software Engineer @ ELT Global, Bangalore.",
		"  Open to " + green.Render("Frontend") + " / " + green.Render("Full-stack") + " / " + green.Render("Product Engineer") + " roles.",
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
	link   := styleLink(r)

	return []string{
		orange.Render("Triple i Admin Portal"),
		dim.Render("Next.js, NestJS, PostgreSQL, MongoDB, Redis, TanStack, Zustand,"),
		dim.Render("Zod, Storybook"),
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
		dim.Render("Next.js, NestJS, PostgreSQL, Redis, HLS, WebRTC, SSE, SDUI"),
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
		"  " + dim.Render("Live  ") + "  " + link.Render(hyperlink("https://app.eltglobal.in/", "app.eltglobal.in")),
		"",
		"",
		orange.Render("SSH Portfolio"),
		dim.Render("Go, Wish, Bubbletea, Docker, Nginx, GCP, GitHub Actions"),
		"",
		"Engineered an SSH-based terminal portfolio using Go with",
		"Charm's Wish and Bubbletea libraries, served from a GCP",
		"free-tier e2-micro VM. Same domain serves a Next.js website",
		"over HTTPS and the terminal UI over SSH simultaneously —",
		"containerized with Docker, reverse proxied through Nginx,",
		"deployed via GitHub Actions CI/CD.",
		"  " + dim.Render("GitHub") + "  " + link.Render(hyperlink("https://github.com/arshadakl/ssh-portfolio", "github.com/arshadakl/ssh-portfolio")),
		"  " + dim.Render("Blog  ") + "  " + link.Render(hyperlink("https://blog.arshadakl.in/my-portfolio-has-no-url-just-an-ssh-command", "blog.arshadakl.in — writeup")),
		"  " + dim.Render("Try   ") + "  " + green.Render("ssh arshadakl.in"),
		"",
		"",
		orange.Render("Minecraft Portfolio"),
		dim.Render("Three.js, Voxel, Next.js, Cloudflare Workers"),
		"",
		"A 3D portfolio as a walkable Minecraft-style voxel house.",
		"Scroll moves a camera along a fixed path through the garden",
		"and rooms — About, Hall of Fame, Experience, Projects,",
		"Skills, Contact — inside a living scene: a farmer working",
		"crops, a dog circling the lawn, day/night lighting, passing",
		"rain, and a hidden bug-hunt mini-game.",
		"  " + dim.Render("GitHub") + "  " + link.Render(hyperlink("https://github.com/arshadakl/Minecraft-portfolio", "github.com/arshadakl/Minecraft-portfolio")),
		"  " + dim.Render("Live  ") + "  " + link.Render(hyperlink("https://minecraft.arshadakl.in", "minecraft.arshadakl.in")),
		"",
		"",
		orange.Render("Freelance Marketplace"),
		dim.Render("Next.js, MongoDB, WebRTC, Stripe"),
		"",
		"Developed a freelance platform matching clients with top",
		"talent. Implemented a ranking system based on client feedback.",
		"Offers real-time chat, video conferencing, and secure",
		"payments. Admins manage users and platform activity.",
		"",
		"",
		orange.Render("Specsy — E-commerce Application"),
		dim.Render("Node.js, MongoDB, EJS, Bootstrap"),
		"",
		"Built an eyewear e-commerce platform with MVC architecture.",
		"Secure onboarding with Nodemailer verification, password",
		"reset, session handling. Full admin panel with category",
		"controls, product lifecycle tools, advanced search, filters,",
		"offers, and coupon management.",
		"",
		"",
		orange.Render("Docker Container Crash Monitor"),
		dim.Render("Bash, Docker, Slack Webhooks, Linux"),
		"",
		"Lightweight Bash monitor for Docker container health —",
		"watches for silent container exits and sends Slack alerts",
		"with container name, ID, image, exit code, runtime, and",
		"host details. Simple production visibility without log-",
		"tailing.",
		"  " + dim.Render("GitHub") + "  " + link.Render(hyperlink("https://github.com/arshadakl/Docker-Crash-Monitor", "github.com/arshadakl/Docker-Crash-Monitor")),
		"",
		"",
		orange.Render("Indian Stock Backtester"),
		dim.Render("Python, Angel One SmartAPI, Pandas, NSE"),
		"",
		"30-minute Opening Range Breakout backtesting system for NSE",
		"stocks. Simulates strict intraday long-only ORB strategy:",
		"one trade/stock/day, breakout entries after first 30 min,",
		"fixed 1.5R target, OR-low stop loss, 15:15 IST square-off.",
		"",
		"",
		orange.Render("Bulk Image Compressor"),
		dim.Render("Python, Pillow, Image Processing"),
		"",
		"Production-ready Python image compression tool that targets",
		"exact output file sizes with minimal quality loss. Binary",
		"search finds the best quality setting, preserves dimensions,",
		"supports JPEG/PNG/WebP, handles batch folders, reports",
		"compression stats, gracefully skips corrupted files.",
		"  " + dim.Render("GitHub") + "  " + link.Render(hyperlink("https://github.com/arshadakl/Bulk-Image-Compressor", "github.com/arshadakl/Bulk-Image-Compressor")),
		"",
		"",
		orange.Render("UTF2TTF"),
		dim.Render("HTML, JavaScript, Static API, GitHub Pages"),
		"",
		"Static Malayalam Unicode to ASCII/TTF converter for legacy",
		"font workflows in DaVinci Resolve, Premiere Pro, Photoshop,",
		"and CapCut. Instant conversion, one-click copy, JSON API",
		"endpoint — works with Apple Shortcuts for fast conversion",
		"from any app.",
		"  " + dim.Render("GitHub") + "  " + link.Render(hyperlink("https://github.com/arshadakl/UTF2TTF", "github.com/arshadakl/UTF2TTF")),
		"  " + dim.Render("Live  ") + "  " + link.Render(hyperlink("https://arshadakl.github.io/UTF2TTF/", "arshadakl.github.io/UTF2TTF")),
		"",
		"",
		orange.Render("StreamHub"),
		dim.Render("Next.js 16, TypeScript, Streaming UI, Vercel"),
		"",
		"Live TV streaming directory with 20,000+ channels from",
		"180+ countries. Apple TV-inspired interface for browsing",
		"and discovering channels across regions.",
		"  " + dim.Render("GitHub") + "  " + link.Render(hyperlink("https://github.com/arshadakl/streamhub", "github.com/arshadakl/streamhub")),
		"  " + dim.Render("Live  ") + "  " + link.Render(hyperlink("https://streamhub-arshad.vercel.app/", "streamhub-arshad.vercel.app")),
	}
}

func buildRecognition(r *lipgloss.Renderer) []string {
	orange := styleOrange(r).Bold(true)
	dim    := styleDim(r)
	green  := styleGreen(r)
	text   := styleText(r)
	link   := styleLink(r)

	return []string{
		stylePink(r).Bold(true).Render("CERT-In Hall of Fame") + text.Render(" — Government of India"),
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
		"  " + green.Render("AWS") + " (EC2/S3) · " + green.Render("GCP") + " · " + green.Render("Docker") + " · CI/CD",
		"  " + green.Render("Cloudflare") + " (Workers · R2 · D1)",
		"  " + green.Render("Supabase") + " · Sentry · Coolify",
		"",
		orange.Render("> Data & Tooling"),
		"  " + green.Render("PostgreSQL") + ", " + green.Render("MongoDB") + "  —  OpenAPI Swagger + Codegen,",
		"  " + green.Render("Socket.IO") + ", Git",
	}
}

func buildContact(r *lipgloss.Renderer) []string {
	orange := styleOrange(r).Bold(true)
	dim    := styleDim(r)
	text   := styleText(r)
	link   := styleLink(r)
	green  := styleGreen(r)

	btn := func(label string) string {
		return stylePill(r).Padding(0, 1).Render(label)
	}

	return []string{
		orange.Render("💬 Let's Talk"),
		text.Render("I reply to every note — usually within 24h."),
		"",
		"  " + btn(" ENTER ") + text.Render("  opens the message form") + dim.Render("   · or type ") + btn(" :message "),
		"",
		text.Render("  Available: ") + dim.Render("Full-time · Contract · Selected freelance"),
		"",
		orange.Render("> Elsewhere"),
		"  " + dim.Render("Resume  ") + "  " + link.Render(hyperlink(
			"https://arshadakl.in/docs/arshad_2026.pdf",
			"arshadakl.in/docs/arshad_2026.pdf",
		)),
		"  " + dim.Render("LinkedIn") + "  " + link.Render(hyperlink(
			"https://linkedin.com/in/arshad-akl",
			"linkedin.com/in/arshad-akl",
		)),
		"  " + dim.Render("GitHub  ") + "  " + link.Render(hyperlink(
			"https://github.com/arshadakl",
			"github.com/arshadakl",
		)),
		"  " + dim.Render("LeetCode") + "  " + link.Render(hyperlink(
			"https://leetcode.com/u/arshadakl/",
			"leetcode.com/u/arshadakl",
		)),
		"  " + dim.Render("Blog    ") + "  " + link.Render(hyperlink(
			"https://blog.arshadakl.in",
			"blog.arshadakl.in",
		)),
		"  " + dim.Render("Website ") + "  " + link.Render(hyperlink(
			"https://arshadakl.in",
			"arshadakl.in",
		)),
		"",
		text.Render("  Prefer email? ") + link.Render(hyperlink(
			"mailto:arshadayanikkal@gmail.com",
			"arshadayanikkal@gmail.com",
		)),
		"",
		green.Bold(true).Render("  Status: Open to work — Frontend / Full-stack roles"),
	}
}

// buildContactForm renders the in-section message form. The active field is
// highlighted with a blinking cursor; Enter on the message field submits.
func buildContactForm(m Model) []string {
	r := m.renderer
	orange := styleOrange(r).Bold(true)
	dim    := styleDim(r)
	text   := styleText(r)

	if m.contactSubmitting {
		return []string{
			orange.Render("> Sending your message..."),
			"",
			dim.Render("  one moment — this goes straight to Arshad's inbox."),
		}
	}

	labels := []string{"name", "email", "subject", "message"}
	values := []string{m.contactName, m.contactEmail, m.contactSubject, m.contactMessage}

	lines := []string{
		orange.Render("💬 Send a Note"),
		text.Render("  Fill in the fields below. Enter or Tab moves to the next,"),
		text.Render("  Enter on message sends, Esc cancels."),
		"",
	}

	for i := range labels {
		label := text.Render(labels[i])
		value := values[i]
		cursor := " "
		if i == m.contactField {
			label = orange.Render(labels[i])
			if m.blinkOn {
				cursor = orange.Render("\u258A")
			}
		}
		lines = append(lines, "  "+label+strings.Repeat(" ", 8-len(labels[i]))+": "+value+cursor)
	}

	if m.contactResult != "" {
		lines = append(lines, "", dim.Render("  "+m.contactResult))
	}

	lines = append(lines,
		"",
		dim.Render("  enter/tab next · enter on message sends · esc cancel"),
	)
	return lines
}
