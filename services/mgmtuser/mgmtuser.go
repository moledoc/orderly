package mgmtuser

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/moledoc/orderly/middleware"
	"github.com/moledoc/orderly/models"
	"github.com/moledoc/orderly/services/common"
	"github.com/moledoc/orderly/storage"
	"github.com/moledoc/orderly/storage/local"
	"github.com/moledoc/orderly/utils"
)

var (
	strg storage.StorageUserAPI = nil
)

func handlePostUser(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handlePostUser")
	defer middleware.SpanStop(ctx, "handlePostUser")

	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		err := models.NewError(http.StatusBadRequest, "invalid payload")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	{ // TODO: move validation to _api
		// validation
		if user.ID != nil {
			err := models.NewError(http.StatusBadRequest, "id not allowed")
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}
	}

	resp, err := HandlePostUser(ctx, &user)
	common.WriteResponse(ctx, w, resp, err, http.StatusCreated)
}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleGetUserByID")
	defer middleware.SpanStop(ctx, "handleGetUserByID")

	id, errAtoi := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if errAtoi != nil {
		err := models.NewError(http.StatusBadRequest, "invalid id")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	resp, err := HandleGetUserByID(ctx, uint(id))
	common.WriteResponse(ctx, w, resp, err, http.StatusOK)
}

func handleGetUserVersions(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleGetUserVersions")
	defer middleware.SpanStop(ctx, "handleGetUserVersions")

	id, errAtoi := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if errAtoi != nil {
		err := models.NewError(http.StatusBadRequest, "invalid id")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	resp, err := HandleGetUserVersions(ctx, uint(id))
	common.WriteResponse(ctx, w, resp, err, http.StatusOK)
}

func handleGetUserSubOrdinates(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleGetUserSubOrdinates")
	defer middleware.SpanStop(ctx, "handleGetUserSubOrdinates")

	id, errAtoi := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if errAtoi != nil {
		err := models.NewError(http.StatusBadRequest, "invalid id")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	resp, err := HandleGetUserSubOrdinates(ctx, uint(id))
	common.WriteResponse(ctx, w, resp, err, http.StatusOK)
}

func handlePatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handlePatchUser")
	defer middleware.SpanStop(ctx, "handlePatchUser")

	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		err := models.NewError(http.StatusBadRequest, "invalid payload")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	{ // TODO: move validation to _api
		// validation
		if user.ID == nil {
			err := models.NewError(http.StatusBadRequest, "id must be provided")
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}
	}

	resp, err := HandlePatchUser(ctx, &user)
	common.WriteResponse(ctx, w, resp, err, http.StatusOK)
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleDeleteUser")
	defer middleware.SpanStop(ctx, "handleDeleteUser")

	id, errAtoi := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if errAtoi != nil {
		err := models.NewError(http.StatusBadRequest, "invalid id")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	var resp *models.User
	var err models.IError

	deleteType := r.URL.Query().Get("type")
	if deleteType == "hard" {
		resp, err = HandleDeleteUserHard(ctx, &models.User{
			ID: utils.Ptr(uint(id)),
		})
	} else {
		resp, err = HandleDeleteUserSoft(ctx, &models.User{
			ID: utils.Ptr(uint(id)),
		})
	}
	common.WriteResponse(ctx, w, resp, err, http.StatusNoContent)
}

func New() {

	strg = local.NewStorageUser()

	http.HandleFunc("POST /user", handlePostUser)
	http.HandleFunc("GET /user/{id}", handleGetUserByID)

	http.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		ctx := middleware.AddTrace(context.Background(), w)

		resp, err := HandleGetUsers(ctx)
		if err != nil {
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}

		bs, jsonerr := json.Marshal(resp)
		if jsonerr != nil {
			err := models.NewError(http.StatusInternalServerError, "marshalling user failed")
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bs)
	})

	http.HandleFunc("GET /user/{id}/versions", handleGetUserVersions)
	http.HandleFunc("GET /user/{id}/subordinates", handleGetUserSubOrdinates)
	http.HandleFunc("PATCH /user", handlePatchUser)
	http.HandleFunc("DELETE /user/{id}", handleDeleteUser)
}
