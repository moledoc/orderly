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

func (r *PutDelegatedTaskRequest) GetTasks() []*order.Task {
	if r == nil {
		return nil
	}
	return r.Tasks
}

////////////////

func (r *PatchDelegatedTaskRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PatchDelegatedTaskRequest) GetTasks() []*order.Task {
	if r == nil {
		return nil
	}
	return r.Tasks
}

////////////////

func (r *DeleteDelegatedTaskRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *DeleteDelegatedTaskRequest) GetDelegatedTaskIDs() []meta.ID {
	if r == nil {
		return []meta.ID{}
	}
	return r.DelegatedTaskIDs
}

////////////////

func (r *PutSitRepRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PutSitRepRequest) GetSitReps() []*order.SitRep {
	if r == nil {
		return nil
	}
	return r.SitReps
}

////////////////

func (r *PatchSitRepRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PatchSitRepRequest) GetSitReps() []*order.SitRep {
	if r == nil {
		return nil
	}
	return r.SitReps
}

////////////////

func (r *DeleteSitRepRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *DeleteSitRepRequest) GetSitRepIDs() []meta.ID {
	if r == nil {
		return []meta.ID{}
	}
	return r.SitRepIDs
}
