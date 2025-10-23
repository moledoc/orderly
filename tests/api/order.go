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
	PatchOrder(t *testing.T, ctx context.Context, req *request.PatchOrderRequest) (*response.PatchOrderResponse, errwrap.Error)
	DeleteOrder(t *testing.T, ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error)
	////
	PutDelegatedTasks(t *testing.T, ctx context.Context, req *request.PutDelegatedTasksRequest) (*response.PutDelegatedTasksResponse, errwrap.Error)
	PatchDelegatedTasks(t *testing.T, ctx context.Context, req *request.PatchDelegatedTasksRequest) (*response.PatchDelegatedTasksResponse, errwrap.Error)
	DeleteDelegatedTasks(t *testing.T, ctx context.Context, req *request.DeleteDelegatedTasksRequest) (*response.DeleteDelegatedTasksResponse, errwrap.Error)
	////
	PutSitReps(t *testing.T, ctx context.Context, req *request.PutSitRepsRequest) (*response.PutSitRepsResponse, errwrap.Error)
	PatchSitReps(t *testing.T, ctx context.Context, req *request.PatchSitRepsRequest) (*response.PatchSitRepsResponse, errwrap.Error)
	DeleteSitReps(t *testing.T, ctx context.Context, req *request.DeleteSitRepsRequest) (*response.DeleteSitRepsResponse, errwrap.Error)
}
