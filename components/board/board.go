package board

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Alex-Merrill/sudoku-tui/components/inputs"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	generator "github.com/forfuns/sudoku-go/generator"
)

// using int8 for our board as we don't need more space
// also our sudoku library gives us an array with int8's so less work xD
type Model struct {
    board [9][9]struct { // contains current game, answer key, and given cells
        game      int8
        answerKey int8
        given     bool
        pencils map[int8]bool
    }
    keyMap inputs.KeyMap // contains all inputs - uses bubbles/key to do fancy things for us
    currCell coordinate // current cell player is on
    selectedCells map[coordinate]bool // keeps track of all selected cells
    wrongCells map[coordinate]bool // cells which contain the wrong number, shown upon puzzle completion
    gameWon bool // are ya winnin' son?
    cellsLeft int // keep track of this so we know when to display error highlighting
}

type coordinate struct {
    row, col int
}

// Using shift+[1-9] for penciling, thus we need to map
// these characters to the proper number
var pencilMap = map[string]int8{
    "!": 1,
    "@": 2,
    "#": 3,
    "$": 4,
    "%": 5,
    "^": 6,
    "&": 7,
    "*": 8,
    "(": 9,
}



// Initializes board model
func NewModel(mode int) Model {
    // Generates sudoku board
    // Generate takes int 0-3 for easy, medium, hard, expert
    // medium is broken in the package I am using, and I can't
    // find a suitable library to replace it, so we are using
    // 0,2,3 for easy, medium, hard - this is defined in main.go
    sudoku,err := generator.Generate(mode)
    game, answerKey := sudoku.Puzzle(), sudoku.Answer()
   
    if err != nil {
        fmt.Println(err)
        os.Exit(0)
    }

    // populate board struct
    // game is state of sudoku
    // answerKey is solution
    // given marks given tiles, the user cannot change them
    var board [9][9]struct{
        game      int8
        answerKey int8
        given     bool
        pencils map[int8]bool
    }
    cellsLeft := 0
    for i := 0; i < 9; i++ {
        for j := 0; j < 9; j++ {
            board[i][j].game = game[(i*9)+j]
            board[i][j].answerKey = answerKey[(i*9)+j]
            board[i][j].given = game[(i*9)+j] != -1
            if given := game[(i*9)+j] != -1; given {
                board[i][j].given = given
            } else {
                cellsLeft++
            }
            board[i][j].pencils = make(map[int8]bool)
        }
    }

    startCell := coordinate{0, 0}
    selectedCells := make(map[coordinate]bool)
    selectedCells[startCell] = true

    return Model {
        board: board,
        keyMap: inputs.Controls,
        currCell: startCell,
        selectedCells: selectedCells,
        wrongCells: make(map[coordinate]bool),
        gameWon: false,
        cellsLeft: cellsLeft,
    }
}

func (m Model) Init() tea.Cmd {
    return nil
}

// TODO: add new keybindings in inputs.go for highlighting multiple cells
// as well as add new cases here for those keybindings, pointing to new cursorHighlight[Dir] funcs
// also look through changes and see if this is all that is needed for this change
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, inputs.Controls.Down):
            m.cursorDown()
        case key.Matches(msg, inputs.Controls.Up):
            m.cursorUp()

        case key.Matches(msg, inputs.Controls.Left):
            m.cursorLeft()

        case key.Matches(msg, inputs.Controls.Right):
            m.cursorRight()

        case key.Matches(msg, inputs.Controls.Number):
            num,_ := strconv.Atoi(msg.String())
            m.setCell(int8(num), m.currCell)

        case key.Matches(msg, inputs.Controls.PencilNumber):
            num := pencilMap[msg.String()]
            m.setPencilCell(int8(num), m.currCell)

        case key.Matches(msg, inputs.Controls.Delete):
            m.deleteCell(m.currCell)
        }
    }

    if m.cellsLeft == 0 {
        if m.checkWon() {
            m.gameWon = true
        }
    }

    return m, nil
}

func (m Model) View() string {
    
    if m.gameWon {
        return "You Won!"
    }

    // converts board.game cell to string for draw
    var convertToString = func(num int8) string {
        if num == -1 {
            return " "
        }
        return fmt.Sprintf("%d", num)
    }

    // check if there is any wrong cells and add text to top of board
    var err string
    if len(m.wrongCells) > 1 {
        err = "You need to fix " + strconv.Itoa(m.cellsLeft) + " cells!"
    } else if len(m.wrongCells) > 0{
        err = "You need to fix " + strconv.Itoa(m.cellsLeft) + " cells!"
    } else {
        err = ""
    }

    // iterates through board to add to draw string
    bLen := len(m.board)
    boardString := err + "\n\n"
    for i := 0; i < bLen; i++ {
        rowString := ""
        for j := 0; j < bLen; j++ {
            // TODO: make isSelected check if curr coordinate is in m.selectedCells
            _,cellWrong := m.wrongCells[coordinate{i,j}]
            isSelected := m.currCell.row == i && m.currCell.col == j
                       
            // add cell to row
            cell := drawCell(cellWrong, isSelected, m.board[i][j].given, convertToString(m.board[i][j].game), m.board[i][j].pencils)
            rowString = lipgloss.JoinHorizontal(lipgloss.Center, rowString, cell)
            // if we are at column where box border goes, add border
            if j == 2 || j == 5 {
                rowString = lipgloss.JoinHorizontal(lipgloss.Center, rowString, drawBorder("vert", ""))
            }
        } 

        // add row to board
        boardString = lipgloss.JoinVertical(lipgloss.Center, boardString, rowString) 

        // if we are at a row where box border goes, add border
        if i == 2 || i == 5 {
            boardString = lipgloss.JoinVertical(lipgloss.Center, boardString, drawBorder("hor", rowString)) 
        }    
    } 

    return boardString
}

