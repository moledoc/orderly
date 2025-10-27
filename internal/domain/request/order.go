package request

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
)

type PostOrderRequest struct {
	Order *order.Order `json:"order,omitempty"`
}

type GetOrderByIDRequest struct {
	ID meta.ID `json:"id,omitempty"`
}

type GetOrdersRequest struct {
	ParentOrderID meta.ID `json:"parent_order_id,omitempty"`
	AccountableID meta.ID `json:"accountable_id,omitempty"`
}

type PatchOrderRequest struct {
	Order *order.Order `json:"order,omitempty"`
}

type DeleteOrderRequest struct {
	ID meta.ID `json:"id,omitempty"`
}

////////////////

type PutDelegatedOrdersRequest struct {
	OrderID meta.ID        `json:"order_id,omitempty"`
	Orders  []*order.Order `json:"tasks,omitempty"`
}

type PatchDelegatedOrdersRequest struct {
	OrderID meta.ID        `json:"order_id,omitempty"`
	Orders  []*order.Order `json:"tasks,omitempty"`
}

type DeleteDelegatedOrdersRequest struct {
	OrderID           meta.ID   `json:"order_id,omitempty"`
	DelegatedOrderIDs []meta.ID `json:"delegated_task_ids,omitempty"`
}

////////////////

type PutSitRepsRequest struct {
	OrderID meta.ID         `json:"order_id,omitempty"`
	SitReps []*order.SitRep `json:"sitreps,omitempty"`
}

type PatchSitRepsRequest struct {
	OrderID meta.ID         `json:"order_id,omitempty"`
	SitReps []*order.SitRep `json:"sitreps,omitempty"`
}

type DeleteSitRepsRequest struct {
	OrderID   meta.ID   `json:"order_id,omitempty"`
	SitRepIDs []meta.ID `json:"sitrep_ids,omitempty"`
}
