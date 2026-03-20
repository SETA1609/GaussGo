package ports

type ViewConfig struct {
	AltScreen bool
	MouseMode string
}

type ViewPort interface {
	Render(text string, config ViewConfig) any
}
