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
	OrderID meta.ID       `json:"order_id"`
	Tasks   []*order.Task `json:"tasks"`
}

type PatchDelegatedTaskRequest struct {
	OrderID meta.ID       `json:"order_id"`
	Tasks   []*order.Task `json:"tasks"`
}

type DeleteDelegatedTaskRequest struct {
	OrderID          meta.ID   `json:"order_id"`
	DelegatedTaskIDs []meta.ID `json:"delegated_task_ids"`
}

////////////////

type PutSitRepRequest struct {
	OrderID meta.ID         `json:"order_id"`
	SitReps []*order.SitRep `json:"sitreps"`
}

type PatchSitRepRequest struct {
	OrderID meta.ID         `json:"order_id"`
	SitReps []*order.SitRep `json:"sitreps"`
}

type DeleteSitRepRequest struct {
	OrderID   meta.ID   `json:"order_id"`
	SitRepIDs []meta.ID `json:"sitrep_ids"`
}
