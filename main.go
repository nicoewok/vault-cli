package main

import (
	"embed"
	"flag"
	"strings"
	"vault-cli/game" // Import your local package

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed data/*.txt
var dataFiles embed.FS

func main() {
	diff := flag.String("d", "easy", "difficulty level")
	flag.Parse()

	// Read the specific file based on flag
	content, _ := dataFiles.ReadFile("data/" + *diff + ".txt")
	words := strings.Split(string(content), "\n")

	p := tea.NewProgram(game.InitialModel(words))
	p.Run()
}
