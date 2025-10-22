package request

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/user"
)

type PostOrderRequest struct {
	Order *order.Order `json:"order,omitempty"`
}

type GetOrderByIDRequest struct {
	ID meta.ID `json:"id,omitempty"`
}

type GetOrdersRequest struct {
	ParentOrderID meta.ID    `json:"parent_order_id,omitempty"`
	Accountable   user.Email `json:"accountable,omitempty"`
}

type GetOrderSubOrdersRequest struct {
	ID meta.ID `json:"id,omitempty"`
}

type PatchOrderRequest struct {
	Order *order.Order `json:"order,omitempty"`
}

type DeleteOrderRequest struct {
	ID meta.ID `json:"id,omitempty"`
}

////////////////

type PutDelegatedTasksRequest struct {
	OrderID meta.ID       `json:"order_id,omitempty"`
	Tasks   []*order.Task `json:"tasks,omitempty"`
}

type PatchDelegatedTasksRequest struct {
	OrderID meta.ID       `json:"order_id,omitempty"`
	Tasks   []*order.Task `json:"tasks,omitempty"`
}

type DeleteDelegatedTasksRequest struct {
	OrderID          meta.ID   `json:"order_id,omitempty"`
	DelegatedTaskIDs []meta.ID `json:"delegated_task_ids,omitempty"`
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
