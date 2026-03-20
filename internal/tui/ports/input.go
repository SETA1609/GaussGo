package ports

type InputEventType string

const (
	InputEventKey        InputEventType = "key"
	InputEventWindowSize InputEventType = "window_size"
)

type Key string

const (
	KeyUp        Key = "up"
	KeyDown      Key = "down"
	KeyK         Key = "k"
	KeyJ         Key = "j"
	KeyEnter     Key = "enter"
	KeyEsc       Key = "esc"
	KeyBackspace Key = "backspace"
	KeyQ         Key = "q"
	KeyCtrlC     Key = "ctrl+c"
	KeyLeftBracket  Key = "["
	KeyRightBracket Key = "]"
)

type InputEvent struct {
	Type   InputEventType
	Key    Key
	Width  int
	Height int
}

type InputMapper interface {
	Map(msg any) (InputEvent, bool)
}
