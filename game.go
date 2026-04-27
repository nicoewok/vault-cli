// game.go
package main

import (
	"embed"
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed data/*.txt
var dataFiles embed.FS

type wordLocation struct {
	text   string
	row    int
	startX int
	endX   int
	column int // 0 for left, 1 for right
}

type model struct {
	state     string
	leftGrid  []string
	rightGrid []string
	words     []wordLocation
	cursorX   int
	cursorY   int
	activeCol int
	secret    string
	attempts  int
	maxAttempts int
	output    []string
	rollIndex int
	rollSpeed time.Duration
}

func initialModel(diff string, speed int) model {
	content, err := dataFiles.ReadFile("data/" + diff + ".txt")
	if err != nil {
		content = []byte("ERROR\nVACUUM\nVOID\nNULL")
	}

	rawWords := strings.Split(strings.TrimSpace(string(content)), "\n")
	var allWords []string
	for _, w := range rawWords {
		if strings.TrimSpace(w) != "" {
			allWords = append(allWords, strings.TrimSpace(w))
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(allWords), func(i, j int) {
		allWords[i], allWords[j] = allWords[j], allWords[i]
	})

	// Change these values to higher targets
	desiredCount := 10
	attemptsCount := 4
	if diff == "medium" {
		desiredCount = 18
		attemptsCount = 6
	} else if diff == "hard" {
		desiredCount = 20
		attemptsCount = 8
	}

	actualCount := len(allWords)
	if actualCount > desiredCount {
		actualCount = desiredCount
	}

	selectedWords := allWords[:actualCount]
	secretWord := ""
	if actualCount > 0 {
		secretWord = selectedWords[r.Intn(actualCount)]
	}



	m := model{
		state:     "opening",
		secret:    secretWord,
		attempts:  attemptsCount,
		maxAttempts: attemptsCount,
		activeCol: 0,
		cursorX:   0,
		cursorY:   0,
		leftGrid:  make([]string, 15),
		rightGrid: make([]string, 15),
		rollIndex: 0,
		rollSpeed: time.Duration(speed) * time.Millisecond,
	}

	// Initialize grids with garbage
	for i := 0; i < 15; i++ {
		m.leftGrid[i] = randomGarbage(12, r)
		m.rightGrid[i] = randomGarbage(12, r)
	}

	// Place words randomly
	for _, w := range selectedWords {
		placed := false
		for attempts := 0; attempts < 100; attempts++ {
			column := r.Intn(2)
			row := r.Intn(15)
			startX := r.Intn(12 - len(w) + 1)

			// Check for overlap with existing words (including 1-char buffer)
			overlap := false
			for _, existing := range m.words {
				if existing.column == column && existing.row == row {
					// Check if [startX, startX+len(w)-1] overlaps with [existing.startX-1, existing.endX+1]
					if startX <= existing.endX+1 && startX+len(w)-1 >= existing.startX-1 {
						overlap = true
						break
					}
				}
			}

			if !overlap {
				// Place the word
				var grid []string
				if column == 0 {
					grid = m.leftGrid
				} else {
					grid = m.rightGrid
				}

				line := []rune(grid[row])
				for i, char := range w {
					line[startX+i] = char
				}
				grid[row] = string(line)

				m.words = append(m.words, wordLocation{
					text:   w,
					row:    row,
					startX: startX,
					endX:   startX + len(w) - 1,
					column: column,
				})
				placed = true
				break
			}
		}
		if !placed {
			// This could happen if the grid is very full, but with our settings it's unlikely.
			// We just skip the word if it doesn't fit after 100 tries.
		}
	}

	return m
}

func randomGarbage(width int, r *rand.Rand) string {
	symbols := "!@#$%^&*()[]{}<>/\\|:;,.?"
	s := ""
	for i := 0; i < width; i++ {
		s += string(symbols[r.Intn(len(symbols))])
	}
	return s
}

func (m model) Init() tea.Cmd {
	return tick(m.rollSpeed)
}

type tickMsg time.Time

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.state {
		case "opening":
			if msg.String() == "enter" {
				m.state = "hacking"
				m.rollIndex = 0 // Reset roll for the hacking view
			}
			return m, nil

		case "success", "lockout":
			return m, tea.Quit

		case "hacking":
			switch msg.String() {
			case "up":
				if m.cursorY > 0 {
					m.cursorY--
				}
			case "down":
				if m.cursorY < 14 {
					m.cursorY++
				}
			case "left":
				if m.cursorX > 0 {
					m.cursorX--
				} else if m.activeCol == 1 {
					m.activeCol = 0
					m.cursorX = 11
				}
			case "right":
				if m.cursorX < 11 {
					m.cursorX++
				} else if m.activeCol == 0 {
					m.activeCol = 1
					m.cursorX = 0
				}
			case "enter":
				for _, w := range m.words {
					if m.activeCol == w.column && m.cursorY == w.row &&
						m.cursorX >= w.startX && m.cursorX <= w.endX {
						return m.checkGuess(w.text), nil
					}
				}
			}
		}

	case tickMsg:
		m.rollIndex++
		return m, tick(m.rollSpeed)
	}
	return m, nil
}

func (m model) checkGuess(guess string) model {
	if guess == m.secret {
		m.state = "success"
		m.rollIndex = 0
		return m
	}

	guess_string := ""
	likeness := 0
	for i := 0; i < len(guess) && i < len(m.secret); i++ {
		if guess[i] == m.secret[i] {
			likeness++
			//draw in green
			guess_string += brightStyle.Render(string(guess[i]))
		} else {
			//draw in red
			guess_string += failedStyle.Render("?")
		}
	}

	m.attempts--
	m.output = append(m.output, fmt.Sprintf("> %s", guess))
	m.output = append(m.output, fmt.Sprintf("  %s", guess_string))

	if m.attempts <= 0 {
		m.state = "lockout"
		m.rollIndex = 0
	}

	return m
}
