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

// we make currBoardState a pointer to the current BoardState,
// this allows us to modify currBoardState instead of m.boardStates[m.currBoardStateIdx]
// makes our code a little cleaner
type Model struct {
    boardStates []BoardState
    currBoardStateIdx int
    currBoardState *BoardState
    keyMap inputs.KeyMap // contains all inputs - uses bubbles/key to do fancy things for us
    currCell coordinate // current cell player is on
    selectedCells map[coordinate]bool // keeps track of all selected cells
}

// using int8 for our board as our sudoku library gives us an array with int8's so less work xD
type BoardState struct {
    board [9][9]struct { // contains current game, answer key, and given cells
        game      int8
        answerKey int8
        given     bool
        pencils map[int8]bool
    }
    wrongCells map[coordinate]bool // cells which contain the wrong number, shown upon puzzle completion
    gameWon bool // are ya winnin' son?
    cellsLeft int // keep track of this so we know when to display error highlighting
}

type coordinate struct {
    row, col int
}

/*
    BoardState method that will copy a board for the next board state
    we need to initialize new maps for b.board.pencils and b.wrongCells
    and copy the old maps to the new ones - we do this as to maintain different
    wrongCells and different pencils states in each BoardState
*/
func (b BoardState) copyBoard() BoardState {
    // make new map for wrongCells
    newWrongCells := make(map[coordinate]bool)
    for k,v := range b.wrongCells {
        newWrongCells[k] = v
    }
    b.wrongCells = newWrongCells

    // make new map for pencils in each cell
    for i := 0; i < 9; i++ {
        for j := 0; j < 9; j++ {
            newPencils := make(map[int8]bool)
            for k,v := range b.board[i][j].pencils {
                newPencils[k] = v
            }
            b.board[i][j].pencils = newPencils
        }
    } 

    return b
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
    /*
        Generates sudoku board
        Generate takes int 0-3 for easy, medium, hard, expert
        medium is broken in the package I am using, and I can't
        find a suitable library to replace it, so we are using
        0,2,3 for easy, medium, hard - this is defined in main.go
    */
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

    // make top left cell starting cell for cursor
    startCell := coordinate{0, 0}
    selectedCells := make(map[coordinate]bool)
    selectedCells[startCell] = true

    // initialize the first board state
    boardState := BoardState {
        board: board,
        wrongCells: make(map[coordinate]bool),
        gameWon: false,
        cellsLeft: cellsLeft,
    }

    return Model{
        boardStates: []BoardState{ boardState },
        currBoardStateIdx: 0,
        currBoardState: &boardState,
        keyMap: inputs.Controls,
        currCell: startCell,
        selectedCells: selectedCells,
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

        case key.Matches(msg, inputs.Controls.Undo):
           m.UndoBoardAction() 

        case key.Matches(msg, inputs.Controls.Redo):
           m.RedoBoardAction() 

        }
    }

    // check if game is won when 0 zeros cell to fill
    if m.currBoardState.cellsLeft == 0 {
        if m.checkWon() {
            m.currBoardState.gameWon = true
        }
    }

    return m, nil
}

