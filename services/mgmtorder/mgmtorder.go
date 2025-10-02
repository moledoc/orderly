package mgmtorder

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

	{ // TODO: move validation to _api file
		// validation
		if order.Task.ID != nil {
			err := models.NewError(http.StatusBadRequest, "id not allowed")
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}
	}

	resp, err := HandlePostOrder(ctx, &order)

	common.WriteResponse(ctx, w, resp, err, http.StatusCreated)
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

	resp, err := HandleGetOrderByID(ctx, uint(id))
	common.WriteResponse(ctx, w, resp, err, http.StatusOK)

}

func handleGetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleGetOrderByID")
	defer middleware.SpanStop(ctx, "handleGetOrderByID")

	resp, err := HandleGetOrders(ctx)
	common.WriteResponse(ctx, w, resp, err, http.StatusOK)
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

	resp, err := HandleGetOrderVersions(ctx, uint(id))
	common.WriteResponse(ctx, w, resp, err, http.StatusOK)
}

func handleGetOrderSubOrders(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "handleGetOrderSubOrders")
	defer middleware.SpanStop(ctx, "handleGetOrderSubOrders")

	id, errAtoi := strconv.ParseUint(r.PathValue("id"), 10, 0)
	if errAtoi != nil {
		err := models.NewError(http.StatusBadRequest, "invalid id")
		w.WriteHeader(err.StatusCode())
		w.Write([]byte(err.String()))
		return
	}

	resp, err := HandleGetOrderSubOrders(ctx, uint(id))
	common.WriteResponse(ctx, w, resp, err, http.StatusOK)
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

	{ // TODO: move validation to _api
		// validation
		if order.Task.ID == nil {
			err := models.NewError(http.StatusBadRequest, "id must be provided")
			w.WriteHeader(err.StatusCode())
			w.Write([]byte(err.String()))
			return
		}
	}

	resp, err := HandlePatchOrder(ctx, &order)

	common.WriteResponse(ctx, w, resp, err, http.StatusOK)
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

	order := &models.Order{
		Task: &models.Task{
			ID: utils.Ptr(uint(id)),
		},
	}
	var resp *models.Order
	var err models.IError

	deleteType := r.URL.Query().Get("type")
	if deleteType == "hard" {
		resp, err = HandleDeleteOrderHard(ctx, order)
	} else {
		resp, err = HandleDeleteOrderSoft(ctx, order)

	}

	common.WriteResponse(ctx, w, resp, err, http.StatusNoContent)
}

func New() {
	strg = local.NewStorageOrder()

	http.HandleFunc("POST /order", handlePostOrder)
	http.HandleFunc("GET /order/{id}", handleGetOrderByID)
	http.HandleFunc("GET /orders", handleGetOrders)
	http.HandleFunc("GET /order/{id}/versions", handleGetOrderVersions)
	http.HandleFunc("GET /order/{id}/suborders", handleGetOrderSubOrders)
	http.HandleFunc("PATCH /order", handlePatchOrder)
	http.HandleFunc("DELETE /order/{id}", handleDeleteOrder)
}
