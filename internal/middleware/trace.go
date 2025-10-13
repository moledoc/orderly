package middleware

import (
	"context"
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/pkg/consts"
	"github.com/moledoc/orderly/pkg/utils"
)

func AddTraceToCtx(ctx context.Context) context.Context {
	if ctx.Value(consts.CtxKeyTrace) == nil {
		ctx = context.WithValue(ctx, consts.CtxKeyTrace, utils.RandAlphanum())
	}
	return ctx
}

func AddTraceToCtxFromWriter(ctx context.Context, w http.ResponseWriter) context.Context {
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

func AddTraceToErrFromCtx(err errwrap.Error, ctx context.Context) errwrap.Error {
	if trace := ctx.Value(consts.CtxKeyTrace); trace != nil {
		err.SetTraceID(trace.(string))
	}
	return err
}

func GetTrace(w http.ResponseWriter) string {
	if len(w.Header().Get(consts.TraceID)) == 0 {
		return ""
	}
	return w.Header().Get(consts.TraceID)
}
