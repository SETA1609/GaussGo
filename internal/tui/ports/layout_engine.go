package ports

type Align string

const (
	AlignLeft   Align = "left"
	AlignCenter Align = "center"
	AlignRight  Align = "right"
	AlignTop    Align = "top"
	AlignBottom Align = "bottom"
)

type LayoutEngine interface {
	Place(width int, height int, hAlign Align, vAlign Align, content string) string
	PlaceHorizontal(width int, hAlign Align, content string) string
	JoinHorizontal(vAlign Align, parts ...string) string
	JoinVertical(hAlign Align, parts ...string) string
	Block(content string, width int, hAlign Align) string
}
