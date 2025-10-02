package common

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/moledoc/orderly/middleware"
	"github.com/moledoc/orderly/models"
)

func WriteResponse(ctx context.Context, w http.ResponseWriter, resp any, err models.IError, successCode int) {
	middleware.SpanStart(ctx, "WriteResponse")
	defer middleware.SpanStop(ctx, "WriteResponse")

	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(resp)
	if jsonerr != nil {
		err := models.NewError(http.StatusInternalServerError, "marshalling resp failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	w.WriteHeader(successCode)
	w.Write(bs)
}
