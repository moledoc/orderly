package span

import (
	"encoding/json"
	"fmt"
	"time"
)

type CtxKey struct {
	Key string
}

type SpanDuration time.Duration

func (sd SpanDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(sd).String())
}

type Span struct {
	FuncName string       `json:"func_name,omitempty"`
	Filename string       `json:"filename,omitempty"`
	Line     int          `json:"line,omitempty"`
	TraceID  string       `json:"trace_id,omitempty"`
	Start    time.Time    `json:"start,omitempty"`
	End      time.Time    `json:"end,omitempty"`
	Duration SpanDuration `json:"duration,omitempty"`
	Desc     string       `json:"desc,omitempty"`
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

func (s *Span) GetDuration() SpanDuration {
	if s == nil {
		return 0
	}
	return s.Duration
}

func (s *Span) GetDesc() string {
	if s == nil {
		return ""
	}
	return s.Desc
}

func (s *Span) String() string {
	bs, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("marshalling error: %s", err)
	}
	return string(bs)
}
