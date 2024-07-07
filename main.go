package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vschmidt94/openapi-tui/lib/config"
	"github.com/vschmidt94/openapi-tui/tui/models"
	"os"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	p := tea.NewProgram(models.New(*cfg), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
