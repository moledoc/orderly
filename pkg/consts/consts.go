package consts

import "github.com/moledoc/orderly/internal/domain/span"

type Action int

const (
	CREATE Action = iota
	UPDATE
	DELETESOFT
	DELETEHARD
	READ
	READALL
	READVERSIONS
	READSUBORDINATES
	READSUBORDERS
)

const (
	TraceID string = "Trace-Id"
)

var (
	CtxKeyTrace = span.CtxKey{Key: TraceID}
)
