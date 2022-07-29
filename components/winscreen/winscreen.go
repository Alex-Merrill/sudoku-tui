package winscreen

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hisamafahri/coco"
)

type Model struct {
    animationState string
    animationOver bool  
    animationStarted bool
    messageCharIdx int
}

type StopAnim struct {}

type frameMsg struct {
    t time.Time
}

const (
    fps = 60
    animationLength = 5*time.Second
    messageToDisplay = "You Won!"
    lineLength = 5
)

func NewModel() Model {
    return Model{
        //animationState: "You Won! You Won! You Won! You Won! You Won!\nYou Won! You Won! You Won! You Won! You Won!\nYou Won! You Won! You Won! You Won! You Won!\nYou Won! You Won! You Won! You Won! You Won!\nYou Won! You Won! You Won! You Won! You Won!",
        animationState: getTextToDisplay(),
        animationOver: false,
        animationStarted: false,
        messageCharIdx: 0,
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    /*  if we do something more dynamic or takes places over certain time,
        we might want user to be able to stop animation
        case tea.KeyMsg:
            if m.animationStarted {
                m.animationOver = true
                return m, stopAnim()
            } else {
                return m, nil
            }
    */

    // new frame
    case frameMsg:
        // change animation started state and set animation start time
        if !m.animationStarted {
            m.animationStarted = true
        }
        // if animation should be over, don't request another frame
        if m.animationOver {
            return m, nil
        }

        /*  building animation string by frame
            if len(m.animationState) == 225 {
                m.animationOver = true
                return m, stopAnim()
            }
       
            // updates animation state
            m.updateAnimationState()
        */

        // request next frame
        return m, animate()

    default:
        return m, nil
    }
}

func (m Model) View() string { 
    stringToDisplay := ""

    charStyle := lipgloss.NewStyle()

    for _,c := range m.animationState {
        if string(c) == " " {
            randCol := getRandomColor()
            stringToDisplay += charStyle.Background(randCol).Render(string(c))
        } else {
            stringToDisplay += string(c)
        }
    }

    return stringToDisplay
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


func getTextToDisplay() string {
    data,err := os.ReadFile("orderedwinscreen.txt")
    check(err)

     return string(data)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}








func (m *Model) updateAnimationState() {
    r := regexp.MustCompile("\n")
    animationStateNoNL := r.ReplaceAllString(m.animationState, "")

    lineCharCount := (len(messageToDisplay)*lineLength + lineLength-1)

    // add character
    if m.messageCharIdx < len(messageToDisplay) {
        charToAdd := string([]rune(messageToDisplay)[m.messageCharIdx])
        m.animationState += charToAdd
        m.messageCharIdx++
    } else if len(animationStateNoNL) % lineCharCount == 0 { // add new line after lineLength messages 
        m.animationState += "\n"
        m.messageCharIdx = 0
    } else { // add space between message
        m.animationState += " "
        m.messageCharIdx = 0
    } 
}

func stopAnim() tea.Cmd {
    return func() tea.Msg {
        return StopAnim{}
    }
}

