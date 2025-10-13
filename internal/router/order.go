package router

import (
	"context"
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/middleware"
)

func postOrder(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "postOrder")
	defer middleware.SpanStop(ctx, "postOrder")

	var req request.PostOrderRequest
	var resp *response.PostOrderResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "PostOrderRequest", &req)
		resp, err = mgmtordersvc.PostOrder(ctx, &req)
	}

	writeResponse(ctx, w, resp, err, http.StatusCreated)
}

func getOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "getOrderByID")
	defer middleware.SpanStop(ctx, "getOrderByID")

	req := &request.GetOrderByIDRequest{
		ID: meta.ID(r.PathValue(orderID)),
	}
	middleware.SpanLog(ctx, "GetOrderByIDRequest", req)
	resp, err := mgmtordersvc.GetOrderByID(ctx, req)
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "getOrderByID")
	defer middleware.SpanStop(ctx, "getOrderByID")

	req := &request.GetOrdersRequest{}
	middleware.SpanLog(ctx, "GetOrdersRequest", req)
	resp, err := mgmtordersvc.GetOrders(ctx, req)
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func getOrderSubOrders(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "getOrderSubOrders")
	defer middleware.SpanStop(ctx, "getOrderSubOrders")

	req := &request.GetOrderSubOrdersRequest{
		ID: meta.ID(r.PathValue(orderID)),
	}
	middleware.SpanLog(ctx, "GetOrderSubOrdersRequest", req)
	resp, err := mgmtordersvc.GetOrderSubOrders(ctx, req)
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func patchOrder(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "patchOrder")
	defer middleware.SpanStop(ctx, "patchOrder")

	var req request.PatchOrderRequest
	var resp *response.PatchOrderResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "PatchOrderRequest", &req)
		resp, err = mgmtordersvc.PatchOrder(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "deleteOrder")
	defer middleware.SpanStop(ctx, "deleteOrder")

	req := &request.DeleteOrderRequest{
		ID: meta.ID(r.PathValue(orderID)),
	}
	middleware.SpanLog(ctx, "DeleteOrderRequest", req)
	resp, err := mgmtordersvc.DeleteOrder(ctx, req)
	writeResponse(ctx, w, resp, err, http.StatusNoContent)
}

////////////////

func putDelegatedTasks(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "putDelegatedTask")
	defer middleware.SpanStop(ctx, "putDelegatedTask")

	var req request.PutDelegatedTasksRequest = request.PutDelegatedTasksRequest{
		OrderID: meta.ID(r.PathValue(orderID)),
	}
	var resp *response.PutDelegatedTasksResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "PutDelegatedTaskRequest", &req)
		resp, err = mgmtordersvc.PutDelegatedTasks(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func patchDelegatedTasks(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "patchDelegatedTask")
	defer middleware.SpanStop(ctx, "patchDelegatedTask")

	var req request.PatchDelegatedTasksRequest = request.PatchDelegatedTasksRequest{
		OrderID: meta.ID(r.PathValue(orderID)),
	}
	var resp *response.PatchDelegatedTasksResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "PatchDelegatedTaskRequest", &req)
		resp, err = mgmtordersvc.PatchDelegatedTasks(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func deleteDelegatedTasks(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "deleteDelegatedTask")
	defer middleware.SpanStop(ctx, "deleteDelegatedTask")

	req := &request.DeleteDelegatedTasksRequest{
		OrderID: meta.ID(r.PathValue(orderID)),
	}
	var resp *response.DeleteDelegatedTasksResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "DeleteDelegatedTaskRequest", req)
		resp, err = mgmtordersvc.DeleteDelegatedTasks(ctx, req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

////////////////

func putSitReps(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "putSitRep")
	defer middleware.SpanStop(ctx, "putSitRep")

	var req request.PutSitRepsRequest = request.PutSitRepsRequest{
		OrderID: meta.ID(r.PathValue(orderID)),
	}
	var resp *response.PutSitRepsResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "PutSitRepRequest", &req)
		resp, err = mgmtordersvc.PutSitReps(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func patchSitReps(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "patchSitRep")
	defer middleware.SpanStop(ctx, "patchSitRep")

	var req request.PatchSitRepsRequest = request.PatchSitRepsRequest{
		OrderID: meta.ID(r.PathValue(orderID)),
	}
	var resp *response.PatchSitRepsResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "PatchSitRepRequest", &req)
		resp, err = mgmtordersvc.PatchSitReps(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func deleteSitReps(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "deleteSitRep")
	defer middleware.SpanStop(ctx, "deleteSitRep")

	req := &request.DeleteSitRepsRequest{
		OrderID: meta.ID(r.PathValue(orderID)),
	}
	var resp *response.DeleteSitRepsResponse
	var err errwrap.Error
	err = decodeBody(ctx, r, req)
	if err == nil {
		middleware.SpanLog(ctx, "DeleteSitRepRequest", req)
		resp, err = mgmtordersvc.DeleteSitReps(ctx, req)
	}

	writeResponse(ctx, w, resp, err, http.StatusOK)
}
