package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO:
// - [ ] get path from arg
// - [ ] read from stdin
func read() any {
	f, err := os.ReadFile("test.json")
	if err != nil {
		log.Fatal(err)
	}

	var result any
	if err := json.Unmarshal(f, &result); err != nil {
		log.Fatal(err)
	}

	return result
}

func main() {
	p := tea.NewProgram(
		initialModel(read()),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
