package winscreen

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hisamafahri/coco"
)

type Model struct {
    animationState string
    sourceText []string
    color lipgloss.Color

    textColStartIdx int
    textColEndIdx int

    width, height int
}

type StopAnim struct {}

type frameMsg struct {
    t time.Time
}

const (
    fps = 60
    animationLength = 5*time.Second
    bannerSpeed = 1 //how many cols move per frame
    filePath = "wedge.txt"
)

func NewModel(w, h int) Model {
    return Model{
        animationState: "",
        sourceText: getTextToDisplay(w),
        color: getRandomColor(),
        textColStartIdx: 0,
        textColEndIdx: 2,
        width: w,
        height: h,
    }
}

func animate() tea.Cmd {
    return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
        return frameMsg {
            t: t,
        }
    })
}

func wait(d time.Duration) tea.Cmd {
    return func() tea.Msg {
        time.Sleep(d)
        return nil
    }
} 

func (m Model) Init() tea.Cmd {
    return tea.Sequentially(wait(time.Second/8), animate())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg.(type) {
    // new frame
    case frameMsg:
        
        // update animationState
        m.updateAnim()

        // request next frame
        return m, animate()

    default:
        return m, nil
    }
}

func (m Model) View() string { 
    return lipgloss.Place(m.width, m.height, lipgloss.Right, lipgloss.Center, m.animationState)

}

// slides banner over by bannerSpeed columns
func (m *Model) updateAnim() {
    m.animationState = "" 
    
    for i := 0; i < len(m.sourceText); i++ {
        for j := m.textColStartIdx; j < m.textColEndIdx; j++ {
            line := string([]rune(m.sourceText[i])[j])
            m.animationState += lipgloss.NewStyle().Foreground(m.color).Render(line)
        }
        m.animationState += "\n"
    }
    
    // iterate columns to draw
    m.textColEndIdx += bannerSpeed
    
    // if we are taking up the full width of the window, also iterate the start column
    if m.textColEndIdx - m.textColStartIdx > m.width {
        m.textColStartIdx += bannerSpeed
    }
    
    // once banner has passed screen, start over
    if m.textColEndIdx >= len([]rune(m.sourceText[0])) {
        m.color = getRandomColor()
        m.textColStartIdx = 0
        m.textColEndIdx = 2
    }

}

func getRandomColor() lipgloss.Color {
    // Get random HSV values where s and v are in range 70-100
    h := rand.Float64()
    s := rand.Float64()
    v := rand.Float64()
    h *= 360
    s = s*30 + 70
    v = v*30 + 70
   
    // convert hsv value to rgb then to hex
    rgb := coco.Hsv2Rgb(h, s, v)
    r,g,b := int(rgb[0]), int(rgb[1]), int(rgb[2])
    hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)

    return lipgloss.Color(hex)
}

// read win text from file
func getTextToDisplay(w int) []string {
    lines := []string{}

    file,err := os.Open(filePath)
    check(err)

    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text() + strings.Repeat(" ", w))
    }
    
    err = scanner.Err()
    check(err)

    return lines
}

// throw error
func check(e error) {
    if e != nil {
        panic(e)
    }
}
