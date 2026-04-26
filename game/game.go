package game

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	state      string   // "opening", "hacking", "success", "lockout"
	words      []string // The list of potential passwords
	cursor     int      // Which word the user is pointing at
	secretIdx  int      // The index of the correct password
	attempts   int      // Remaining tries
	output     []string // History of "Likeness" feedback
	difficulty string
}

func initialModel(diff string) model {
	// Logic to pick words based on difficulty:
	// Easy: 5 letters, 8 words
	// Hard: 10 letters, 15 words
	return model{
		state:      "opening",
		words:      []string{"BATTERY", "STATION", "REACTOR", "PROTECT"}, // Example
		secretIdx:  2,
		attempts:   4,
		difficulty: diff,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.words)-1 {
				m.cursor++
			}

		case "enter":
			if m.state == "opening" {
				m.state = "hacking"
			} else if m.state == "hacking" {
				return m.checkGuess(), nil
			}
		}
	}
	return m, nil
}

func getLikeness(word, secret string) int {
	score := 0
	for i := range word {
		if word[i] == secret[i] {
			score++
		}
	}
	return score
}

func (m model) checkGuess() model {
	guess := m.words[m.cursor]
	secret := m.words[m.secretIdx]

	if guess == secret {
		m.state = "success"
		return m
	}

	// Calculate Likeness
	likeness := 0
	for i := 0; i < len(guess) && i < len(secret); i++ {
		if guess[i] == secret[i] {
			likeness++
		}
	}

	// Update game state
	m.attempts--
	feedback := fmt.Sprintf("> %s\n> Entry denied.\n> Likeness=%d", guess, likeness)
	m.output = append(m.output, feedback)

	if m.attempts <= 0 {
		m.state = "lockout"
	}

	return m
}
