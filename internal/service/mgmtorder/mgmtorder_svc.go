package mgmtorder

import (
	"context"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/repository"
)

type ServiceMgmtOrderAPI interface {
	PostOrder(ctx context.Context, req *request.PostOrderRequest) (*response.PostOrderResponse, errwrap.Error)
	GetOrderByID(ctx context.Context, req *request.GetOrderByIDRequest) (*response.GetOrderByIDResponse, errwrap.Error)
	GetOrders(ctx context.Context, req *request.GetOrdersRequest) (*response.GetOrdersResponse, errwrap.Error)
	GetOrderSubOrders(ctx context.Context, req *request.GetOrderSubOrdersRequest) (*response.GetOrderSubOrdersResponse, errwrap.Error)
	PatchOrder(ctx context.Context, req *request.PatchOrderRequest) (*response.PatchOrderResponse, errwrap.Error)
	DeleteOrder(ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error)
	////
	PutDelegatedTask(ctx context.Context, req *request.PutDelegatedTaskRequest) (*response.PutDelegatedTaskResponse, errwrap.Error)
	PatchDelegatedTask(ctx context.Context, req *request.PatchDelegatedTaskRequest) (*response.PatchDelegatedTaskResponse, errwrap.Error)
	DeleteDelegatedTask(ctx context.Context, req *request.DeleteDelegatedTaskRequest) (*response.DeleteDelegatedTaskResponse, errwrap.Error)
	////
	PutSitRep(ctx context.Context, req *request.PutSitRepRequest) (*response.PutSitRepResponse, errwrap.Error)
	PatchSitRep(ctx context.Context, req *request.PatchSitRepRequest) (*response.PatchSitRepResponse, errwrap.Error)
	DeleteSitRep(ctx context.Context, req *request.DeleteSitRepRequest) (*response.DeleteSitRepResponse, errwrap.Error)
}

type serviceMgmtOrder struct {
	Repository repository.RepositoryOrderAPI
}

var (
	_ ServiceMgmtOrderAPI = (*serviceMgmtOrder)(nil)
)

func NewServiceMgmtOrder(repo repository.RepositoryOrderAPI) ServiceMgmtOrderAPI {
	return &serviceMgmtOrder{
		Repository: repo,
	}
}
