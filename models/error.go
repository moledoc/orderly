package models

import "fmt"

type IError interface {
	String() string
	StatusCode() int
	StatusMessage() string
}

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) String() string {
	return fmt.Sprintf(`{"code": %d, "message": %q}`, e.Code, e.Message)
}

func (e *Error) StatusCode() int {
	return e.Code
}
func (e *Error) StatusMessage() string {
	return e.Message
}

func NewError(code uint, format string, a ...any) IError {
	return &Error{
		Code:    int(code),
		Message: fmt.Sprintf(format, a...),
	}
}
