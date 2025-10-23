package response

import (
	"github.com/moledoc/orderly/internal/domain/order"
)

type PostOrderResponse struct {
	Order *order.Order `json:"order"`
}

type GetOrderByIDResponse struct {
	Order *order.Order `json:"order"`
}

type GetOrdersResponse struct {
	Orders []*order.Order `json:"orders"`
}

type PatchOrderResponse struct {
	Order *order.Order `json:"order"`
}

type DeleteOrderResponse struct{}

////////////////

type PutDelegatedTasksResponse struct {
	Order *order.Order `json:"order"`
}

type PatchDelegatedTasksResponse struct {
	Order *order.Order `json:"order"`
}

type DeleteDelegatedTasksResponse struct {
	Order *order.Order `json:"order"`
}

////////////////

type PutSitRepsResponse struct {
	Order *order.Order `json:"order"`
}

type PatchSitRepsResponse struct {
	Order *order.Order `json:"order"`
}

type DeleteSitRepsResponse struct {
	Order *order.Order `json:"order"`
}

////////////////

type GetOrdersByAccountableResponse struct {
	Orders []*order.Order `json:"orders"`
}
