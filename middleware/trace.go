package middleware

import (
	"context"
	"net/http"

	"github.com/moledoc/orderly/models"
	"github.com/moledoc/orderly/utils"
)

func AddTrace(ctx context.Context, w http.ResponseWriter) context.Context {
	var trace string
	if w != nil {
		trace = w.Header().Get("trace")
	}
	if len(trace) == 0 {
		trace = utils.RandAlphanum()
		if w != nil {
			w.Header().Add("trace", trace)
		}
	}
	if ctx.Value(models.CtxKeyTrace) == nil {
		ctx = context.WithValue(ctx, models.CtxKeyTrace, trace)
	}
	return ctx
}

func GetTrace(w http.ResponseWriter) string {
	if len(w.Header().Get("trace")) == 0 {
		return ""
	}
	return w.Header().Get("trace")
}
