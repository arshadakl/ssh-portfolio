package main

import "github.com/charmbracelet/lipgloss"

// Palette: warm, muted terminal colors inspired by Claude Code.
var (
	colorBg         = lipgloss.Color("#151513") // warm charcoal background
	colorPanel      = lipgloss.Color("#1D1D1A") // cards and raised surfaces
	colorText       = lipgloss.Color("#E7D9C4") // warm cream primary text
	colorDim        = lipgloss.Color("#A79C89") // warm taupe secondary text
	colorBorder     = lipgloss.Color("#615D51") // muted warm borders and dividers
	colorStatusBg   = lipgloss.Color("#1D1D1A") // footer/command bar background
	colorOrange     = lipgloss.Color("#E97B55") // primary coral accent
	colorAccentSoft = lipgloss.Color("#EF875F") // selected navigation highlight
	colorGreen      = lipgloss.Color("#B6D56A") // online / open-to-work status
	colorLink       = lipgloss.Color("#81938D") // muted teal hyperlinks
	colorPurple     = lipgloss.Color("#81938D") // muted teal secondary highlights
	colorPink       = lipgloss.Color("#E97B55") // recognition accent
)

func styleOrange(r *lipgloss.Renderer) lipgloss.Style { return r.NewStyle().Foreground(colorOrange) }
func styleText(r *lipgloss.Renderer) lipgloss.Style   { return r.NewStyle().Foreground(colorText) }
func styleDim(r *lipgloss.Renderer) lipgloss.Style    { return r.NewStyle().Foreground(colorDim) }
func styleGreen(r *lipgloss.Renderer) lipgloss.Style  { return r.NewStyle().Foreground(colorGreen) }
func styleLink(r *lipgloss.Renderer) lipgloss.Style   { return r.NewStyle().Foreground(colorLink) }
func stylePurple(r *lipgloss.Renderer) lipgloss.Style { return r.NewStyle().Foreground(colorPurple) }
func stylePink(r *lipgloss.Renderer) lipgloss.Style   { return r.NewStyle().Foreground(colorPink) }

// stylePill renders the selected navigation item: dark text on soft coral.
func stylePill(r *lipgloss.Renderer) lipgloss.Style {
	return r.NewStyle().Background(colorAccentSoft).Foreground(colorBg).Bold(true)
}
