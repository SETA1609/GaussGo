package ports

type Zone interface {
	Init()
	Mark(id string, content string) string
	Scan(content string) string
}
