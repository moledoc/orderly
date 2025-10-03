package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
)

var (
	orderID    = "order_id"
	subtaskID  = "subtask_id"
	sitrepID   = "sitrep_id"
	hardDelete = "hard_delete"
	userID     = "user_id"
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

func RouteOrder(svc mgmtorder.ServiceMgmtOrderAPI) {
	mgmtordersvc = svc

	if mgmtordersvc == nil {
		panic("router: order management service is not initialized")
	}

	http.HandleFunc("POST /v1/mgmt/order", postOrder)
	http.HandleFunc(fmt.Sprintf("GET /v1/mgmt/order/{%v}", orderID), getOrderByID)
	http.HandleFunc("GET /v1/mgmt/orders", getOrders)
	http.HandleFunc(fmt.Sprintf("GET /v1/mgmt/order/{%v}/versions", orderID), getOrderVersions)
	http.HandleFunc(fmt.Sprintf("GET /v1/mgmt/order/{%v}/suborders", orderID), getOrderSubOrders)
	http.HandleFunc("PATCH /v1/mgmt/order", patchOrder)
	http.HandleFunc(fmt.Sprintf("DELETE /v1/mgmt/order/{%v}", orderID), deleteOrder)

	http.HandleFunc(fmt.Sprintf("PUT /v1/mgmt/order/{%v}/subtask", orderID), putSubTask)
	http.HandleFunc(fmt.Sprintf("PATCH /v1/mgmt/order/{%v}/subtask", orderID), patchSubTask)
	http.HandleFunc(fmt.Sprintf("DELETE /v1/mgmt/order/{%v}/subtask/{%v}", orderID, subtaskID), deleteSubTask)

	http.HandleFunc(fmt.Sprintf("PUT /v1/mgmt/order/{%v}/sitrep", orderID), putSitRep)
	http.HandleFunc(fmt.Sprintf("PATCH /v1/mgmt/order/{%v}/sitrep", orderID), patchSitRep)
	http.HandleFunc(fmt.Sprintf("DELETE /v1/mgmt/order/{%v}/sitrep/{%v}", orderID, sitrepID), deleteSitRep)
}

func RouteUser(svc mgmtuser.ServiceMgmtUserAPI) {
	mgmtusersvc = svc

	if mgmtusersvc == nil {
		panic("router: user management service is not initialized")
	}

	http.HandleFunc("POST /v1/mgmt/user", postUser)
	http.HandleFunc(fmt.Sprintf("GET /v1/mgmt/user/{%v}", userID), getUserByID)
	http.HandleFunc("GET /v1/mgmt/users", getUsers)
	http.HandleFunc(fmt.Sprintf("GET /v1/mgmt/user/{%v}/versions", userID), getUserVersions)
	http.HandleFunc(fmt.Sprintf("GET /v1/mgmt/user/{%v}/subordinates", userID), getUserSubOrdinates)
	http.HandleFunc("PATCH /v1/mgmt/user", patchUser)
	http.HandleFunc(fmt.Sprintf("DELETE /v1/mgmt/user/{%v}", userID), deleteUser)
}

func Route(svcs *Service) {
	RouteOrder(svcs.MgmtOrder)
	RouteUser(svcs.MgmtUser)
}
