package main

import "github.com/charmbracelet/lipgloss"

var (
	colorBg       = lipgloss.Color("235") // near-black background
	colorOrange   = lipgloss.Color("208") // primary accent — Claude brand orange
	colorText     = lipgloss.Color("252") // main body text
	colorDim      = lipgloss.Color("241") // secondary/dim text
	colorBorder   = lipgloss.Color("237") // panel borders and dividers
	colorSelBg    = lipgloss.Color("94")  // selected item background (dark amber)
	colorGreen    = lipgloss.Color("71")  // success / open-to-work
	colorStatusBg = lipgloss.Color("236") // footer status bar background
	colorLink     = lipgloss.Color("75")  // clickable hyperlinks (soft blue)
)

func styleOrange(r *lipgloss.Renderer) lipgloss.Style { return r.NewStyle().Foreground(colorOrange) }
func styleText(r *lipgloss.Renderer) lipgloss.Style   { return r.NewStyle().Foreground(colorText) }
func styleDim(r *lipgloss.Renderer) lipgloss.Style    { return r.NewStyle().Foreground(colorDim) }
func styleGreen(r *lipgloss.Renderer) lipgloss.Style  { return r.NewStyle().Foreground(colorGreen) }
func styleLink(r *lipgloss.Renderer) lipgloss.Style   { return r.NewStyle().Foreground(colorLink) }
