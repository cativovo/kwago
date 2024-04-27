run:
	@go run .

install_deps:
	go get github.com/charmbracelet/bubbles
	go get github.com/charmbracelet/bubbletea
	github.com/charmbracelet/lipgloss

install_dev_deps: install_deps
	go install github.com/go-delve/delve/cmd/dlv@latest

debug_start:
	dlv debug --headless --api-version=2 --listen=127.0.0.1:43000 .

debug_connect:
	dlv connect 127.0.0.1:43000
