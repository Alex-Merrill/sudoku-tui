package main

import (
	"fmt"
	"os"

	model "github.com/Alex-Merrill/sudoku-tui/components"

	tea "github.com/charmbracelet/bubbletea"
)

// sudoku generator library is broken for true medium difficulty
// instead using hard as medium and expert as hard
const (
	LEVEL_EASY   = 0
	LEVEL_MEDIUM = 2
	LEVEL_HARD   = 3
)

func main() {
	var mode int
	modeMap := map[string]int{
		"easy":   LEVEL_EASY,
		"medium": LEVEL_MEDIUM,
		"hard":   LEVEL_HARD,
	}

	// incorrect amount of args
	if len(os.Args) < 2 || len(os.Args) > 2 {
		fmt.Println(printArgHelp())
		os.Exit(0)
	} else { // handle mode checking
		if _, ok := modeMap[os.Args[1]]; !ok {
			fmt.Println(printArgHelp())
			os.Exit(0)
		}
	}

	mode = modeMap[os.Args[1]]

	p := tea.NewProgram(model.NewModel(mode), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		panic(err)
	}
}

func printArgHelp() string {
	return `go run main.go [mode]
               [mode] - [easy, medium, hard]`
}
