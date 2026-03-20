package apperrors

import "fmt"

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
	Cause     error
}

func New(code string, errorType ErrorType, message string) Error {
	return Error{Code: code, ErrorType: errorType, Message: message}
}

func Wrap(code string, errorType ErrorType, message string, cause error) Error {
	return Error{Code: code, ErrorType: errorType, Message: message, Cause: cause}
}

func (e Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s/%s] %s: %v", e.Code, e.ErrorType, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s/%s] %s", e.Code, e.ErrorType, e.Message)
}

func (e Error) Unwrap() error {
	return e.Cause
}
