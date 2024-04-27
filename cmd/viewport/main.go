package main

import (
	"log"

	viewport "github.com/cativovo/kwago/internal/sample"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(viewport.InitialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
