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

func TestOrderSvcPerformanceSuite(t *testing.T) {

	t.Run("OrderAPISvcPerformance", func(t *testing.T) {
		suite.Run(t, &OrderPerformanceSuite{
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

func (api *OrderAPISvc) PutDelegatedTasks(t *testing.T, ctx context.Context, req *request.PutDelegatedTasksRequest) (*response.PutDelegatedTasksResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PutDelegatedTasks(ctx, req)
}

func (api *OrderAPISvc) PatchDelegatedTasks(t *testing.T, ctx context.Context, req *request.PatchDelegatedTasksRequest) (*response.PatchDelegatedTasksResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PatchDelegatedTasks(ctx, req)
}

func (api *OrderAPISvc) DeleteDelegatedTasks(t *testing.T, ctx context.Context, req *request.DeleteDelegatedTasksRequest) (*response.DeleteDelegatedTasksResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.DeleteDelegatedTasks(ctx, req)
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
