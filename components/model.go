package model

import (
	"github.com/Alex-Merrill/sudoku-tui/components/board"
	"github.com/Alex-Merrill/sudoku-tui/components/inputs"
	"github.com/Alex-Merrill/sudoku-tui/components/menu"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Defines App Model
// Consists of board component and menu component
type Model struct {
    mode int
    board board.Model
    menu menu.Model

    width, height int
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // switch to check for quit command or window sizing
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, inputs.Controls.Quit):
            return m, tea.Quit
        case key.Matches(msg, inputs.Controls.NewGame):
            m.board = board.NewModel(m.mode)
        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

    }   

    // update board and menu models
    newBoardModel,_ := m.board.Update(msg)
    newMenuModel,_ := m.menu.Update(msg)
    
    m.board = newBoardModel.(board.Model)
    m.menu = newMenuModel.(menu.Model)

    return m, nil
}

func (m Model) View() string {
  
    // make composite view of app
    // board view on top, menu view on bottom
    compositeView := lipgloss.JoinVertical(lipgloss.Center, m.board.View(), "\n", m.menu.View())

    return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, compositeView)
}

func NewModel(mode int) Model {
    return Model {
        mode: mode,
        board: board.NewModel(mode),
        menu: menu.NewModel(),
    }
}
