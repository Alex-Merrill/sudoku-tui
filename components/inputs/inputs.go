package inputs

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	ShiftUp      key.Binding
	ShiftDown    key.Binding
	ShiftLeft    key.Binding
	ShiftRight   key.Binding
	Number       key.Binding
	PencilNumber key.Binding
	Delete       key.Binding
	Undo         key.Binding
	Redo         key.Binding
	Quit         key.Binding
	Help         key.Binding
	NewGame      key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.NewGame, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.ShiftUp, k.ShiftDown, k.ShiftLeft, k.ShiftRight}, // first column
		{k.Number, k.PencilNumber, k.Delete, k.Undo, k.Redo},                               // third column
		{k.Help, k.Quit, k.NewGame},                                                        // fifth column
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
	ShiftUp: key.NewBinding(
		key.WithKeys("K", "shift+up"),
		key.WithHelp("shift+↑/shift+k", "highlight up"),
	),
	ShiftDown: key.NewBinding(
		key.WithKeys("J", "shift+down"),
		key.WithHelp("shift+↓/shift+j", "highlight down"),
	),
	ShiftLeft: key.NewBinding(
		key.WithKeys("H", "shift+left"),
		key.WithHelp("shift+←/shift+h", "highlight left"),
	),
	ShiftRight: key.NewBinding(
		key.WithKeys("L", "shift+right"),
		key.WithHelp("shift+→/shift+l", "highlight right"),
	),
	Number: key.NewBinding(
		key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
		key.WithHelp("1-9", "input number"),
	),
	PencilNumber: key.NewBinding(
		key.WithKeys("!", "@", "#", "$", "%", "^", "&", "*", "("),
		key.WithHelp("shift+[1-9]", "pencil mark/unmark number"),
	),
	Delete: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "clear cell value"),
	),
	Undo: key.NewBinding(
		key.WithKeys("ctrl+z"),
		key.WithHelp("ctrl+z", "undo action"),
	),
	Redo: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "redo action"),
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
