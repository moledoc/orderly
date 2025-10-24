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

type PutDelegatedOrdersResponse struct {
	Order *order.Order `json:"order"`
}

type PatchDelegatedOrdersResponse struct {
	Order *order.Order `json:"order"`
}

type DeleteDelegatedOrdersResponse struct {
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
