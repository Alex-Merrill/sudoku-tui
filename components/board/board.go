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

        case key.Matches(msg, inputs.Controls.ShiftDown):
            m.cursorHighlightDown()

        case key.Matches(msg, inputs.Controls.ShiftUp):
            m.cursorHighlightUp()

        case key.Matches(msg, inputs.Controls.ShiftLeft):
            m.cursorHighlightLeft()

        case key.Matches(msg, inputs.Controls.ShiftRight):
            m.cursorHighlightRight()

        case key.Matches(msg, inputs.Controls.Number):
            num,_ := strconv.Atoi(msg.String())
            m.setCell(int8(num))

        case key.Matches(msg, inputs.Controls.PencilNumber):
            num := pencilMap[msg.String()]
            m.setPencilCell(int8(num))

        case key.Matches(msg, inputs.Controls.Delete):
            m.deleteCell()
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
            _,cellWrong := m.wrongCells[coordinate{i,j}]
            _,isSelected := m.selectedCells[coordinate{i,j}]
            isCurrCell := m.currCell.row == i && m.currCell.col == j

            // add cell to row
            cell := drawCell(cellWrong, isSelected, isCurrCell, m.board[i][j].given, convertToString(m.board[i][j].game), m.board[i][j].pencils)
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

// sets cell at all selected cells
func (m *Model) setCell(num int8) {
    for k := range m.selectedCells {
        row := k.row
        col := k.col
        if !m.board[row][col].given {
            // if marking an empty cell or a wrong cell, decrement cellsLeft
            cellEmpty := m.board[row][col].game == -1
            _,cellWrong := m.wrongCells[coordinate{row, col}]
            if cellEmpty || cellWrong {
                m.cellsLeft--
            } 

            m.board[row][col].game = num
            delete(m.wrongCells, coordinate{row, col})
            m.updatePencilCells(num, coordinate{row, col})
        }
    }
}

// clears all selected cells
func (m *Model) deleteCell() {
    for k := range m.selectedCells {
        row := k.row
        col := k.col
        given := m.board[row][col].given
        if !given && m.board[row][col].game != -1 { // delete cell value
            m.board[row][col].game = -1
            delete(m.wrongCells, coordinate{row, col})
            m.cellsLeft++
        } else if !given { // delete pencil marks if no cell value
            for i := 1; i < 10; i++ {
                m.board[row][col].pencils[int8(i)] = false
            }
        }
    }
}

// sets/removes pencil mark at all selected cells if cell is not given or value is not set
func (m *Model) setPencilCell(num int8) {
    for k := range m.selectedCells {
        row := k.row
        col := k.col
        given := m.board[row][col].given
        set := m.board[row][col].game != -1
        if !given && !set {
            m.board[row][col].pencils[num] = !m.board[row][col].pencils[num]
        }
    }
}

// updates pencil cells in given row/box/col based on new "num" in cell currCell
func (m *Model) updatePencilCells(num int8, currCell coordinate) {
    row := currCell.row
    col := currCell.col
    // row
    for i := 0; i < 9; i++ {
        coord := coordinate{row, i}
        m.removePencilCell(num, coord)
    } 

    // col
    for i := 0; i < 9; i++ {
        coord := coordinate{i, col}
        m.removePencilCell(num, coord)
    } 

    // box
    // convert (row,col) to the box index (0,0 0,1 0,2 for top three boxes etc)
    // then convert those coordinates to the actual coords of the top left cell in box
    bRow := row/3
    bCol := col/3
    realRow := bRow * 3
    realCol := bCol * 3
    for i := realRow; i < realRow + 3; i++ {
        for j := realCol; j < realCol + 3; j++ {
            coord := coordinate{i, j}
            m.removePencilCell(num, coord)
        }
    }
}

// takes a number 1-9 and a coordinate
// if there is a pencil mark of number "num" at coordiante currCell
// and the cell is not given and a value is not set, removes pencil mark
func (m *Model) removePencilCell(num int8, currCell coordinate) {
    row := currCell.row
    col := currCell.col
    given := m.board[row][col].given
    pencilsContainsNum := m.board[row][col].pencils[num]
    set := m.board[row][col].game != -1
    if !given && !set && pencilsContainsNum {
        m.board[row][col].pencils[num] = false
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
