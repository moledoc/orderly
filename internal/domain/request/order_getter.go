package request

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/user"
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

func (r *GetOrdersRequest) GetAccountable() user.Email {
	if r == nil {
		return ""
	}
	return r.Accountable
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

func (r *PutDelegatedTasksRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PutDelegatedTasksRequest) GetTasks() []*order.Task {
	if r == nil {
		return nil
	}
	return r.Tasks
}

////////////////

func (r *PatchDelegatedTasksRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *PatchDelegatedTasksRequest) GetTasks() []*order.Task {
	if r == nil {
		return nil
	}
	return r.Tasks
}

////////////////

func (r *DeleteDelegatedTasksRequest) GetOrderID() meta.ID {
	if r == nil {
		return ""
	}
	return r.OrderID
}

func (r *DeleteDelegatedTasksRequest) GetDelegatedTaskIDs() []meta.ID {
	if r == nil {
		return []meta.ID{}
	}
	return r.DelegatedTaskIDs
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
