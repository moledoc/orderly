package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/moledoc/orderly/internal/domain/span"
	"github.com/moledoc/orderly/pkg/consts"
)

var (
	spanss     map[string][]*span.Span = make(map[string][]*span.Span)
	spansMutex sync.Mutex
)

func SpanFlushTrace(ctx context.Context) {
	spansMutex.Lock()
	defer spansMutex.Unlock()

	trace := ctx.Value(consts.CtxKeyTrace).(string)
	spans := spanss[trace]
	emptyTime := time.Time{}
	var prevEnd time.Time
	for _, span := range spans {
		if span.End.Equal(emptyTime) {
			if prevEnd.Equal(emptyTime) {
				span.End = time.Now().UTC()
			} else {
				span.End = prevEnd
			}
			prevEnd = span.End
		}
		bs, err := json.Marshal(span)
		if err == nil {
			fmt.Println(string(bs))
		}
	}
	delete(spanss, trace)
}

func SpanStart(ctx context.Context, desc string) {
	spansMutex.Lock()
	defer spansMutex.Unlock()

	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	trace := ctx.Value(consts.CtxKeyTrace).(string)
	s := &span.Span{
		FuncName: fn.Name(),
		Filename: file,
		Line:     line,
		TraceID:  trace,
		Start:    time.Now().UTC(),
		Desc:     desc,
	}
	spanss[trace] = append(spanss[trace], s)
}

func SpanStop(ctx context.Context, desc string) {
	spansMutex.Lock()
	defer spansMutex.Unlock()

	pc, file, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	trace := ctx.Value(consts.CtxKeyTrace).(string)
	spans, ok := spanss[trace]
	if !ok {
		return
	}
	for _, span := range spans {
		if span.Filename == file && span.FuncName == fn.Name() && span.Desc == desc {
			span.End = time.Now().UTC()
			break
		}
	}
}
