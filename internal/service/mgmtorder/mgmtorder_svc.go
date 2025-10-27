package mgmtorder

import (
	"context"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/internal/repository"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/moledoc/orderly/pkg/utils"
)

type ServiceMgmtOrderAPI interface {
	GetRootOrder(context.Context) *order.Order
	////
	PostOrder(ctx context.Context, req *request.PostOrderRequest) (*response.PostOrderResponse, errwrap.Error)
	GetOrderByID(ctx context.Context, req *request.GetOrderByIDRequest) (*response.GetOrderByIDResponse, errwrap.Error)
	GetOrders(ctx context.Context, req *request.GetOrdersRequest) (*response.GetOrdersResponse, errwrap.Error)
	PatchOrder(ctx context.Context, req *request.PatchOrderRequest) (*response.PatchOrderResponse, errwrap.Error)
	DeleteOrder(ctx context.Context, req *request.DeleteOrderRequest) (*response.DeleteOrderResponse, errwrap.Error)
	////
	PutDelegatedOrders(ctx context.Context, req *request.PutDelegatedOrdersRequest) (*response.PutDelegatedOrdersResponse, errwrap.Error)
	PatchDelegatedOrders(ctx context.Context, req *request.PatchDelegatedOrdersRequest) (*response.PatchDelegatedOrdersResponse, errwrap.Error)
	DeleteDelegatedOrders(ctx context.Context, req *request.DeleteDelegatedOrdersRequest) (*response.DeleteDelegatedOrdersResponse, errwrap.Error)
	////
	PutSitReps(ctx context.Context, req *request.PutSitRepsRequest) (*response.PutSitRepsResponse, errwrap.Error)
	PatchSitReps(ctx context.Context, req *request.PatchSitRepsRequest) (*response.PatchSitRepsResponse, errwrap.Error)
	DeleteSitReps(ctx context.Context, req *request.DeleteSitRepsRequest) (*response.DeleteSitRepsResponse, errwrap.Error)
}

type serviceMgmtOrder struct {
	Repository repository.RepositoryOrderAPI
	RootOrder  *order.Order
}

var (
	_   ServiceMgmtOrderAPI = (*serviceMgmtOrder)(nil)
	svc ServiceMgmtOrderAPI = nil
)

func postRootOrder(ctx context.Context, repo repository.RepositoryOrderAPI) (*order.Order, errwrap.Error) {
	now := time.Now().UTC()
	id := meta.NewID()
	order := &order.Order{
		ID:            id,
		ParentOrderID: id,
		AccountableID: mgmtuser.RootUserID,
		Objective:     "Root Order",
		State:         utils.Ptr(order.InProgress),
		Deadline:      time.Now().UTC().Add(100 * 365 * 24 * time.Hour),
		Meta: &meta.Meta{
			Version: 1,
			Created: now,
			Updated: now,
		},
	}

	o, err := repo.Write(ctx, order)
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return o, nil
}

func GetServiceMgmtOrder() ServiceMgmtOrderAPI {
	return svc
}

func NewServiceMgmtOrder(repo repository.RepositoryOrderAPI) ServiceMgmtOrderAPI {
	rootOrder, err := postRootOrder(context.Background(), repo)
	if err != nil {
		panic(err)
	}
	svc = &serviceMgmtOrder{
		RootOrder:  rootOrder,
		Repository: repo,
	}
	return svc
}
