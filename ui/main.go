package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Start() {
	configureLogging()

	p := tea.NewProgram(NewModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("an error occurred: %v", err)
		os.Exit(1)
	}
}

func configureLogging() {
	debugFlag, exists := os.LookupEnv("DEBUG")
	if !exists || debugFlag != "true" {
		return
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Printf("could not start logger: %v", err)
		os.Exit(1)
	}

	defer f.Close()
}
