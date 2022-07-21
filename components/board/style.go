package board

import (

	"github.com/charmbracelet/lipgloss"
)

// style definitions and ui rendering
var (
    
    // defines base cell style
    baseCell = func(given bool) lipgloss.Style {
        // defines style for cell that is given to user
        if given {
           //return lipgloss.NewStyle().Padding(0, 1, 0, 1).Background(lipgloss.Color("#595959")) 
           return lipgloss.NewStyle().Padding(0, 1, 0, 1).Background(lipgloss.Color("#02182B")) 
        } else { // defines syle for cell that is empty/user controlled
           return lipgloss.NewStyle().Padding(0, 1, 0, 1).Background(lipgloss.Color("#68C5DB")) 
        }
    }

    // defines selected cell style
    selectedCell = func(given bool) lipgloss.Style {
        // defines style for cell that is given to user
        if given {
           return lipgloss.NewStyle().Padding(0, 1, 0, 1).Background(lipgloss.Color("#0164A2")) 
        } else { // defines syle for cell that is empty/user controlled
           return lipgloss.NewStyle().Padding(0, 1, 0, 1).Background(lipgloss.Color("#0197F6")) 
        }
    }

    // defines wrong cell style
    wrongCell = func(isSelected bool) lipgloss.Style {
        if isSelected {
            return lipgloss.NewStyle().Padding(0, 1, 0, 1).Background(lipgloss.Color("#E77484"))
        } else {
            return lipgloss.NewStyle().Padding(0, 1, 0, 1).Background(lipgloss.Color("#D7263D"))
        }
    }

    // renders cell
    drawCell = func(isWrong, isSelected, given bool, cell string) string {
        var style lipgloss.Style
        if isWrong {
            style = wrongCell(isSelected)
        } else if isSelected {
            style = selectedCell(given)
        } else {
            style = baseCell(given)
        }
        
        return style.Render(cell)
    }

    drawBorder = func(dir string) string {
        if dir == "vert" {
            return lipgloss.NewStyle().Padding(0, 1, 0, 1).Render("│")
        } else {
            return lipgloss.NewStyle().Padding(0, 1, 0, 1).Render("──────────┼───────────┼──────────")
        }
    }
)

