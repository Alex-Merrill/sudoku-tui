package menu

import (
	"github.com/Alex-Merrill/sudoku-tui/components/inputs"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
    rules string
    help help.Model
    keys inputs.KeyMap
}

func NewModel() Model {
    return Model {
        rules: "Normal sudoku rules apply. Each number 1-9 must be placed exactly once in each row, column, and box.",
        help: help.New(),
        keys: inputs.Controls,
    }
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.help.Width = msg.Width
    
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, inputs.Controls.Help):
            m.help.ShowAll = !m.help.ShowAll
        }
    }

    return m, nil
}

func (m Model) View() string {
    
    helpView := m.help.View(m.keys)

    return helpView
}

