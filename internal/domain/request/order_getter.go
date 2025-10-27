package request

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
)

func (r *PostOrderRequest) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *GetOrderByIDRequest) GetID() meta.ID {
	if r == nil {
		return ""
	}
	return r.ID
}

////////////////

func (r *GetOrdersRequest) GetParentOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.ParentOrderID
}

func (r *GetOrdersRequest) GetAccountableID() meta.ID {
	if r == nil {
		return ""
	}
	return r.AccountableID
}

////////////////

func (r *PatchOrderRequest) GetOrder() *order.Order {
	if r == nil {
		return nil
	}
	return r.Order
}

////////////////

func (r *DeleteOrderRequest) GetID() meta.ID {
	if r == nil {
		return ""
	}
	return r.ID
}

////////////////

func (r *PutDelegatedOrdersRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PutDelegatedOrdersRequest) GetOrders() []*order.Order {
	if r == nil {
		return nil
	}
	return r.Orders
}

////////////////

func (r *PatchDelegatedOrdersRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PatchDelegatedOrdersRequest) GetOrders() []*order.Order {
	if r == nil {
		return nil
	}
	return r.Orders
}

////////////////

func (r *DeleteDelegatedOrdersRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *DeleteDelegatedOrdersRequest) GetDelegatedOrderIDs() []meta.ID {
	if r == nil {
		return []meta.ID{}
	}
	return r.DelegatedOrderIDs
}

////////////////

func (r *PutSitRepsRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PutSitRepsRequest) GetSitReps() []*order.SitRep {
	if r == nil {
		return nil
	}
	return r.SitReps
}

////////////////

func (r *PatchSitRepsRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PatchSitRepsRequest) GetSitReps() []*order.SitRep {
	if r == nil {
		return nil
	}
	return r.SitReps
}

////////////////

func (r *DeleteSitRepsRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *DeleteSitRepsRequest) GetSitRepIDs() []meta.ID {
	if r == nil {
		return []meta.ID{}
	}
	return r.SitRepIDs
}
