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

func postOrder(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "postOrder")
	defer middleware.SpanStop(ctx, "postOrder")

	var req request.PostOrderRequest
	var resp *response.PostOrderResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		resp, err = mgmtordersvc.PostOrder(ctx, &req)
	}

	writeResponse(ctx, w, resp, err, http.StatusCreated)
}

func getOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "getOrderByID")
	defer middleware.SpanStop(ctx, "getOrderByID")

	resp, err := mgmtordersvc.GetOrderByID(ctx, &request.GetOrderByIDRequest{
		ID: utils.Ptr(meta.ID(r.PathValue(orderID))),
	})
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "getOrderByID")
	defer middleware.SpanStop(ctx, "getOrderByID")

	resp, err := mgmtordersvc.GetOrders(ctx, &request.GetOrdersRequest{})
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func getOrderVersions(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "getOrderVersions")
	defer middleware.SpanStop(ctx, "getOrderVersions")

	resp, err := mgmtordersvc.GetOrderVersions(ctx, &request.GetOrderVersionsRequest{
		ID: utils.Ptr(meta.ID(r.PathValue(orderID))),
	})
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func getOrderSubOrders(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "getOrderSubOrders")
	defer middleware.SpanStop(ctx, "getOrderSubOrders")

	resp, err := mgmtordersvc.GetOrderSubOrders(ctx, &request.GetOrderSubOrdersRequest{
		ID: utils.Ptr(meta.ID(r.PathValue(orderID))),
	})
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func patchOrder(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "patchOrder")
	defer middleware.SpanStop(ctx, "patchOrder")

	var req request.PatchOrderRequest
	var resp *response.PatchOrderResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		resp, err = mgmtordersvc.PatchOrder(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "deleteOrder")
	defer middleware.SpanStop(ctx, "deleteOrder")

	resp, err := mgmtordersvc.DeleteOrder(ctx, &request.DeleteOrderRequest{
		ID:   utils.Ptr(meta.ID(r.PathValue(orderID))),
		Hard: utils.Ptr(r.PathValue(hardDelete) == "true"),
	})
	writeResponse(ctx, w, resp, err, http.StatusNoContent)
}

////////////////

func putSubTask(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "putSubTask")
	defer middleware.SpanStop(ctx, "putSubTask")

	var req request.PutSubTaskRequest = request.PutSubTaskRequest{
		OrderID: utils.Ptr(meta.ID(r.PathValue(orderID))),
	}
	var resp *response.PutSubTaskResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		resp, err = mgmtordersvc.PutSubTask(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func patchSubTask(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "patchSubTask")
	defer middleware.SpanStop(ctx, "patchSubTask")

	var req request.PatchSubTaskRequest = request.PatchSubTaskRequest{
		OrderID: utils.Ptr(meta.ID(r.PathValue(orderID))),
	}
	var resp *response.PatchSubTaskResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		resp, err = mgmtordersvc.PatchSubTask(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func deleteSubTask(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "deleteSubTask")
	defer middleware.SpanStop(ctx, "deleteSubTask")

	resp, err := mgmtordersvc.DeleteSubTask(ctx, &request.DeleteSubTaskRequest{
		OrderID:   utils.Ptr(meta.ID(r.PathValue(orderID))),
		SubTaskID: utils.Ptr(meta.ID(r.PathValue(subtaskID))),
		Hard:      r.URL.Query().Get(hardDelete) == "true",
	})
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

////////////////

func putSitRep(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "putSitRep")
	defer middleware.SpanStop(ctx, "putSitRep")

	var req request.PutSitRepRequest = request.PutSitRepRequest{
		OrderID: utils.Ptr(meta.ID(r.PathValue(orderID))),
	}
	var resp *response.PutSitRepResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		resp, err = mgmtordersvc.PutSitRep(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func patchSitRep(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "patchSitRep")
	defer middleware.SpanStop(ctx, "patchSitRep")

	var req request.PatchSitRepRequest = request.PatchSitRepRequest{
		OrderID: utils.Ptr(meta.ID(r.PathValue(orderID))),
	}
	var resp *response.PatchSitRepResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		resp, err = mgmtordersvc.PatchSitRep(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func deleteSitRep(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "deleteSitRep")
	defer middleware.SpanStop(ctx, "deleteSitRep")

	resp, err := mgmtordersvc.DeleteSitRep(ctx, &request.DeleteSitRepRequest{
		OrderID:  utils.Ptr(meta.ID(r.PathValue(orderID))),
		SitRepID: utils.Ptr(meta.ID(r.PathValue(sitrepID))),
		Hard:     utils.Ptr(r.URL.Query().Get(hardDelete) == "true"),
	})

	writeResponse(ctx, w, resp, err, http.StatusOK)
}
