package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func read() any {
	var file *os.File
	if len(os.Args) < 2 {
		file = os.Stdin
	} else {
		filepath := os.Args[1]

		var err error
		file, err = os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer file.Close()

	data := make([]byte, 0)

	reader := bufio.NewReader(file)

readByte:
	for {
		b, err := reader.ReadByte()
		if err != nil {
			switch err {
			case io.EOF:
				break readByte
			default:
				log.Fatal(err)
			}
		}

		data = append(data, b)
	}

	var result any
	if err := json.Unmarshal(data, &result); err != nil {
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
