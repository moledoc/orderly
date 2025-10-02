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
	strg storage.StorageOrderAPI = nil
)

func handlePostOrder(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handlePostOrder")
	defer middleware.SpanStop(ctx, "handlePostOrder")

	var order models.Order
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&order); err != nil {
		err := models.NewError(http.StatusBadRequest, "invalid payload")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	{
		// validation
		if order.Task.ID != nil {
			err := models.NewError(http.StatusBadRequest, "id not allowed")
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}
	}

	u, err := strg.Write(ctx, actions.CREATE, &order)
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(u)
	if jsonerr != nil {
		err := models.NewError(http.StatusInternalServerError, "marshalling order failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(bs)
}

func handleGetOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleGetOrderByID")
	defer middleware.SpanStop(ctx, "handleGetOrderByID")

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
		err := models.NewError(http.StatusInternalServerError, "marshalling order failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func handleGetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleGetOrders")
	defer middleware.SpanStop(ctx, "handleGetOrders")

	us, err := strg.Read(ctx, actions.READALL, 0)
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(us)
	if jsonerr != nil {
		err := models.NewError(http.StatusInternalServerError, "marshalling order failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func handleGetOrderVersions(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleGetOrderVersions")
	defer middleware.SpanStop(ctx, "handleGetOrderVersions")

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
		err := models.NewError(http.StatusInternalServerError, "marshalling order failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func handleGetOrderSubTasks(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleGetOrderSubTasks")
	defer middleware.SpanStop(ctx, "handleGetOrderSubTasks")

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
		err := models.NewError(http.StatusInternalServerError, "marshalling order failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func handlePatchOrder(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handlePatchOrder")
	defer middleware.SpanStop(ctx, "handlePatchOrder")

	var order models.Order
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&order); err != nil {
		err := models.NewError(http.StatusBadRequest, "invalid payload")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	{
		// validation
		if order.Task.ID == nil {
			err := models.NewError(http.StatusBadRequest, "id must be provided")
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}
	}

	u, err := strg.Write(ctx, actions.UPDATE, &order)
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(u)
	if jsonerr != nil {
		err := models.NewError(http.StatusInternalServerError, "marshalling order failed")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func handleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleDeleteOrder")
	defer middleware.SpanStop(ctx, "handleDeleteOrder")

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

	_, err := strg.Write(ctx, action, &models.Order{Task: &models.Task{ID: utils.Ptr(uint(id))}})
	if err != nil {
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {

	strg = local.NewStorageOrder()

	http.HandleFunc("POST /order", handlePostOrder)
	http.HandleFunc("GET /order/{id}", handleGetOrderByID)
	http.HandleFunc("GET /orders", handleGetOrders)
	http.HandleFunc("GET /order/{id}/versions", handleGetOrderVersions)
	http.HandleFunc("GET /order/{id}/subtasks", handleGetOrderSubTasks)
	http.HandleFunc("PATCH /order", handlePatchOrder)
	http.HandleFunc("DELETE /order/{id}", handleDeleteOrder)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
