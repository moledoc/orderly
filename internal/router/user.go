package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/middleware"
)

func handlePostUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "postUser")
	defer middleware.SpanStop(ctx, "postUser")

	var req request.PostUserRequest
	var resp *response.PostUserResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "PostUserRequest", &req)
		resp, err = mgmtusersvc.PostUser(ctx, &req)
	}

	if err == nil {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/user/%v", resp.GetUser().GetID()))
	}
	writeResponse(ctx, w, resp, err, http.StatusCreated)
}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "getUserByID")
	defer middleware.SpanStop(ctx, "getUserByID")

	req := &request.GetUserByIDRequest{
		ID: meta.ID(r.PathValue(userID)),
	}
	middleware.SpanLog(ctx, "GetUserByIDRequest", req)
	resp, err := mgmtusersvc.GetUserByID(ctx, req)

	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "getUsers")
	defer middleware.SpanStop(ctx, "getUsers")

	queryEmails := r.URL.Query()["emails"]
	emails := make([]user.Email, len(queryEmails))
	for i, em := range queryEmails {
		emails[i] = user.Email(em)
	}
	req := &request.GetUsersRequest{
		Emails:     emails,
		Supervisor: user.Email(r.URL.Query().Get("supervisor")),
	}
	middleware.SpanLog(ctx, "GetUsersRequest", req)
	resp, err := mgmtusersvc.GetUsers(ctx, req)
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func handlePatchUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "patchUser")
	defer middleware.SpanStop(ctx, "patchUser")

	var req request.PatchUserRequest
	var resp *response.PatchUserResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "PatchUserRequest", &req)
		resp, err = mgmtusersvc.PatchUser(ctx, &req)
	}

	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "deleteUser")
	defer middleware.SpanStop(ctx, "deleteUser")

	req := &request.DeleteUserRequest{
		ID: meta.ID(r.PathValue(userID)),
	}
	middleware.SpanLog(ctx, "DeleteUserRequest", req)
	resp, err := mgmtusersvc.DeleteUser(ctx, req)
	writeResponse(ctx, w, resp, err, http.StatusNoContent)
}
