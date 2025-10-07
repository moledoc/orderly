package request

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
)

type PostOrderRequest struct {
	Order *order.Order `json:"order"`
}

type GetOrderByIDRequest struct {
	ID meta.ID `json:"id"`
}

type GetOrdersRequest struct{}

type GetOrderSubOrdersRequest struct {
	ID meta.ID `json:"id"`
}

type PatchOrderRequest struct {
	Order *order.Order `json:"order"`
}

type DeleteOrderRequest struct {
	ID meta.ID `json:"id"`
}

////////////////

type PutDelegatedTaskRequest struct {
	OrderID meta.ID     `json:"order_id"`
	Task    *order.Task `json:"task"`
}

type PatchDelegatedTaskRequest struct {
	OrderID meta.ID     `json:"order_id"`
	Task    *order.Task `json:"task"`
}

type DeleteDelegatedTaskRequest struct {
	OrderID         meta.ID `json:"order_id"`
	DelegatedTaskID meta.ID `json:"delegatedTask_id"`
}

////////////////

type PutSitRepRequest struct {
	OrderID meta.ID       `json:"order_id"`
	SitRep  *order.SitRep `json:"sitrep"`
}

type PatchSitRepRequest struct {
	OrderID meta.ID       `json:"order_id"`
	SitRep  *order.SitRep `json:"sitrep"`
}

type DeleteSitRepRequest struct {
	OrderID  meta.ID `json:"order_id"`
	SitRepID meta.ID `json:"sitrep_id"`
}
