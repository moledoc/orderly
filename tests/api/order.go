package api

import (
	"context"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
)

type Order interface {
	PostOrder(t *testing.T, ctx context.Context, req *request.PostOrderRequest) (*response.PostOrderResponse, errwrap.Error)
	GetOrderByID(t *testing.T, ctx context.Context, req *request.GetOrderByIDRequest) (*response.GetOrderByIDResponse, errwrap.Error)
	GetOrders(t *testing.T, ctx context.Context, req *request.GetOrdersRequest) (*response.GetOrdersResponse, errwrap.Error)
	GetOrderSubOrders(t *testing.T, ctx context.Context, req *request.GetOrderSubOrdersRequest) (*response.GetOrderSubOrdersResponse, errwrap.Error)
	PatchOrder(t *testing.T, ctx context.Context, req *request.PatchOrderRequest) (*response.PatchOrderResponse, errwrap.Error)
	DeleteOrder(t *testing.T, ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error)
	////
	PutDelegatedTask(t *testing.T, ctx context.Context, req *request.PutDelegatedTasksRequest) (*response.PutDelegatedTasksResponse, errwrap.Error)
	PatchDelegatedTask(t *testing.T, ctx context.Context, req *request.PatchDelegatedTasksRequest) (*response.PatchDelegatedTasksResponse, errwrap.Error)
	DeleteDelegatedTask(t *testing.T, ctx context.Context, req *request.DeleteDelegatedTasksRequest) (*response.DeleteDelegatedTasksResponse, errwrap.Error)
	////
	PutSitRep(t *testing.T, ctx context.Context, req *request.PutSitRepsRequest) (*response.PutSitRepsResponse, errwrap.Error)
	PatchSitRep(t *testing.T, ctx context.Context, req *request.PatchSitRepsRequest) (*response.PatchSitRepsResponse, errwrap.Error)
	DeleteSitRep(t *testing.T, ctx context.Context, req *request.DeleteSitRepsRequest) (*response.DeleteSitRepsResponse, errwrap.Error)
}
