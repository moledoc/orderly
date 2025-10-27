package router

import (
	"context"
	"fmt"
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

	if err == nil {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/order/%v", resp.GetOrder().GetID()))
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

	middleware.SpanStart(ctx, "getOrders")
	defer middleware.SpanStop(ctx, "getOrders")

	req := &request.GetOrdersRequest{
		ParentOrderID: meta.ID(r.URL.Query().Get("parent_order_id")),
		AccountableID: meta.ID(r.URL.Query().Get("accountable_id")),
	}

	middleware.SpanLog(ctx, "GetOrdersRequest", req)
	resp, err := mgmtordersvc.GetOrders(ctx, req)
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

func putDelegatedOrders(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "putDelegatedOrder")
	defer middleware.SpanStop(ctx, "putDelegatedOrder")

	var req request.PutDelegatedOrdersRequest = request.PutDelegatedOrdersRequest{
		OrderID: meta.ID(r.PathValue(orderID)),
	}
	var resp *response.PutDelegatedOrdersResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "PutDelegatedOrderRequest", &req)
		resp, err = mgmtordersvc.PutDelegatedOrders(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func patchDelegatedOrders(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "patchDelegatedOrder")
	defer middleware.SpanStop(ctx, "patchDelegatedOrder")

	var req request.PatchDelegatedOrdersRequest = request.PatchDelegatedOrdersRequest{
		OrderID: meta.ID(r.PathValue(orderID)),
	}
	var resp *response.PatchDelegatedOrdersResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "PatchDelegatedOrderRequest", &req)
		resp, err = mgmtordersvc.PatchDelegatedOrders(ctx, &req)
	}
	writeResponse(ctx, w, resp, err, http.StatusOK)
}

func deleteDelegatedOrders(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.AddTraceToCtxFromWriter(context.Background(), w)
	defer func() { go middleware.SpanFlushTrace(ctx) }()

	middleware.SpanStart(ctx, "deleteDelegatedOrder")
	defer middleware.SpanStop(ctx, "deleteDelegatedOrder")

	req := &request.DeleteDelegatedOrdersRequest{
		OrderID: meta.ID(r.PathValue(orderID)),
	}
	var resp *response.DeleteDelegatedOrdersResponse
	var err errwrap.Error

	err = decodeBody(ctx, r, &req)
	if err == nil {
		middleware.SpanLog(ctx, "DeleteDelegatedOrderRequest", req)
		resp, err = mgmtordersvc.DeleteDelegatedOrders(ctx, req)
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
