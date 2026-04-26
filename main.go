// main.go
package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	diff := flag.String("d", "easy", "difficulty level (easy, medium, hard)")
	speed := flag.Int("s", 5, "rolling speed in milliseconds")
	flag.Parse()

	p := tea.NewProgram(initialModel(*diff, *speed))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
