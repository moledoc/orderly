package models

import (
	"time"
)

type CtxKey struct {
	key string
}

type Span struct {
	FuncName string    `json:"func_name,omitempty"`
	Filename string    `json:"filename,omitempty"`
	Line     int       `json:"line,omitempty"`
	Trace    string    `json:"trace,omitempty"`
	Start    time.Time `json:"start,omitempty"`
	End      time.Time `json:"end,omitempty"`
	Desc     string    `json:"desc,omitempty"`
}

type Spans map[string][]*Span

var (
	CtxKeyTrace = CtxKey{key: "trace"}
)
