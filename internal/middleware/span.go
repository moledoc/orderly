package middleware

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/moledoc/orderly/internal/domain/span"
	"github.com/moledoc/orderly/pkg/consts"
)

var (
	spanss          map[string][]*span.Span = make(map[string][]*span.Span)
	spansMutex      sync.Mutex
	spanLogFilename = "/tmp/orderly.spans.log"
)

func fixTraceSpanEndTimes(spans []*span.Span) {
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

func SpanFlushTrace(ctx context.Context) {
	spansMutex.Lock()
	defer spansMutex.Unlock()

	fptr, err := os.OpenFile(spanLogFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		fmt.Printf("[WARNING]: failed to open %v: %v\n", spanLogFilename, err)
	}
	defer fptr.Close()

	var traceID string
	ctxTraceID := ctx.Value(consts.CtxKeyTrace)
	if ctxTraceID != nil {
		traceID = ctxTraceID.(string)
	}

	spans := spanss[traceID]
	fixTraceSpanEndTimes(spans)
	delete(spanss, traceID)

	var buf string
	for _, s := range spans {
		buf += fmt.Sprintf("%s\n", s)
	}
	n, err := fptr.Write([]byte(buf))
	if n != len(buf) {
		fmt.Fprintf(os.Stderr, "[WARNING]: wrote diff nr bytes to '%s': expected %v, wrote %v\n", spanLogFilename, len(buf), n)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "[WARNING]: failed to write span logs to '%s': %s\n", spanLogFilename, err)
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
