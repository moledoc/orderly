package svc

import (
	"context"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/tests/api"
)

type OrderAPISvc struct { // NOTE: tests service layer methods
	Svc mgmtorder.ServiceMgmtOrderAPI
}

func NewOrderAPISvc() *OrderAPISvc {
	// TODO: local vs db
	return &OrderAPISvc{
		Svc: mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()),
	}
}

var (
	_ api.Order = (*OrderAPISvc)(nil)
)

func (api *OrderAPISvc) PostOrder(t *testing.T, ctx context.Context, req *request.PostOrderRequest) (*response.PostOrderResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PostOrder(ctx, req)
}

func (api *OrderAPISvc) GetOrderByID(t *testing.T, ctx context.Context, req *request.GetOrderByIDRequest) (*response.GetOrderByIDResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.GetOrderByID(ctx, req)
}

func (api *OrderAPISvc) GetOrders(t *testing.T, ctx context.Context, req *request.GetOrdersRequest) (*response.GetOrdersResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.GetOrders(ctx, req)
}

func (api *OrderAPISvc) PatchOrder(t *testing.T, ctx context.Context, req *request.PatchOrderRequest) (*response.PatchOrderResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PatchOrder(ctx, req)
}

func (api *OrderAPISvc) DeleteOrder(t *testing.T, ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.DeleteOrder(ctx, req)
}

////

func (api *OrderAPISvc) PutDelegatedOrders(t *testing.T, ctx context.Context, req *request.PutDelegatedOrdersRequest) (*response.PutDelegatedOrdersResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PutDelegatedOrders(ctx, req)
}

func (api *OrderAPISvc) PatchDelegatedOrders(t *testing.T, ctx context.Context, req *request.PatchDelegatedOrdersRequest) (*response.PatchDelegatedOrdersResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PatchDelegatedOrders(ctx, req)
}

func (api *OrderAPISvc) DeleteDelegatedOrders(t *testing.T, ctx context.Context, req *request.DeleteDelegatedOrdersRequest) (*response.DeleteDelegatedOrdersResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.DeleteDelegatedOrders(ctx, req)
}

////

func (api *OrderAPISvc) PutSitReps(t *testing.T, ctx context.Context, req *request.PutSitRepsRequest) (*response.PutSitRepsResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PutSitReps(ctx, req)
}

func (api *OrderAPISvc) PatchSitReps(t *testing.T, ctx context.Context, req *request.PatchSitRepsRequest) (*response.PatchSitRepsResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PatchSitReps(ctx, req)
}

func (api *OrderAPISvc) DeleteSitReps(t *testing.T, ctx context.Context, req *request.DeleteSitRepsRequest) (*response.DeleteSitRepsResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.DeleteSitReps(ctx, req)
}
