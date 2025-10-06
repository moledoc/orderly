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
		ID: utils.Ptr(meta.ID(r.PathValue(orderID))),
	})
	writeResponse(ctx, w, resp, err, http.StatusNoContent)
}

////////////////

func putDelegatedTask(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "putDelegatedTask")
	defer middleware.SpanStop(ctx, "putDelegatedTask")

	var req request.PutDelegatedTaskRequest = request.PutDelegatedTaskRequest{
		OrderID: utils.Ptr(meta.ID(r.PathValue(orderID))),
	}
	var resp *response.PutDelegatedTaskResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		resp, err = mgmtordersvc.PutDelegatedTask(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func patchDelegatedTask(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "patchDelegatedTask")
	defer middleware.SpanStop(ctx, "patchDelegatedTask")

	var req request.PatchDelegatedTaskRequest = request.PatchDelegatedTaskRequest{
		OrderID: utils.Ptr(meta.ID(r.PathValue(orderID))),
	}
	var resp *response.PatchDelegatedTaskResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		resp, err = mgmtordersvc.PatchDelegatedTask(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func deleteDelegatedTask(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTrace(context.Background(), w)
	defer middleware.SpanFlushTrace(ctx)

	middleware.SpanStart(ctx, "deleteDelegatedTask")
	defer middleware.SpanStop(ctx, "deleteDelegatedTask")

	resp, err := mgmtordersvc.DeleteDelegatedTask(ctx, &request.DeleteDelegatedTaskRequest{
		OrderID:         utils.Ptr(meta.ID(r.PathValue(orderID))),
		DelegatedTaskID: utils.Ptr(meta.ID(r.PathValue(delegatedTaskID))),
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
	})

	writeResponse(ctx, w, resp, err, http.StatusOK)
}
