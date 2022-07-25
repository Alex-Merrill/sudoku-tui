package board

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
    GIVEN_BASE_COLOR = lipgloss.Color("#042B4E")
    GIVEN_SELECTED_COLOR = lipgloss.Color("#0164A2")
    GIVEN_CURRENT_CELL_COLOR = lipgloss.Color("#318DCA")
    NOT_GIVEN_BASE_COLOR = lipgloss.Color("#68C5DB")
    NOT_GIVEN_SELECTED_COLOR = lipgloss.Color("#0197F6")
    NOT_GIVEN_CURRENT_CELL_COLOR = lipgloss.Color("#1677AD")
    WRONG_BASE_COLOR = lipgloss.Color("#D7263D")
    WRONG_SELECTED_COLOR = lipgloss.Color("#E77484")
    WRONG_CURRENT_CELL_COLOR = lipgloss.Color("#E15164")
    BOLD_BORDER_COLOR = lipgloss.Color("#F26419")
    PENCIL_MARK_COLOR = lipgloss.Color("#F58347")
    FINAL_VALUE_COLOR = lipgloss.Color("#ffffff")
)

// style definitions and ui rendering
var (
   
    /* 
        draws full cell, which is a 3x3 grid of 1 character cells.
        this allows us to put pencil markings in each cell of the grid.
        there are three different cell types, a highlighted cell, a normal cell,
        and an error cell. Each type has two states:
        
        error cell: highlighted, non-highlighted
        normal cell: given, modifiable
        highlighted cell: given, modifiable

        Each of these states has a different color, this is a generic function
        that takes the primaryColor boolean, which determines which state the
        given cell is, and the primary color and secondary color associated
        with those given states.
    */
    drawFullCell = func(cellColor lipgloss.Color, cell string, pencils map[int8]bool) string {
        cellString := ""
        for i := 0; i < 3; i++ {
            currRow := ""
            for j := 0; j < 3; j++ {
                // checks whether to render pencil marks or cell value
                var valToRender string
                var foregroundColor lipgloss.Color
                if cell == " " { // cell not marked render, pencil marks
                    if pencils[int8(i*3 + j + 1)] {
                        valToRender = fmt.Sprintf("%d", i*3 + j + 1)
                    } else {
                        valToRender = " "
                    }
                    foregroundColor = PENCIL_MARK_COLOR
                } else { // cell marked, dont render pencil marks, only render cell val on middle cell
                    if i == 1 && j == 1 {
                        valToRender = cell
                    } else {
                        valToRender = " "
                    }
                    foregroundColor = FINAL_VALUE_COLOR
                }
                
                // creates cell string
                currRow += lipgloss.NewStyle().
                                    Padding(0, 1, 0, 1).
                                    Foreground(foregroundColor).
                                    Background(cellColor).
                                    Render(valToRender)
            }
            // only add a new line on the first two rows of the cell
            if i < 2 {
                cellString += currRow + "\n"
            } else {
                cellString += currRow
            }
        }
        return cellString
    }

    // renders cell
    drawCell = func(cellWrong, isSelected, isCurrCell, given bool, cell string, pencils map[int8]bool) string {
        if given { // given cells
            if isCurrCell { // cell is current cell
                return drawFullCell(GIVEN_CURRENT_CELL_COLOR, cell, pencils)
            } else if isSelected { // cell is a selected cell and not the current cell
                return drawFullCell(GIVEN_SELECTED_COLOR, cell, pencils)
            } else { // cell is a base cell 
                return drawFullCell(GIVEN_BASE_COLOR, cell, pencils)
            }
        } else if cellWrong { // error cells
            if isCurrCell { // cell is current cell
                return drawFullCell(WRONG_CURRENT_CELL_COLOR, cell, pencils)
            } else if isSelected { // cell is a selected cell and not the current cell
                return drawFullCell(WRONG_SELECTED_COLOR, cell, pencils)
            } else { // cell is a base cell 
                return drawFullCell(WRONG_BASE_COLOR, cell, pencils)
            }
        } else { // modifiable cells
            if isCurrCell { // cell is current cell
                return drawFullCell(NOT_GIVEN_CURRENT_CELL_COLOR, cell, pencils)
            } else if isSelected { // cell is a selected cell and not the current cell
                return drawFullCell(NOT_GIVEN_SELECTED_COLOR, cell, pencils)
            } else { // cell is a base cell 
                return drawFullCell(NOT_GIVEN_BASE_COLOR, cell, pencils)
            }
        }
    }

    // takes string with direction(vert, hor) and a rowString, rowString only needed
    // for horizontal border
    drawBorder = func(dir string, rowString string) string {
        if dir == "vert" {
            return drawVerticalBorder()
        } else {
            return drawHorizontalBorder(rowString)
        }
    }

    // returns vertical border string for one cell
    drawVerticalBorder = func() string {
        border := lipgloss.NewStyle().
                           Padding(0, 1, 0, 1).
                           Foreground(BOLD_BORDER_COLOR).
                           Render("│")

        return lipgloss.JoinVertical(lipgloss.Center, border, border, border)
    }
        
    // returns horizontal border string for one row
    drawHorizontalBorder = func(rowString string) string {
        rowWidth,_ := lipgloss.Size(rowString)
        renderChar := "─"
        /* 
            the middle box border is one longer than the outside box borders
            since the middle box border has to meet the joint border on both sides,
            whereas the outside box borders have to meet the joint border only on
            one side.(each cell, including the borders, have a width of 3,
            as there are two padded cells on either side of the char)
            Thus, the outside borders have a length of boxWidth + 1  and the inside
            border has a length of boxWidth + 2
            rowWidth/3 - 1 = boxWidth + 1, as there are six cells that are not on the
            box (the border cells)
            ex:
            Each box is 27 cells wide, with 3 cells in between the boxes.
            Thus, the entire board is 87 cells wide. Outside box borders are 28 cells wide
            and the middle box border must be 29 cells wide.
        */
        renderChar = strings.Repeat(renderChar, rowWidth/3 - 1)
        middleBoxBorder := lipgloss.NewStyle().
                           Padding(0, 0, 0, 0).
                           Foreground(BOLD_BORDER_COLOR).
                           Render(renderChar + "─")
        outsideBoxesBorders := lipgloss.NewStyle().
                           Padding(0, 0, 0, 0).
                           Foreground(BOLD_BORDER_COLOR).
                           Render(renderChar)
        borderJoint := lipgloss.NewStyle().
                           Padding(0, 0, 0, 0).
                           Foreground(BOLD_BORDER_COLOR).
                           Render("┼")

        return lipgloss.JoinHorizontal(lipgloss.Left,
                                       outsideBoxesBorders, borderJoint,
                                       middleBoxBorder, borderJoint,
                                       outsideBoxesBorders)
    }

)
