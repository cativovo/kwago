package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func logToFile(s string) {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := f.WriteString(s); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	filterTextInput textinput.Model
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return model{
		filterTextInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// TODO: add basic vim keybindings
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.filterTextInput, cmd = m.filterTextInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	sampleJson := `{"test": ["uno", "dos"]}`
	var input any
	json.Unmarshal([]byte(sampleJson), &input)

	p, _ := prettifyJson(input)

	return fmt.Sprintf(
		"Interactive JQ?\n\n%s\n\norig: %s\n\nresult:\n %s",
		m.filterTextInput.View(),
		p,
		processJson(input, m.filterTextInput.Value()),
	) + "\n"
}
