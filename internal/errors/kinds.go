package errors

type Kind string

const (
	KindUnknown        Kind = "unknown"
	KindInput          Kind = "input"
	KindDomain         Kind = "domain"
	KindInfrastructure Kind = "infrastructure"
)

type Error struct {
	Code    string
	Kind    Kind
	Message string
}

func (e Error) Error() string {
	return e.Message
}
