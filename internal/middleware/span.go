package middleware

import (
	"context"
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

func fixTraceSpanEndTimes(spans []*span.Span) {
	spansMutex.Lock()
	defer spansMutex.Unlock()

	emptyTime := time.Time{}
	var prevEnd time.Time
	for _, spn := range spans {
		if spn.End.Equal(emptyTime) {
			if prevEnd.Equal(emptyTime) {
				spn.End = time.Now().UTC()
			} else {
				spn.End = prevEnd
			}
			prevEnd = spn.End
			spn.Duration = span.SpanDuration(spn.End.Sub(spn.Start))
		}
	}
}

func GetSpansByTrace(traceID string) []*span.Span {
	spans := spanss[traceID]
	fixTraceSpanEndTimes(spans)
	return spans
}

func SpanFlushTrace(ctx context.Context) {

	var traceID string
	ctxTraceID := ctx.Value(consts.CtxKeyTrace)
	if ctxTraceID != nil {
		traceID = ctxTraceID.(string)
	}

	spans := GetSpansByTrace(traceID)

	spansMutex.Lock()
	delete(spanss, traceID)
	spansMutex.Unlock()

	for _, s := range spans {
		fmt.Println(s)
	}

}

func SpanStart(ctx context.Context, desc string) {
	spansMutex.Lock()
	defer spansMutex.Unlock()

	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	var traceID string
	ctxTraceID := ctx.Value(consts.CtxKeyTrace)
	if ctxTraceID != nil {
		traceID = ctxTraceID.(string)
	}

	s := &span.Span{
		FuncName: fn.Name(),
		Filename: file,
		Line:     line,
		TraceID:  traceID,
		Start:    time.Now().UTC(),
		Desc:     desc,
	}
	spanss[traceID] = append(spanss[traceID], s)
}

func SpanStop(ctx context.Context, desc string) {
	spansMutex.Lock()
	defer spansMutex.Unlock()

	pc, file, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	var traceID string
	ctxTraceID := ctx.Value(consts.CtxKeyTrace)
	if ctxTraceID != nil {
		traceID = ctxTraceID.(string)
	}

	spans, ok := spanss[traceID]
	if !ok {
		return
	}
	for _, spn := range spans {
		if spn.Filename == file && spn.FuncName == fn.Name() && spn.Desc == desc {
			spn.End = time.Now().UTC()
			spn.Duration = span.SpanDuration(spn.End.Sub(spn.Start))
			break
		}
	}
}

func SpanLog(ctx context.Context, desc string, val any) {
	spansMutex.Lock()
	defer spansMutex.Unlock()

	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	var traceID string
	ctxTraceID := ctx.Value(consts.CtxKeyTrace)
	if ctxTraceID != nil {
		traceID = ctxTraceID.(string)
	}

	now := time.Now().UTC()
	s := &span.Span{
		FuncName: fn.Name(),
		Filename: file,
		Line:     line,
		TraceID:  traceID,
		Start:    now,
		End:      now,
		Desc:     desc,
		Val:      val,
	}
	spanss[traceID] = append(spanss[traceID], s)
}
