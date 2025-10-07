package tests

import (
	"context"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/tests/api"
	"github.com/stretchr/testify/suite"
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

func TestOrderSvcSuite(t *testing.T) {
	t.Run("OrderAPISvc", func(t *testing.T) {
		suite.Run(t, &OrderSuite{
			API: NewOrderAPISvc(),
		})
	})
}

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

func (api *OrderAPISvc) GetOrderSubOrders(t *testing.T, ctx context.Context, req *request.GetOrderSubOrdersRequest) (*response.GetOrderSubOrdersResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.GetOrderSubOrders(ctx, req)
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

func (api *OrderAPISvc) PutDelegatedTask(t *testing.T, ctx context.Context, req *request.PutDelegatedTaskRequest) (*response.PutDelegatedTaskResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PutDelegatedTask(ctx, req)
}

func (api *OrderAPISvc) PatchDelegatedTask(t *testing.T, ctx context.Context, req *request.PatchDelegatedTaskRequest) (*response.PatchDelegatedTaskResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PatchDelegatedTask(ctx, req)
}

func (api *OrderAPISvc) DeleteDelegatedTask(t *testing.T, ctx context.Context, req *request.DeleteDelegatedTaskRequest) (*response.DeleteDelegatedTaskResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.DeleteDelegatedTask(ctx, req)
}

////

func (api *OrderAPISvc) PutSitRep(t *testing.T, ctx context.Context, req *request.PutSitRepRequest) (*response.PutSitRepResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PutSitRep(ctx, req)
}

func (api *OrderAPISvc) PatchSitRep(t *testing.T, ctx context.Context, req *request.PatchSitRepRequest) (*response.PatchSitRepResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PatchSitRep(ctx, req)
}

func (api *OrderAPISvc) DeleteSitRep(t *testing.T, ctx context.Context, req *request.DeleteSitRepRequest) (*response.DeleteSitRepResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.DeleteSitRep(ctx, req)
}
