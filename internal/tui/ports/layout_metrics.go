package ports

type LayoutMetrics interface {
	Height(content string) int
	Width(content string) int
}
