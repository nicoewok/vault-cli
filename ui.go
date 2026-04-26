// ui.go
package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	terminalGreen = lipgloss.Color("#009940ff")
	darkGreen     = lipgloss.Color("#416345ff")
	darkRed       = lipgloss.Color("#881104ff")

	brightStyle = lipgloss.NewStyle().Foreground(terminalGreen)
	dimStyle    = lipgloss.NewStyle().Foreground(darkGreen)
	failedStyle = lipgloss.NewStyle().Foreground(darkRed)
	cursorStyle = lipgloss.NewStyle().
			Background(terminalGreen).
			Foreground(lipgloss.Color("#000000"))

	successStyle = lipgloss.NewStyle().Foreground(terminalGreen).Bold(true).Padding(1)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Padding(1)
)

const robcoLogo = `
  _____   ____  ____   _____  ____  
 |  __ \ / __ \|  _ \ / ____|/ __ \ 
 | |__) | |  | | |_) | |    | |  | |
 |  _  /| |  | |  _ <| |    | |  | |
 | | \ \| |__| | |_) | |____| |__| |
 |_|  \_\\____/|____/ \_____|\____/ 
      I  N  D  U  S  T  R  I  E  S`

func (m model) View() string {
	var output string
	switch m.state {
	case "lockout":
		output = errorStyle.Render(robcoLogo + "\n\n  [ TERMINAL LOCKED ]\n\n  PLEASE CONTACT AN ADMINISTRATOR\n\n  Press any key to exit...\n")
	case "success":
		output = successStyle.Render(robcoLogo + "\n\n  [ ACCESS GRANTED ]\n\n  WELCOME, OVERSEER\n\n  Press any key to exit...\n")
	case "opening":
		output = brightStyle.Render(robcoLogo + "\n\n  ROBCO INDUSTRIES (TM) TERMLINK \n\n  > Press [ENTER] to begin...\n")
	default:
		output = m.renderHackingView()
	}
	return m.applyRoll(output)
}

// applyRoll takes the full rendered string and only returns up to rollIndex characters,
// skipping ANSI escape codes in the count so styling doesn't break.
func (m model) applyRoll(s string) string {
	var out strings.Builder
	count := 0
	inAnsi := false

	for _, r := range s {
		if r == '\x1b' {
			inAnsi = true
		}

		if inAnsi {
			out.WriteRune(r)
			if r == 'm' { // End of ANSI color code
				inAnsi = false
			}
			continue
		}

		if count < m.rollIndex {
			out.WriteRune(r)
			count++
		} else {
			break
		}
	}
	return out.String()
}

func (m model) renderHackingView() string {
	startAddr := 0xF82C
	leftView := m.renderGrid(m.leftGrid, 0, startAddr)
	rightView := m.renderGrid(m.rightGrid, 1, startAddr+256)

	gameGrid := lipgloss.JoinHorizontal(lipgloss.Top, leftView, "    ", rightView)
	history := lipgloss.NewStyle().MarginLeft(4).Render(strings.Join(m.output, "\n"))

	return lipgloss.JoinHorizontal(lipgloss.Top, gameGrid, history)
}

func (m model) renderGrid(grid []string, colIdx int, startAddr int) string {
	var out strings.Builder
	for y, line := range grid {
		// Add the Hex address
		out.WriteString(dimStyle.Render(fmt.Sprintf("[0x%X]  ", startAddr+(y*12))))

		for x, char := range line {
			charStr := string(char)
			// Highlight ONLY the specific character the cursor is over
			if m.activeCol == colIdx && m.cursorX == x && m.cursorY == y {
				out.WriteString(cursorStyle.Render(charStr))
			} else {
				out.WriteString(brightStyle.Render(charStr))
			}
		}
		out.WriteString("\n")
	}
	return out.String()
}
