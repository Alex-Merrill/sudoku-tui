package winscreen

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hisamafahri/coco"
)

type Model struct {
	animationState string
	sourceText     []string
	color          lipgloss.Color

	textColStartIdx int
	textColEndIdx   int

	width, height int
}

type StopAnim struct{}

type frameMsg struct {
	t time.Time
}

const (
	fps             = 60
	animationLength = 5 * time.Second
	bannerSpeed     = 1 // how many cols move per frame
)

func NewModel(w, h int) Model {
	return Model{
		animationState:  "",
		sourceText:      getTextToDisplay(w),
		color:           getRandomColor(),
		textColStartIdx: 0,
		textColEndIdx:   2,
		width:           w,
		height:          h,
	}
}

func animate() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return frameMsg{
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
	if m.textColEndIdx-m.textColStartIdx > m.width {
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
	r, g, b := int(rgb[0]), int(rgb[1]), int(rgb[2])
	hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)

	return lipgloss.Color(hex)
}

func getTextToDisplay(w int) []string {
	lines := []string{}

	lines = append(lines, "YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!Y"+strings.Repeat(" ", w))
	lines = append(lines, "ouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!Yo"+strings.Repeat(" ", w))
	lines = append(lines, "uWon!         YouWo n        !YouWon                      !YouWon!YouWon         !YouWo        n!YouWon!You         Won!Yo        uWon!Y         ouWon!YouW                     on!YouWon!You          Won!Y         ouWon!YouWon!YouWon!YouWon!YouWon!YouWon"+strings.Repeat(" ", w))
	lines = append(lines, "!Yo             uW            on!Yo                         uWon!YouWon                          !YouWon!Y            ouW           on!            YouWon!                        YouWon!You                          Won!YouWon           !YouWon!YouWon!You"+strings.Repeat(" ", w))
	lines = append(lines, "Won    !YouWon!    YouWon!    YouWo  n!YouWon!YouWon!YouWo   n!YouWon!Y   ouWon!Yo      uWon!Yo   uWon!You    Won!You   W  on!YouWo       n!YouWon  !YouWo  n!YouWon!YouWon!YouW   on!YouWon    !YouWon!      YouWon!   YouWon!Y    ouWon!   YouWon!YouWon!Yo"+strings.Repeat(" ", w))
	lines = append(lines, "uWo     n!YouWon   !YouWon     !You    Won!YouWon   !YouWon    !YouWon!     YouWon!      YouWon!Y  ouWon!Y     ouWon!Yo     uWon!YouW      on!YouWo  n!You    Won!YouWo    n!YouWo   n!YouWo     n!YouWon!Y    ouWon!Y   ouWon!Y     ouWon!Y  ouWon!YouWon!Yo"+strings.Repeat(" ", w))
	lines = append(lines, "uWo       n!YouWon !YouWon     !You     Won!YouW      on!YouW   on!YouW      on!YouWo      n!YouWo   n!You       Won!YouW   on!YouWon!Yo    uWon!You  Won!     YouWon!Yo     uWon!Yo  uWon!Y       ouWon!YouWo  n!YouWon   !YouW      on!YouW  on!YouWon!YouW"+strings.Repeat(" ", w))
	lines = append(lines, "on!         YouWon!YouWon!     YouW      on!YouWon     !YouWon    !YouW       on!YouWo      n!YouWo   n!Yo         uWon!Yo  uWon!YouWon!Yo  uWon!YouW  on!      YouWon!Yo     uWon!Yo   uWon!       YouWon!YouWon!YouWon!Y  ouWon       !YouWo   n!YouWon!You"+strings.Repeat(" ", w))
	lines = append(lines, "Won             !YouWon!You     Won        !YouWon!      YouWon!   YouW        on!YouWo      n!YouWon  !Yo          uWon!YouWon!YouWo n!YouWon!YouWon!  Yo        uWon!You      Won!You  Won!        YouWon!YouWon!YouWon!Y   ouWo       n!YouWo  n!YouWon!Yo"+strings.Repeat(" ", w))
	lines = append(lines, "uWo                 n!YouWo     n!Y         ouWon!Yo      uWon!Yo    uWo         n!YouWon     !YouWon!   Yo           uWon!YouWon!You   Won!YouWon!YouW  on        !YouWon!      YouWon!   You         Won!YouWon!YouWon!YouW  on!Y       ouWon!Y  ouWon!YouW"+strings.Repeat(" ", w))
	lines = append(lines, "on!Y                  ouWon!     You          Won!YouW     on!YouWo   n!Y         ouWon!Yo      uWon!You  Wo            n!YouWon!YouW      on!YouWon!You  Wo        n!YouWon!     YouWon!Y  ouW         on!YouWon  !YouWon!You   Won                !YouWon!Y"+strings.Repeat(" ", w))
	lines = append(lines, "ouWon!Y                ouWon!     YouW         on!YouWo      n!YouWo   n!Y         ouWon!Yo      uWon!You  Won            !YouWon!You        Won!YouWon!Y  ou         Won!YouW     on!YouWo   n!Y         ouWon!Y    ouWon!YouWo  n!Yo                uWon!Yo"+strings.Repeat(" ", w))
	lines = append(lines, "uWon!YouWon             !YouWo      n!Yo        uWon!YouW    on!YouWon   !Yo         uWon!You    Won!YouWo   n!Yo          uWon!YouWo          n!YouWon!Yo  uWo        n!YouWon!    YouWon!Yo  uWon        !YouWon!    YouWon!You   Won      !YouWon!  YouWon"+strings.Repeat(" ", w))
	lines = append(lines, "!YouWon!YouWon!           YouWon     !YouWo       n!YouWon!YouWon!YouWo    n!Y        ouWon!YouWon!YouWon!Yo  uWon!Y         ouWon!Yo             uWon!YouW  on!Y        ouWon!YouWon!YouWon!Y  ouWon       !YouWon!      YouWon!Yo  uWon     !YouWon!  YouWo"+strings.Repeat(" ", w))
	lines = append(lines, "n!YouWon!YouWon!You        Won!Yo      uWon!Y      ouWon!YouWon!YouWon!Yo   uWon!      YouWon!YouWon!YouWon!Y  ouWon!Yo        uWon!     Y          ouWon!Yo  uWon!Y      ouWon!YouWon!YouWon!Y  ouWon!       YouWon!       YouWon!Y   ouWo    n!YouWon  !You"+strings.Repeat(" ", w))
	lines = append(lines, "Won!YouWon!YouWon!YouWo                n!YouWon                             !YouWon                            !YouWon!You               Won!Y                 ouWon!Y                           ouWon!You                             Won!Yo            uWon"+strings.Repeat(" ", w))
	lines = append(lines, "!YouWon!YouWon!YouWon!YouW             on!YouWon!Y                         ouWon!YouWo                        n!YouWon!YouWon            !YouWon!Y             ouWon!YouW                        on!YouWon!Y             ou            Won!You          Won!Y"+strings.Repeat(" ", w))
	lines = append(lines, "ouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!Yo"+strings.Repeat(" ", w))
	lines = append(lines, "uWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!YouWon!You"+strings.Repeat(" ", w))

	return lines
}

// throw error
func check(e error) {
	if e != nil {
		panic(e)
	}
}
