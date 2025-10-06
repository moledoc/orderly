package errwrap

import "fmt"

type Error interface {
	Error() string
	String() string
	GetStatusCode() int
	GetStatusMessage() string
}

type Err struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e Err) Error() string {
	return e.String()
}

func (e Err) String() string {
	return fmt.Sprintf(`{"code": %d, "message": %q}`, e.Code, e.Message)
}

func (e Err) GetStatusCode() int {
	return e.Code
}
func (e Err) GetStatusMessage() string {
	return e.Message
}

func NewError(code uint, format string, a ...any) Error {
	return Err{
		Code:    int(code),
		Message: fmt.Sprintf(format, a...),
	}
}
