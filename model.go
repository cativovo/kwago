package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type model struct {
	filterTextInput textinput.Model
	json            any
	ready           bool
	viewport        viewport.Model
}

func (m model) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func initialModel(j any) model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 999
	ti.Width = 100
	return model{
		filterTextInput: ti,
		json:            j,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		inputCmd    tea.Cmd
		viewPortCmd tea.Cmd
		cmds        []tea.Cmd
	)
	m.filterTextInput, inputCmd = m.filterTextInput.Update(msg)
	m.viewport, viewPortCmd = m.viewport.Update(msg)
	cmds = append(cmds, inputCmd, viewPortCmd)

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// TODO: add basic vim keybindings
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			m.viewport.SetContent(processJson(m.json, m.filterTextInput.Value()))
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight + 10

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			b, _ := prettifyJson(m.json)
			m.viewport.SetContent(string(b))
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	viewPortComponent := fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	inputComponent := fmt.Sprintf(
		"Interactive JQ?\n\n%s",
		m.filterTextInput.View(),
	)

	return fmt.Sprintf("%s\n%s", viewPortComponent, inputComponent)
}
