package errwrap

import (
	"encoding/json"
	"fmt"
)

type Error interface {
	Error() string
	String() string
	GetStatusCode() int
	GetStatusMessage() string
	GetTraceID() string
	SetTraceID(traceID string)
}

type Err struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
}

func (e *Err) Error() string {
	if e == nil {
		return ""
	}
	return e.String()
}

func (e *Err) String() string {
	if e == nil {
		return ""
	}
	bs, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("marshalling err failed: %s", err)
	}
	return string(bs)
}

func (e *Err) GetStatusCode() int {
	if e == nil {
		return 0
	}
	return e.Code
}

func (e *Err) GetStatusMessage() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *Err) GetTraceID() string {
	if e == nil {
		return ""
	}
	return e.TraceID
}

func (e *Err) SetTraceID(traceID string) {
	if e == nil {
		return
	}
	e.TraceID = traceID
}

func NewError(code uint, format string, a ...any) Error {
	return &Err{
		Code:    int(code),
		Message: fmt.Sprintf(format, a...),
	}
}
