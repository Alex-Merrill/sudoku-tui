package model

import (
	//"fmt"

	"github.com/Alex-Merrill/sudoku-tui/components/board"
	"github.com/Alex-Merrill/sudoku-tui/components/inputs"
	"github.com/Alex-Merrill/sudoku-tui/components/menu"
	"github.com/Alex-Merrill/sudoku-tui/components/winscreen"

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
    winscreen winscreen.Model
    gameWon bool
    winscreenDone bool

    width, height int
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var initCmd tea.Cmd
    // switch to check for quit command or window sizing
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {

        case key.Matches(msg, inputs.Controls.Quit):
            return m, tea.Quit

        case key.Matches(msg, inputs.Controls.NewGame):
            m.board = board.NewModel(m.mode)
            m.gameWon = false

        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

    case board.GameWon:
        m.gameWon = true
        m.winscreen = winscreen.NewModel()
        initCmd = m.winscreen.Init()

    /* if we allow user to stop animation, we can enable this
        case winscreen.StopAnim:
            m.winscreenDone = true
    */
    }

    // update board, menu, and winscreen models
    newBoardModel,boardCmd := m.board.Update(msg)
    newMenuModel,_ := m.menu.Update(msg)
    newWinScreenModel, winScreenCmd := m.winscreen.Update(msg)

    m.board = newBoardModel.(board.Model)
    m.menu = newMenuModel.(menu.Model)
    m.winscreen = newWinScreenModel.(winscreen.Model)

    return m, tea.Batch(boardCmd, winScreenCmd, initCmd)
}

func (m Model) View() string {
  
    // make composite view of app
    // board view on top, menu view on bottom

    if m.gameWon {
        compositeView := m.winscreen.View() +
                        "\n\n" +
                        "Press 'n' to start a new game" +
                        "\n" +
                        "Press 'q' or 'ctrl+c' to quit"
        return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, compositeView)

        /* if we allow user to stop animation, we can use this instead
            if m.winscreenDone {
                return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, compositeView)
            } else {
                return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.winscreen.View())
            }
        */
    }

    compositeView := m.board.View() + "\n\n" + m.menu.View() 

    return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, compositeView)
}

func NewModel(mode int) Model {
    return Model {
        mode: mode,
        board: board.NewModel(mode),
        menu: menu.NewModel(),
        winscreen: winscreen.NewModel(),
        gameWon: false,
        winscreenDone: false,
    }
}
