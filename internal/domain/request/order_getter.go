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

func (r *GetOrderSubOrdersRequest) GetID() meta.ID {
	if r == nil {
		return ""
	}
	return r.ID
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

func (r *PutDelegatedTaskRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PutDelegatedTaskRequest) GetTask() *order.Task {
	if r == nil {
		return nil
	}
	return r.Task
}

////////////////

func (r *PatchDelegatedTaskRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PatchDelegatedTaskRequest) GetTask() *order.Task {
	if r == nil {
		return nil
	}
	return r.Task
}

////////////////

func (r *DeleteDelegatedTaskRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *DeleteDelegatedTaskRequest) GetDelegatedTaskID() meta.ID {
	if r == nil {
		return ""
	}
	return r.DelegatedTaskID
}

////////////////

func (r *PutSitRepRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PutSitRepRequest) GetSitRep() *order.SitRep {
	if r == nil {
		return nil
	}
	return r.SitRep
}

////////////////

func (r *PatchSitRepRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PatchSitRepRequest) GetSitRep() *order.SitRep {
	if r == nil {
		return nil
	}
	return r.SitRep
}

////////////////

func (r *DeleteSitRepRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *DeleteSitRepRequest) GetSitRepID() meta.ID {
	if r == nil {
		return ""
	}
	return r.SitRepID
}
