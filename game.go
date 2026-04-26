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

	desiredCount := 10
	if diff == "hard" {
		desiredCount = 16
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
		attempts:  4,
		activeCol: 0,
		cursorX:   0,
		cursorY:   0,
		leftGrid:  make([]string, 15),
		rightGrid: make([]string, 15),
		rollIndex: 0,
		rollSpeed: time.Duration(speed) * time.Millisecond,
	}

	// Build the grids and store word locations
	numWords := len(selectedWords)
	half := (numWords + 1) / 2
	if half == 0 && numWords > 0 {
		half = 1
	}

	for i, w := range selectedWords {
		column := 0
		row := i
		if half > 0 {
			column = i / half
			row = i % half
		}

		if row >= 15 {
			continue
		}

		line, startX := generateScrambledLine(w, 12, r)

		if column == 0 {
			m.leftGrid[row] = line
		} else if column == 1 {
			m.rightGrid[row] = line
		}

		m.words = append(m.words, wordLocation{
			text:   w,
			row:    row,
			startX: startX,
			endX:   startX + len(w) - 1,
			column: column,
		})
	}

	for i := 0; i < 15; i++ {
		if m.leftGrid[i] == "" {
			m.leftGrid[i] = randomGarbage(12, r)
		}
		if m.rightGrid[i] == "" {
			m.rightGrid[i] = randomGarbage(12, r)
		}
	}

	return m
}

func generateScrambledLine(word string, width int, r *rand.Rand) (string, int) {
	symbols := "!@#$%^&*()[]{}<>/\\|:;,.?"
	startX := 0
	if width > len(word) {
		startX = r.Intn(width - len(word) + 1)
	}
	line := ""
	for i := 0; i < width; i++ {
		if i >= startX && i < startX+len(word) {
			line += string(word[i-startX])
		} else {
			line += string(symbols[r.Intn(len(symbols))])
		}
	}
	return line, startX
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
			guess_string += string(guess[i])
		} else {
			guess_string += "?"
		}
	}

	m.attempts--
	m.output = append(m.output, fmt.Sprintf("> %s -> %s", guess, guess_string))
	m.output = append(m.output, "> Entry denied.")
	m.output = append(m.output, fmt.Sprintf("> Likeness=%d. Remaining attempts: %d", likeness, m.attempts))

	if m.attempts <= 0 {
		m.state = "lockout"
		m.rollIndex = 0
	}

	return m
}
