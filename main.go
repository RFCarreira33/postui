package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcarreira33/postui/ui"
)

func main() {
	program := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
