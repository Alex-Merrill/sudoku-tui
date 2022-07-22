package board

import (

	"github.com/charmbracelet/lipgloss"
)


const (
    BASE_GIVEN_COLOR = lipgloss.Color("#02182B")
    BASE_NOT_GIVEN_COLOR = lipgloss.Color("#68C5DB")
    SELECTED_GIVEN_COLOR = lipgloss.Color("#0164A2")
    SELECTED_NOT_GIVEN_COLOR = lipgloss.Color("#0197F6")
    WRONG_SELECTED_COLOR = lipgloss.Color("#E77484")
    WRONG_NOT_SELECTED_COLOR = lipgloss.Color("#D7263D")
    BOLD_BORDER_COLOR = lipgloss.Color("#F26419")
    NORMAL_BORDER_COLOR = lipgloss.Color("#ffffff")
)

// style definitions and ui rendering
var (
   
    // draws full cell, which is a 3x3 grid of 1 character cells
    // this allows us to put pencil markings in each cell of the grid
    drawFullCell = func(primaryColor bool, primaryGlossColor, secondaryGlossColor lipgloss.Color, cell string) string {
        cellString := ""
        if primaryColor {
            for i := 0; i < 3; i++ {
                currRow := ""
                for j := 0; j < 3; j++ {
                    currCell := lipgloss.NewStyle().Padding(0, 1, 0, 1).Background(primaryGlossColor).Render(cell)
                    currRow = lipgloss.JoinHorizontal(lipgloss.Center, currRow, currCell)
                }
                cellString = lipgloss.JoinVertical(lipgloss.Center, cellString, currRow)
            }
            return cellString
        } else {
            for i := 0; i < 3; i++ {
                currRow := ""
                for j := 0; j < 3; j++ {
                    currCell := lipgloss.NewStyle().Padding(0, 1, 0, 1).Background(secondaryGlossColor).Render(cell)
                    currRow = lipgloss.JoinHorizontal(lipgloss.Center, currRow, currCell)
                }
                cellString = lipgloss.JoinVertical(lipgloss.Center, cellString, currRow)
            }
            return cellString
        }
    }

    // renders cell
    drawCell = func(isWrong, isSelected, given bool, cell string) string {
        if isWrong {
            return drawFullCell(isSelected, WRONG_SELECTED_COLOR, WRONG_NOT_SELECTED_COLOR, cell)
        } else if isSelected {
            return drawFullCell(given, SELECTED_GIVEN_COLOR, SELECTED_NOT_GIVEN_COLOR, cell) 
        } else {
            return drawFullCell(given, BASE_GIVEN_COLOR, BASE_NOT_GIVEN_COLOR, cell)
        }
    }

    // takes string with direction(vert, hor) and bool for bold
    drawBorder = func(dir string, bold bool) string {
        if dir == "vert" {
            return drawVerticalBorder(bold)
        } else {
            return drawHorizontalBorder(bold)
        }
    }

    // returns vertical border string
    drawVerticalBorder = func(bold bool) string {
        var foregroundColor lipgloss.Color         
        if bold {
            foregroundColor = BOLD_BORDER_COLOR
        } else {
            foregroundColor = NORMAL_BORDER_COLOR
        }
        border := lipgloss.NewStyle().Padding(0, 1, 0, 1).Foreground(foregroundColor).Render("│")
        return lipgloss.JoinVertical(lipgloss.Center, border, border, border)
    }

    // returns horizontal border string
    drawHorizontalBorder = func(bold bool) string {
        var foregroundColor lipgloss.Color         
        if bold {
            foregroundColor = BOLD_BORDER_COLOR
        } else {
            foregroundColor = NORMAL_BORDER_COLOR
        }
        border := lipgloss.NewStyle().Padding(0, 0, 0, 0).Foreground(foregroundColor).Render("────────────────────────────────────")
        return lipgloss.JoinHorizontal(lipgloss.Center, border, border, border)
    }

//return lipgloss.NewStyle().Padding(0, 1, 0, 1).Bold(bold).Render("├───────────┼───────────┼───────────┤")
    
    
    drawSideBorder = func(dir string, loc string) string {
        if dir == "vert" {
            return lipgloss.NewStyle().Padding(0, 1, 0, 1).Render("│")
        } else {
            if loc == "top" {
                return lipgloss.NewStyle().Padding(0, 1, 0, 1).Render("┌───────────┬───────────┬───────────┐")
            } else {
                return lipgloss.NewStyle().Padding(0, 1, 0, 1).Render("└───────────┴───────────┴───────────┘")
            }
        }
    }
)

