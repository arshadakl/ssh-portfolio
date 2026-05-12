package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ViewMode int

const (
	ViewLoading ViewMode = iota
	ViewNormal
)

type tickMsg time.Time

type Model struct {
	width     int
	height    int
	selected  int
	scrollPos int
	mode      ViewMode
	animFrame int
	sections  []Section
}

func NewModel() Model {
	return Model{
		mode:     ViewLoading,
		sections: buildSections(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) isNarrow() bool { return m.width < 70 }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if m.mode == ViewLoading {
			m.animFrame++
			fullText := "Accessing Arshad's system..."
			if m.animFrame > len(fullText)+5 {
				m.mode = ViewNormal
			}
			return m, tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}

	case tea.MouseMsg:
		if m.mode != ViewLoading && msg.Action == tea.MouseActionPress {
			switch msg.Button {
			case tea.MouseButtonWheelUp:
				if m.scrollPos > 0 {
					m.scrollPos--
				}
			case tea.MouseButtonWheelDown:
				m.scrollPos++
			}
		}

	case tea.KeyMsg:
		if m.mode == ViewLoading {
			return m, nil
		}

		// Arrow keys navigate sections (type-based for SSH compat)
		switch msg.Type {
		case tea.KeyUp:
			m.navigatePrev()
			return m, nil
		case tea.KeyDown:
			m.navigateNext()
			return m, nil
		case tea.KeyLeft:
			if m.isNarrow() {
				m.navigatePrev()
			}
			return m, nil
		case tea.KeyRight:
			if m.isNarrow() {
				m.navigateNext()
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		// section navigation — j/k and ↑/↓ both navigate
		case "j", "down":
			m.navigateNext()
		case "k", "up":
			m.navigatePrev()
		case "right":
			if m.isNarrow() {
				m.navigateNext()
			}
		case "left":
			if m.isNarrow() {
				m.navigatePrev()
			}

		// content scrolling — w/s (mouse wheel also works)
		case "w":
			m.scrollUp()
		case "s":
			m.scrollDown()
		case "pgup", "ctrl+u":
			for i := 0; i < 10; i++ {
				m.scrollUp()
			}
		case "pgdown", "ctrl+d":
			for i := 0; i < 10; i++ {
				m.scrollDown()
			}

		// direct section jump
		case "1":
			m.jumpTo(0)
		case "2":
			m.jumpTo(1)
		case "3":
			m.jumpTo(2)
		case "4":
			m.jumpTo(3)
		case "5":
			m.jumpTo(4)
		case "6":
			m.jumpTo(5)

		case "enter":
			m.scrollPos = 0
		}
	}
	return m, nil
}

func (m *Model) navigateNext() {
	if m.selected < len(m.sections)-1 {
		m.selected++
		m.scrollPos = 0
	}
}

func (m *Model) navigatePrev() {
	if m.selected > 0 {
		m.selected--
		m.scrollPos = 0
	}
}

func (m *Model) jumpTo(idx int) {
	if idx >= 0 && idx < len(m.sections) {
		m.selected = idx
		m.scrollPos = 0
		m.mode = ViewNormal
	}
}

func (m *Model) scrollUp() {
	if m.scrollPos > 0 {
		m.scrollPos--
	}
}

func (m *Model) scrollDown() {
	m.scrollPos++
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	if m.mode == ViewLoading {
		return m.renderSplash()
	}
	if m.isNarrow() {
		return m.renderNarrowLayout()
	}
	return m.renderWideLayout()
}

func (m Model) visibleLines() []string {
	return m.sections[m.selected].Lines
}
