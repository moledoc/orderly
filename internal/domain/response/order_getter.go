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

func (r *GetOrderSubOrdersResponse) GetSubOrders() []*order.Order {
	if r == nil {
		return nil
	}
	return r.SubOrders
}

////////////////

func (r *PatchOrderResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *PutDelegatedTaskResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *PatchDelegatedTaskResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *DeleteDelegatedTaskResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *PutSitRepResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *PatchSitRepResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *DeleteSitRepResponse) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}
