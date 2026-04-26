package game

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	green    = lipgloss.Color("#00FF41")
	dimGreen = lipgloss.Color("#003B00")
	style    = lipgloss.NewStyle().Foreground(green).Bold(true)
	cursor   = lipgloss.NewStyle().Background(green).Foreground(lipgloss.Color("#000000"))
)

func (m model) View() string {

	if m.state == "lockout" {
		return "\n\n  TERMINAL LOCKED\n\n  PLEASE CONTACT AN ADMINISTRATOR\n\n"
	} else if m.state == "success" {
		return "\n\n  ACCESS GRANTED\n\n  WELCOME, OVERSEER\n\n"
	} else if m.state == "opening" {
		return style.Render("ROBCO INDUSTRIES (TM) TERMLINK \n\n Press [ENTER] to begin...")
	}

	s := "ENTER PASSWORD NOW\n\n"
	for i, word := range m.words {
		if m.cursor == i {
			s += cursor.Render("> "+word) + "\n"
		} else {
			s += style.Render("  "+word) + "\n"
		}
	}

	s += fmt.Sprintf("\nAttempts remaining: %d", m.attempts)
	return s
}
