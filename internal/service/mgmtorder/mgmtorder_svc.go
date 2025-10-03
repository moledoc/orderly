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
	GetOrderVersions(ctx context.Context, req *request.GetOrderVersionsRequest) (*response.GetOrderVersionsResponse, errwrap.Error)
	GetOrderSubOrders(ctx context.Context, req *request.GetOrderSubOrdersRequest) (*response.GetOrderSubOrdersResponse, errwrap.Error)
	PatchOrder(ctx context.Context, req *request.PatchOrderRequest) (*response.PatchOrderResponse, errwrap.Error)
	DeleteOrder(ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error)
	////
	PutSubTask(ctx context.Context, req *request.PutSubTaskRequest) (*response.PutSubTaskResponse, errwrap.Error)
	PatchSubTask(ctx context.Context, req *request.PatchSubTaskRequest) (*response.PatchSubTaskResponse, errwrap.Error)
	DeleteSubTask(ctx context.Context, req *request.DeleteSubTaskRequest) (*response.DeleteSubTaskResponse, errwrap.Error)
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
