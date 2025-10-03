package errwrap

import "fmt"

type Error interface {
	String() string
	GetStatusCode() int
	GetStatusMessage() string
}

type err struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *err) String() string {
	return fmt.Sprintf(`{"code": %d, "message": %q}`, e.Code, e.Message)
}

func (e *err) GetStatusCode() int {
	return e.Code
}
func (e *err) GetStatusMessage() string {
	return e.Message
}

func NewError(code uint, format string, a ...any) Error {
	return &err{
		Code:    int(code),
		Message: fmt.Sprintf(format, a...),
	}
}