// TODO: for all cursor funcs, clear selectedCells and add updated currCell to it
// TODO: also create new cursorHighlight[Dir] funcs which change currCell, and add updated currCell to selectedCells
func (m *Model) cursorDown() {
    if m.currCell.row < len(m.board) - 1  {
        m.currCell.row++ 
    } else {
        m.currCell.row = 0
    }
}

func (m *Model) cursorUp() {
    if m.currCell.row > 0 {
        m.currCell.row--
    } else {
        m.currCell.row = len(m.board) - 1
    }
}

func (m *Model) cursorLeft() {
    if m.currCell.col > 0 {
        m.currCell.col--
    } else {
        m.currCell.col = len(m.board[0]) - 1
    }
}

func (m *Model) cursorRight() {
    if m.currCell.col < len(m.board[0]) - 1 {
        m.currCell.col++
    } else {
        m.currCell.col = 0
    }
}

// TODO: make deleteCell, setCell, and setPencilCell funcs do event to all selectedCells instead of currCell
// clears cell at board[currCell.row, currCell.col].game
func (m *Model) deleteCell(currCell coordinate) {
    given := m.board[currCell.row][currCell.col].given
    if !given && m.board[currCell.row][currCell.col].game != -1 {
        m.board[currCell.row][currCell.col].game = -1
        delete(m.wrongCells, coordinate{currCell.row, currCell.col})
        m.cellsLeft++
    }
}

// sets cell at board[currCell.row][currCell.col].game
func (m *Model) setCell(num int8, currCell coordinate) {
    if !m.board[currCell.row][currCell.col].given {
        // if marking an empty cell or a wrong cell, decrement cellsLeft
        cellEmpty := m.board[currCell.row][currCell.col].game == -1
        _,cellWrong := m.wrongCells[coordinate{currCell.row, currCell.col}]
        if cellEmpty || cellWrong {
            m.cellsLeft--
        } 

        m.board[currCell.row][currCell.col].game = num
        delete(m.wrongCells, coordinate{currCell.row, currCell.col})
    }
}

// sets/removes pencil mark
func (m *Model) setPencilCell(num int8, currCell coordinate) {
    row := currCell.row
    col := currCell.col
    if !m.board[currCell.row][currCell.col].given {
        m.board[row][col].pencils[num] = !m.board[row][col].pencils[num]
    }
}

// need to get wrong cells and check for win seperately
// sudoku generator can output puzzles with multiple solutions
// its common for there to be x-wings at the end of the puzzle
// where some/all configurations work.
func (m *Model) checkWon() bool {
    // get wrong cells
    cellsWrong := 0        
    for i := 0; i < len(m.board); i++ {
        for j := 0; j < len(m.board[0]); j++ {
            if m.board[i][j].game != m.board[i][j].answerKey {
                m.wrongCells[coordinate{i,j}] = true
                cellsWrong++
            }
        }
    }

    // check for win
    won := m.checkForWinManual()

    m.cellsLeft = cellsWrong
    return won
}

// We can be a bit clever here and use row,col,box 2D arrays
// with type bool 
// We can then check for correctness with one pass through the
// board, as opposed to 3 (for every row, col, and box)
// ex: if a 1 is placed at row 2 col 2, we will set
// row[2][1] = true, col[2][1] = true, box[0][1] = true
func (m *Model) checkForWinManual() bool {
    // [9][10] because there are 9 cells in each row/col/box
    // but we have numbers 1-9, so [c][0] will never be used
    var row,col,box [9][10] bool

     
    for i := 0; i < len(m.board); i++ {
        for j := 0; j < len(m.board[0]); j++ {
            val := m.board[i][j].game
            if val != -1 {
                // row check
                if row[i][val] { // we have a dupe
                    return false
                }
                row[i][val] = true

                // col check
                if col[j][val] { // we have a dupe
                    return false
                }
                col[j][val] = true
                
                // box check
                /*
                    convert i,j to boxRow/boxCol indices 
                    ie: (0,0)│(0,1)│(0,2)
                        ─────┼─────┼─────
                        (1,0)│(1,1)│(1,2)
                        ─────┼─────┼─────
                        (2,0)│(2,1)│(2,2)
                    then convert to flat index boxIndx [0-8]
                */
                bRow := i/3
                bCol := j/3
                boxIdx := bRow*3 + bCol
                if box[boxIdx][val] { // we have a dupe
                    return false
                }
                box[boxIdx][val] = true

            } else { // i dont think we can ever get here, but better safe than sorry
                return false
            }
        }
    }

    // if we've made it here, there were no dupes, thus win
    return true
}

