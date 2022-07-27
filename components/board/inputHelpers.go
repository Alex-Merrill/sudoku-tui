package board

// Handles all cursor movements

// all cursor[Dir] funcs move currCell, and delete all entries in selectedCells
func (m *Model) cursorDown() {
    m.selectedCells = make(map[coordinate]bool)
    if m.currCell.row < len(m.currBoardState.board) - 1  {
        m.currCell.row++ 
    } else {
        m.currCell.row = 0
    }
    m.selectedCells[m.currCell] = true
}

func (m *Model) cursorUp() {
    m.selectedCells = make(map[coordinate]bool)
    if m.currCell.row > 0 {
        m.currCell.row--
    } else {
        m.currCell.row = len(m.currBoardState.board) - 1
    }
    m.selectedCells[m.currCell] = true
}

func (m *Model) cursorLeft() {
    m.selectedCells = make(map[coordinate]bool)
    if m.currCell.col > 0 {
        m.currCell.col--
    } else {
        m.currCell.col = len(m.currBoardState.board[0]) - 1
    }
    m.selectedCells[m.currCell] = true
}

func (m *Model) cursorRight() {
    m.selectedCells = make(map[coordinate]bool)
    if m.currCell.col < len(m.currBoardState.board[0]) - 1 {
        m.currCell.col++
    } else {
        m.currCell.col = 0
    }
    m.selectedCells[m.currCell] = true
}

// all cursorHighlight[Dir] funcs move currCell in correct direction
// as well as add the new currCell into selected Cells
func (m *Model) cursorHighlightDown() {
    if m.currCell.row < len(m.currBoardState.board) - 1  {
        m.currCell.row++ 
    } else {
        m.currCell.row = 0
    }
    m.selectedCells[m.currCell] = true
}

func (m *Model) cursorHighlightUp() {
    if m.currCell.row > 0 {
        m.currCell.row--
    } else {
        m.currCell.row = len(m.currBoardState.board) - 1
    }
    m.selectedCells[m.currCell] = true
}

func (m *Model) cursorHighlightLeft() {
    if m.currCell.col > 0 {
        m.currCell.col--
    } else {
        m.currCell.col = len(m.currBoardState.board[0]) - 1
    }
    m.selectedCells[m.currCell] = true
}

func (m *Model) cursorHighlightRight() {
    if m.currCell.col < len(m.currBoardState.board[0]) - 1 {
        m.currCell.col++
    } else {
        m.currCell.col = 0
    }
    m.selectedCells[m.currCell] = true
}