func (m Model) View() string {
    
    if m.currBoardState.gameWon {
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
    if len(m.currBoardState.wrongCells) > 1 {
        err = "You need to fix " + strconv.Itoa(m.currBoardState.cellsLeft) + " cells!"
    } else if len(m.currBoardState.wrongCells) > 0{
        err = "You need to fix " + strconv.Itoa(m.currBoardState.cellsLeft) + " cells!"
    } else {
        err = ""
    }

    // iterates through board to add to draw string
    bLen := len(m.currBoardState.board)
    boardString := err + "\n\n"
    for i := 0; i < bLen; i++ {
        rowString := ""
        for j := 0; j < bLen; j++ {
            _,cellWrong := m.currBoardState.wrongCells[coordinate{i,j}]
            _,isSelected := m.selectedCells[coordinate{i,j}]
            isCurrCell := m.currCell.row == i && m.currCell.col == j

            // add cell to row
            cell := drawCell(cellWrong, isSelected, isCurrCell, m.currBoardState.board[i][j].given, convertToString(m.currBoardState.board[i][j].game), m.currBoardState.board[i][j].pencils)
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
    // check if we need to make a new board state
    somethingChanged := false
    for k := range m.selectedCells {
        row := k.row
        col := k.col
        if !m.currBoardState.board[row][col].given {
            somethingChanged = true
        }
    }

    // if we are supposed to do something, make new board state, then do it
    // to the new board state, else return
    if somethingChanged {
        m.makeNewBoardState()
    } else {
        return
    }

    for k := range m.selectedCells {
        row := k.row
        col := k.col
        if !m.currBoardState.board[row][col].given {
            // if marking an empty cell or a wrong cell, decrement cellsLeft
            cellEmpty := m.currBoardState.board[row][col].game == -1
            _,cellWrong := m.currBoardState.wrongCells[coordinate{row, col}]
            if cellEmpty || cellWrong {
                m.currBoardState.cellsLeft--
            } 

            m.currBoardState.board[row][col].game = num
            delete(m.currBoardState.wrongCells, coordinate{row, col})
            m.updatePencilCells(num, coordinate{row, col})
        }
    }
}

// clears all selected cells
func (m *Model) deleteCell() {
    // check if we need to make a new board state
    somethingChanged := false
    for k := range m.selectedCells {
        row := k.row
        col := k.col
        given := m.currBoardState.board[row][col].given
        if !given && m.currBoardState.board[row][col].game != -1 { // delete cell value
            somethingChanged = true
        } else if !given { // delete pencil marks if no cell value
            for i := 1; i < 10; i++ {
                if m.currBoardState.board[row][col].pencils[int8(i)] {
                    somethingChanged = true
                }
            }
        }
    }

    // if we are supposed to do something, make new board state, then do it
    // to the new board state, else return
    if somethingChanged {
        m.makeNewBoardState()
    } else {
        return
    }

    for k := range m.selectedCells {
        row := k.row
        col := k.col
        given := m.currBoardState.board[row][col].given
        if !given && m.currBoardState.board[row][col].game != -1 { // delete cell value
            m.currBoardState.board[row][col].game = -1
            delete(m.currBoardState.wrongCells, coordinate{row, col})
            m.currBoardState.cellsLeft++
        } else if !given { // delete pencil marks if no cell value
            for i := 1; i < 10; i++ {
                m.currBoardState.board[row][col].pencils[int8(i)] = false
            }
        }
    }
}

// sets/removes pencil mark at all selected cells if cell is not given or value is not set
func (m *Model) setPencilCell(num int8) {
    // check if we need to make a new board state
    somethingChanged := false
    for k := range m.selectedCells {
        row := k.row
        col := k.col
        given := m.currBoardState.board[row][col].given
        set := m.currBoardState.board[row][col].game != -1
        if !given && !set {
            somethingChanged = true
        }
    }

    // if we are supposed to do something, make new board state, then do it
    // to the new board state, else return
    if somethingChanged {
        m.makeNewBoardState()
    } else {
        return
    }

    for k := range m.selectedCells {
        row := k.row
        col := k.col
        given := m.currBoardState.board[row][col].given
        set := m.currBoardState.board[row][col].game != -1
        if !given && !set {
            m.currBoardState.board[row][col].pencils[num] = !m.currBoardState.board[row][col].pencils[num]
        }
    }
}

// updates pencil cells in given row/box/col based on new "num" in cell currCell
// we only call this function from setCell
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

/*
    takes a number 1-9 and a coordinate
    if there is a pencil mark of number "num" at coordiante currCell
    and the cell is not given and a value is not set, removes pencil mark
    this is a helper function for updatePencilCells
*/
func (m *Model) removePencilCell(num int8, currCell coordinate) {
    row := currCell.row
    col := currCell.col
    given := m.currBoardState.board[row][col].given
    pencilsContainsNum := m.currBoardState.board[row][col].pencils[num]
    set := m.currBoardState.board[row][col].game != -1
    if !given && !set && pencilsContainsNum {
        m.currBoardState.board[row][col].pencils[num] = false
    }
}

// sets current board state to last board state
func (m *Model) UndoBoardAction() {
    if m.currBoardStateIdx > 0 {
        m.currBoardStateIdx--
        m.currBoardState = &m.boardStates[m.currBoardStateIdx]
    }
}

// sets current board state to next board state
func (m *Model) RedoBoardAction() {
    if m.currBoardStateIdx < len(m.boardStates) - 1 {
        m.currBoardStateIdx++
        m.currBoardState = &m.boardStates[m.currBoardStateIdx]
    }
}

/*
    makes new board state to be called before action

    If our current board state is pointing to the most recent
    board state, then we can just add a new board state and point our
    current board state to it
    However, if the current board state is not pointing to the latest board state,
    then we delete all board states that come after the current one, after which 
    we can append a new board state and point our current state to it
*/
func (m *Model) makeNewBoardState() {
    newBoardState := m.currBoardState.copyBoard()
    if m.currBoardState != &m.boardStates[len(m.boardStates)-1] {
        currIdx := m.currBoardStateIdx
        m.boardStates = m.boardStates[:currIdx + 1]
    }
    m.boardStates = append(m.boardStates, newBoardState)
    m.currBoardStateIdx++
    m.currBoardState = &m.boardStates[m.currBoardStateIdx]
}

/*
    need to get wrong cells and check for win seperately
    sudoku generator can output puzzles with multiple solutions
    its common for there to be x-wings at the end of the puzzle
    where some/all configurations work. Thus, our answer key might be different
    than the board while the board is still a valid solution
*/
func (m *Model) checkWon() bool {
    // get wrong cells
    cellsWrong := 0        
    for i := 0; i < len(m.currBoardState.board); i++ {
        for j := 0; j < len(m.currBoardState.board[0]); j++ {
            if m.currBoardState.board[i][j].game != m.currBoardState.board[i][j].answerKey {
                m.currBoardState.wrongCells[coordinate{i,j}] = true
                cellsWrong++
            }
        }
    }

    // check for win
    won := m.checkForWinManual()

    m.currBoardState.cellsLeft = cellsWrong
    return won
}

/*
    We can be a bit clever here and use row,col,box 2D arrays
    with type bool 
    We can then check for correctness with one pass through the
    board, as opposed to 3 (for every row, col, and box)
    ex: if a 1 is placed at row 2 col 2, we will set
    row[2][1] = true, col[2][1] = true, box[0][1] = true
*/
func (m *Model) checkForWinManual() bool {
    // [9][10] because there are 9 cells in each row/col/box
    // but we have numbers 1-9, so [c][0] will never be used
    var row,col,box [9][10] bool

    sudokuLen := len(m.currBoardState.board)
    for i := 0; i < sudokuLen; i++ {
        for j := 0; j < sudokuLen; j++ {
            val := m.currBoardState.board[i][j].game
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
