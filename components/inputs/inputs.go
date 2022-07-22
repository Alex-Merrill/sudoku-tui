package inputs

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
    Up key.Binding
    Down key.Binding
    Left key.Binding
    Right key.Binding
    Number key.Binding
    Delete key.Binding
    Quit key.Binding
    Help key.Binding
    NewGame key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Help, k.NewGame, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down, k.Left, k.Right}, // first column
        {k.Number, k.Delete}, // second column
        {k.Help, k.Quit, k.NewGame}, // fourth column
    }
}

var Controls = KeyMap{
    Up: key.NewBinding(
        key.WithKeys("k", "up"),
        key.WithHelp("↑/k", "move up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("j", "down"),
        key.WithHelp("↓/j", "move down"),
    ),
    Left: key.NewBinding(
        key.WithKeys("h", "left"),
        key.WithHelp("←/h", "move left"),
    ),
    Right: key.NewBinding(
        key.WithKeys("l", "right"),
        key.WithHelp("→/l", "move right"),
    ),
    Number: key.NewBinding(
        key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
        key.WithHelp("1-9", "enter number in cell"),
    ),
    Delete: key.NewBinding(
        key.WithKeys("backspace"),
        key.WithHelp("backspace", "clear cell"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("ctrl+c", "q"),
        key.WithHelp("ctrl+c/q", "quit"),
    ),
    Help: key.NewBinding(
        key.WithKeys("?"),
        key.WithHelp("?", "toggle help"),
    ),
    NewGame: key.NewBinding(
        key.WithKeys("n"),
        key.WithHelp("n", "new game"),
    ),
}
