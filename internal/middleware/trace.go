package middleware

import (
	"context"
	"net/http"

	"github.com/moledoc/orderly/pkg/consts"
	"github.com/moledoc/orderly/pkg/utils"
)

func AddTrace(ctx context.Context, w http.ResponseWriter) context.Context {
	var trace string
	if w != nil {
		trace = w.Header().Get(consts.TraceID)
	}
	if len(trace) == 0 {
		trace = utils.RandAlphanum()
		if w != nil {
			w.Header().Add(consts.TraceID, trace)
		}
	}
	if ctx.Value(consts.CtxKeyTrace) == nil {
		ctx = context.WithValue(ctx, consts.CtxKeyTrace, trace)
	}
	return ctx
}

func GetTrace(w http.ResponseWriter) string {
	if len(w.Header().Get(consts.TraceID)) == 0 {
		return ""
	}
	return w.Header().Get(consts.TraceID)
}
