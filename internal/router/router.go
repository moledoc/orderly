package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/moledoc/orderly/pkg/consts"
)

var (
	orderID         = "order_id"
	delegatedTaskID = "delegated_task_id"
	sitrepID        = "sitrep_id"
	userID          = "user_id"
)

var (
	onceRouteUser  sync.Once
	onceRouteOrder sync.Once
)

func decodeBody(ctx context.Context, r *http.Request, req any) errwrap.Error {
	middleware.SpanStart(ctx, "DecodeBody")
	defer middleware.SpanStop(ctx, "DecodeBody")

	var err errwrap.Error
	if errj := json.NewDecoder(r.Body).Decode(req); errj != nil {
		err = errwrap.NewError(http.StatusBadRequest, "invalid payload: %s", errj)
	}
	return err
}

func writeResponse(ctx context.Context, w http.ResponseWriter, resp any, err errwrap.Error, successCode int) {
	middleware.SpanStart(ctx, "WriteResponse")
	defer middleware.SpanStop(ctx, "WriteResponse")

	if err != nil {

		var traceID string
		ctxTraceID := ctx.Value(consts.CtxKeyTrace)
		if ctxTraceID != nil {
			traceID = ctxTraceID.(string)
		}
		err.SetTraceID(traceID)

		w.WriteHeader(err.GetStatusCode())
		w.Write([]byte(err.String()))
		return
	}

	bs, jsonerr := json.Marshal(resp)
	if jsonerr != nil {
		err := errwrap.NewError(http.StatusInternalServerError, "marshalling resp failed")
		w.WriteHeader(err.GetStatusCode())
		w.Write([]byte(err.String()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(successCode)
	w.Write(bs)
}

type Service struct {
	MgmtOrder mgmtorder.ServiceMgmtOrderAPI
	MgmtUser  mgmtuser.ServiceMgmtUserAPI
}

var (
	mgmtordersvc mgmtorder.ServiceMgmtOrderAPI = nil
	mgmtusersvc  mgmtuser.ServiceMgmtUserAPI   = nil
)

func RouteOrder(svc mgmtorder.ServiceMgmtOrderAPI) *http.ServeMux {
	mgmtordersvc = svc

	if mgmtordersvc == nil {
		panic("router: order management service is not initialized")
	}

	onceRouteOrder.Do(func() {
		http.HandleFunc("POST /v1/mgmt/order", postOrder)
		http.HandleFunc(fmt.Sprintf("GET /v1/mgmt/order/{%v}", orderID), getOrderByID)
		http.HandleFunc("GET /v1/mgmt/orders", getOrders)
		http.HandleFunc(fmt.Sprintf("GET /v1/mgmt/order/{%v}/suborders", orderID), getOrderSubOrders)
		http.HandleFunc("PATCH /v1/mgmt/order", patchOrder)
		http.HandleFunc(fmt.Sprintf("DELETE /v1/mgmt/order/{%v}", orderID), deleteOrder)

		http.HandleFunc(fmt.Sprintf("PUT /v1/mgmt/order/{%v}/delegated_task", orderID), putDelegatedTasks)
		http.HandleFunc(fmt.Sprintf("PATCH /v1/mgmt/order/{%v}/delegated_task", orderID), patchDelegatedTasks)
		http.HandleFunc(fmt.Sprintf("DELETE /v1/mgmt/order/{%v}/delegated_task", orderID), deleteDelegatedTasks)

		http.HandleFunc(fmt.Sprintf("PUT /v1/mgmt/order/{%v}/sitrep", orderID), putSitReps)
		http.HandleFunc(fmt.Sprintf("PATCH /v1/mgmt/order/{%v}/sitrep", orderID), patchSitReps)
		http.HandleFunc(fmt.Sprintf("DELETE /v1/mgmt/order/{%v}/sitrep", orderID), deleteSitReps)

		// NOTE: handle empty ids
		http.HandleFunc("GET /v1/mgmt/order/", getOrderByID)
		http.HandleFunc("GET /v1/mgmt/order/suborders", getOrderSubOrders)
		http.HandleFunc("DELETE /v1/mgmt/order/", deleteOrder)

		http.HandleFunc("PUT /v1/mgmt/order/delegated_task", putDelegatedTasks)
		http.HandleFunc("PATCH /v1/mgmt/order/delegated_task", patchDelegatedTasks)
		http.HandleFunc("DELETE /v1/mgmt/order/delegated_task", deleteDelegatedTasks)

		http.HandleFunc("PUT /v1/mgmt/order/sitrep", putSitReps)
		http.HandleFunc("PATCH /v1/mgmt/order/sitrep", patchSitReps)
		http.HandleFunc("DELETE /v1/mgmt/order/sitrep", deleteSitReps)
	})

	return http.DefaultServeMux
}

func RouteUser(svc mgmtuser.ServiceMgmtUserAPI) *http.ServeMux {
	mgmtusersvc = svc

	if mgmtusersvc == nil {
		panic("router: user management service is not initialized")
	}

	onceRouteUser.Do(func() {
		http.HandleFunc("POST /v1/mgmt/user", handlePostUser)
		http.HandleFunc(fmt.Sprintf("GET /v1/mgmt/user/{%v}", userID), handleGetUserByID)
		http.HandleFunc("GET /v1/mgmt/users", handleGetUsers)
		http.HandleFunc(fmt.Sprintf("GET /v1/mgmt/user/{%v}/subordinates", userID), handleGetUserSubOrdinates)
		http.HandleFunc("PATCH /v1/mgmt/user", handlePatchUser)
		http.HandleFunc(fmt.Sprintf("DELETE /v1/mgmt/user/{%v}", userID), handleDeleteUser)

		// NOTE: handle empty ids
		http.HandleFunc("GET /v1/mgmt/user/", handleGetUserByID)
		http.HandleFunc("GET /v1/mgmt/user/subordinates", handleGetUserSubOrdinates)
		http.HandleFunc("DELETE /v1/mgmt/user/", handleDeleteUser)
	})

	return http.DefaultServeMux
}

func Route(svcs *Service) {
	RouteOrder(svcs.MgmtOrder)
	RouteUser(svcs.MgmtUser)
}
