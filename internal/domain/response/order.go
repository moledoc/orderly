package response

import "github.com/moledoc/orderly/internal/domain/order"

type PostOrderResponse struct {
	Order *order.Order `json:"order"`
}

type GetOrderByIDResponse struct {
	Order *order.Order `json:"order"`
}

type GetOrdersResponse struct {
	Orders []*order.Order `json:"orders"`
}

type GetOrderSubOrdersResponse struct {
	SubOrders []*order.Order `json:"sub_orders"`
}

type PatchOrderResponse struct {
	Order *order.Order `json:"order"`
}

type DeleteOrderResponse struct{}

////////////////

type PutDelegatedTaskResponse struct {
	Order *order.Order `json:"order"`
}

type PatchDelegatedTaskResponse struct {
	Order *order.Order `json:"order"`
}

type DeleteDelegatedTaskResponse struct {
	Order *order.Order `json:"order"`
}

////////////////

type PutSitRepResponse struct {
	Order *order.Order `json:"order"`
}

type PatchSitRepResponse struct {
	Order *order.Order `json:"order"`
}

type DeleteSitRepResponse struct {
	Order *order.Order `json:"order"`
}
