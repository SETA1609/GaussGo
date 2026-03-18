package errors

type ErrorType string

const (
	ErrorTypeUnknown        ErrorType = "unknown"
	ErrorTypeInput          ErrorType = "input"
	ErrorTypeDomain         ErrorType = "domain"
	ErrorTypeInfrastructure ErrorType = "infrastructure"
)

type Error struct {
	Code      string
	ErrorType ErrorType
	Message   string
}

func (e Error) Error() string {
	return e.Message
}
