package main

import "github.com/charmbracelet/lipgloss"

// Palette: Tokyo Night neutrals + brand orange hero accent (arshadakl.in).
var (
	colorBg       = lipgloss.Color("#1a1b26") // deep blue-black background
	colorText     = lipgloss.Color("#c0caf5") // soft lavender-white body text
	colorDim      = lipgloss.Color("#565f89") // muted slate secondary text
	colorBorder   = lipgloss.Color("#3b4261") // panel borders and dividers
	colorStatusBg = lipgloss.Color("#1f2335") // footer/command bar background
	colorOrange   = lipgloss.Color("#D9651A") // brand orange — hero accent
	colorGreen    = lipgloss.Color("#9ece6a") // metrics / open-to-work
	colorLink     = lipgloss.Color("#7aa2f7") // clickable hyperlinks
	colorPurple   = lipgloss.Color("#bb9af7") // secondary highlights
	colorPink     = lipgloss.Color("#f7768e") // security/recognition accents
)

func styleOrange(r *lipgloss.Renderer) lipgloss.Style { return r.NewStyle().Foreground(colorOrange) }
func styleText(r *lipgloss.Renderer) lipgloss.Style   { return r.NewStyle().Foreground(colorText) }
func styleDim(r *lipgloss.Renderer) lipgloss.Style    { return r.NewStyle().Foreground(colorDim) }
func styleGreen(r *lipgloss.Renderer) lipgloss.Style  { return r.NewStyle().Foreground(colorGreen) }
func styleLink(r *lipgloss.Renderer) lipgloss.Style   { return r.NewStyle().Foreground(colorLink) }
func stylePurple(r *lipgloss.Renderer) lipgloss.Style { return r.NewStyle().Foreground(colorPurple) }
func stylePink(r *lipgloss.Renderer) lipgloss.Style   { return r.NewStyle().Foreground(colorPink) }

// stylePill renders the selected navigation item: dark text on brand orange.
func stylePill(r *lipgloss.Renderer) lipgloss.Style {
	return r.NewStyle().Background(colorOrange).Foreground(colorBg).Bold(true)
}
