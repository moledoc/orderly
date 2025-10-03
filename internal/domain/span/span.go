package span

import (
	"time"
)

type CtxKey struct {
	Key string
}

type Span struct {
	FuncName string    `json:"func_name,omitempty"`
	Filename string    `json:"filename,omitempty"`
	Line     int       `json:"line,omitempty"`
	TraceID  string    `json:"trace,omitempty"`
	Start    time.Time `json:"start,omitempty"`
	End      time.Time `json:"end,omitempty"`
	Desc     string    `json:"desc,omitempty"`
}

func (s *Span) GetFuncName() string {
	if s == nil {
		return ""
	}
	return s.FuncName
}

func (s *Span) GetFilename() string {
	if s == nil {
		return ""
	}
	return s.Filename
}

func (s *Span) GetLine() int {
	if s == nil {
		return 0
	}
	return s.Line
}

func (s *Span) GetTraceID() string {
	if s == nil {
		return ""
	}
	return s.TraceID
}

func (s *Span) GetStart() time.Time {
	if s == nil {
		return time.Time{}
	}
	return s.Start
}

func (s *Span) GetEnd() time.Time {
	if s == nil {
		return time.Time{}
	}
	return s.End
}

func (s *Span) GetDesc() string {
	if s == nil {
		return ""
	}
	return s.Desc
}

type Spans map[string][]*Span
