package ports

type Spinner interface {
	InitCmd() any
	Update(msg any) (cmd any, handled bool)
	View() string
}
