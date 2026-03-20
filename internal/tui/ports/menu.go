package ports

type MenuRenderer interface {
	Mark(id string, content string) string
}
