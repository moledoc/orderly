package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/moledoc/orderly/actions"
	"github.com/moledoc/orderly/middleware"
	"github.com/moledoc/orderly/models"
	"github.com/moledoc/orderly/storage"
	"github.com/moledoc/orderly/storage/local"
	"github.com/moledoc/orderly/utils"
)

var (
	strg storage.StorageAPI = nil
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

	{
		// validation
		if user.ID != nil {
			err := models.NewError(http.StatusBadRequest, "id not allowed")
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}
	}

	u, err := strg.Write(ctx, actions.CREATE, &user)
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(u)
	if jsonerr != nil {
		err := models.NewError(http.StatusInternalServerError, "marshalling user failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(bs)
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

	u, err := strg.Read(ctx, actions.READ, uint(id))
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(u)
	if jsonerr != nil {
		err := models.NewError(http.StatusInternalServerError, "marshalling user failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleGetUsers")
	defer middleware.SpanStop(ctx, "handleGetUsers")

	us, err := strg.Read(ctx, actions.READALL, 0)
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(us)
	if jsonerr != nil {
		err := models.NewError(http.StatusInternalServerError, "marshalling user failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
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

	u, err := strg.Read(ctx, actions.READVERSIONS, uint(id))
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(u)
	if jsonerr != nil {
		err := models.NewError(http.StatusInternalServerError, "marshalling user failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
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

	u, err := strg.Read(ctx, actions.READSUBORDINATES, uint(id))
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(u)
	if jsonerr != nil {
		err := models.NewError(http.StatusInternalServerError, "marshalling user failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
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

	{
		// validation
		if user.ID == nil {
			err := models.NewError(http.StatusBadRequest, "id must be provided")
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}
	}

	u, err := strg.Write(ctx, actions.UPDATE, &user)
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(u)
	if jsonerr != nil {
		err := models.NewError(http.StatusInternalServerError, "marshalling user failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
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

	action := actions.DELETESOFT

	deleteType := r.URL.Query().Get("type")
	if deleteType == "hard" {
		action = actions.DELETEHARD
	}

	_, err := strg.Write(ctx, action, &models.User{ID: utils.Ptr(uint(id))})
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {

	strg = local.New()

	http.HandleFunc("POST /user", handlePostUser)
	http.HandleFunc("GET /user/{id}", handleGetUserByID)
	http.HandleFunc("GET /users", handleGetUsers)
	http.HandleFunc("GET /user/{id}/versions", handleGetUserVersions)
	http.HandleFunc("GET /user/{id}/subordinates", handleGetUserSubOrdinates)
	http.HandleFunc("PATCH /user", handlePatchUser)
	http.HandleFunc("DELETE /user/{id}", handleDeleteUser)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
