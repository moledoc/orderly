package response

import (
	"github.com/moledoc/orderly/internal/domain/order"
)

func (r *PostOrderResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *GetOrderByIDResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *GetOrdersResponse) GetOrders() []*order.Order {
	if r == nil {
		return nil
	}
	return r.Orders
}

////////////////

func (r *PatchOrderResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *PutDelegatedOrdersResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *PatchDelegatedOrdersResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *DeleteDelegatedOrdersResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *PutSitRepsResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *PatchSitRepsResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *DeleteSitRepsResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}
