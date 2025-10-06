package router

import (
	"context"
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/pkg/utils"
)

func postUser(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "postUser")
	defer middleware.SpanStop(ctx, "postUser")

	var req request.PostUserRequest
	var resp *response.PostUserResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		resp, err = mgmtusersvc.PostUser(ctx, &req)
	}

	writeResponse(ctx, w, resp, err, http.StatusCreated)
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "getUserByID")
	defer middleware.SpanStop(ctx, "getUserByID")

	resp, err := mgmtusersvc.GetUserByID(ctx, &request.GetUserByIDRequest{
		ID: utils.Ptr(meta.ID(r.PathValue(userID))),
	})
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "getUsers")
	defer middleware.SpanStop(ctx, "getUsers")

	resp, err := mgmtusersvc.GetUsers(ctx, &request.GetUsersRequest{})
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func getUserSubOrdinates(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "getUserSubOrdinates")
	defer middleware.SpanStop(ctx, "getUserSubOrdinates")

	resp, err := mgmtusersvc.GetUserSubOrdinates(ctx, &request.GetUserSubOrdinatesRequest{
		ID: utils.Ptr(meta.ID(r.PathValue(userID))),
	})
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func patchUser(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "patchUser")
	defer middleware.SpanStop(ctx, "patchUser")

	var req request.PatchUserRequest
	var resp *response.PatchUserResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		resp, err = mgmtusersvc.PatchUser(ctx, &req)
	}

	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "deleteUser")
	defer middleware.SpanStop(ctx, "deleteUser")

	resp, err := mgmtusersvc.DeleteUser(ctx, &request.DeleteUserRequest{
		ID: utils.Ptr(meta.ID(r.PathValue(userID))),
	})
	writeResponse(ctx, w, resp, err, http.StatusNoContent)
}
