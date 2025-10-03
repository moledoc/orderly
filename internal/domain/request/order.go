package request

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
)

type PostOrderRequest struct {
	Order *order.Order `json:"order"`
}

type GetOrderByIDRequest struct {
	ID *meta.ID `json:"id"`
}

type GetOrdersRequest struct{}

type GetOrderVersionsRequest struct {
	ID *meta.ID `json:"id"`
}

type GetOrderSubOrdersRequest struct {
	ID *meta.ID `json:"id"`
}

type PatchOrderRequest struct {
	Order *order.Order `json:"order"`
}

type DeleteOrderRequest struct {
	ID   *meta.ID `json:"id"`
	Hard *bool    `json:"hard"`
}

////////////////

type PutSubTaskRequest struct {
	OrderID *meta.ID    `json:"order_id"`
	Task    *order.Task `json:"task"`
}

type PatchSubTaskRequest struct {
	OrderID *meta.ID    `json:"order_id"`
	Task    *order.Task `json:"task"`
}

type DeleteSubTaskRequest struct {
	OrderID   *meta.ID `json:"order_id"`
	SubTaskID *meta.ID `json:"subtask_id"`
	Hard      bool     `json:"hard"`
}

////////////////

type PutSitRepRequest struct {
	OrderID *meta.ID      `json:"order_id"`
	SitRep  *order.SitRep `json:"sitrep"`
}

type PatchSitRepRequest struct {
	OrderID *meta.ID      `json:"order_id"`
	SitRep  *order.SitRep `json:"sitrep"`
}

type DeleteSitRepRequest struct {
	OrderID  *meta.ID `json:"order_id"`
	SitRepID *meta.ID `json:"sitrep_id"`
	Hard     *bool    `json:"hard"`
}
