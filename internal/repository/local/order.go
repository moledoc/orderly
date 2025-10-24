package local

import (
	"context"
	"net/http"
	"sync"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/internal/repository"
)

type LocalRepositoryOrder struct {
	mu             sync.Mutex
	OrderHierarchy *order.Order
}

var (
	_ repository.RepositoryOrderAPI = (*LocalRepositoryOrder)(nil)
)

func NewLocalRepositoryOrder() *LocalRepositoryOrder {
	return &LocalRepositoryOrder{
		mu:             sync.Mutex{},
		OrderHierarchy: nil,
	}
}

// func (r *LocalRepositoryOrder) composeOrder(storedOrder *orderInfo) *order.Order {

// 	var delegatedOrders []*order.Order
// 	for _, delegatedID := range storedOrder.DelegatedOrderIDs {
// 		d, ok := r.Orders[delegatedID]
// 		if !ok {
// 			// TODO: log warning
// 			continue
// 		}
// 		delegatedOrders = append(delegatedOrders, d)
// 	}

// 	var sitreps []*order.SitRep
// 	for _, sitrepID := range storedOrder.SitRepIDs {
// 		sr, ok := r.SitReps[sitrepID]
// 		if !ok {
// 			// TODO: log warning
// 			continue
// 		}
// 		sitreps = append(sitreps, sr)
// 	}

// 	if len(delegatedOrders) == 0 {
// 		delegatedOrders = nil
// 	}
// 	if len(sitreps) == 0 {
// 		sitreps = nil
// 	}
// 	resp := &order.Order{
// 		ParentOrderID:   storedOrder.ParentOrderID,
// 		DelegatedOrders: delegatedOrders,
// 		SitReps:         sitreps,
// 		Meta:            storedOrder.Meta,
// 	}
// 	return resp
// }

// func (r *LocalRepositoryOrder) storeOrder(o *order.Order) *orderInfo {

// 	var delegatedOrderIDs []meta.ID
// 	var sitrepIDs []meta.ID
// 	for _, delegated := range o.GetDelegatedOrders() {
// 		delegatedOrderIDs = append(delegatedOrderIDs, delegated.GetID())
// 		r.Orders[delegated.GetID()] = &orderInfo{
// 			ID:            delegated.GetID(),
// 			ParentOrderID: o.GetID(),
// 			Meta:          o.GetMeta(), // NOTE: set meta as order.meta, since they are created at the same time
// 		}
// 	}
// 	for _, sitrep := range o.GetSitReps() {
// 		sitrepIDs = append(sitrepIDs, sitrep.GetID())
// 		r.SitReps[sitrep.GetID()] = sitrep
// 	}
// 	info := &orderInfo{
// 		ID:                o.GetID(),
// 		ParentOrderID:     o.GetParentOrderID(),
// 		DelegatedOrderIDs: delegatedOrderIDs,
// 		SitRepIDs:         sitrepIDs,
// 		Meta:              o.GetMeta(),
// 	}
// 	r.Orders[o.GetID()] = info
// 	parent := r.Orders[o.ParentOrderID]
// 	parent.DelegatedOrderIDs = append(r.Orders[o.ParentOrderID].DelegatedOrderIDs, o.GetID())
// 	return info
// }

// func (r *LocalRepositoryOrder) deleteOrder(storedOrder *orderInfo) {
// 	for _, id := range storedOrder.DelegatedOrderIDs {
// 		delete(r.Orders, id)
// 	}
// 	for _, id := range storedOrder.SitRepIDs {
// 		delete(r.SitReps, id)
// 	}
// 	delete(r.Orders, storedOrder.ID)
// 	return
// }

func (r *LocalRepositoryOrder) Close(ctx context.Context) errwrap.Error {
	middleware.SpanStart(ctx, "LocalStorageOrder:Close")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Close")

	if r == nil {
		return errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r = nil
	return nil
}

func (r *LocalRepositoryOrder) ReadByID(ctx context.Context, id meta.ID) (*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:ReadByID")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:ReadByID")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.OrderHierarchy.GetID() == id {
		return r.OrderHierarchy, nil
	}

	var curDepth []*order.Order
	var nextDepth []*order.Order

	curDepth = append(curDepth, r.OrderHierarchy.GetDelegatedOrders()...)

	for len(curDepth) != 0 {
		for _, do := range curDepth {
			if do.GetID() == id {
				return do, nil
			}
			nextDepth = append(nextDepth, do.GetDelegatedOrders()...)
		}
		curDepth = nextDepth
		nextDepth = []*order.Order{}
	}

	return nil, errwrap.NewError(http.StatusNotFound, "not found")
}

func (r *LocalRepositoryOrder) ReadBy(ctx context.Context, req *request.GetOrdersRequest) ([]*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:ReadBy")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:ReadBy")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	parentOrderID := req.GetParentOrderID()
	accountable := req.GetAccountable()

	var orders []*order.Order
	var curDepth []*order.Order
	var nextDepth []*order.Order

	curDepth = append(curDepth, r.OrderHierarchy.GetDelegatedOrders()...)

	for len(curDepth) != 0 {
		for _, do := range curDepth {
			if (len(parentOrderID) == 0 || parentOrderID == do.GetParentOrderID()) &&
				(len(accountable) == 0 || accountable == do.GetAccountable()) {
				orders = append(orders, do)
			}
			nextDepth = append(nextDepth, do.GetDelegatedOrders()...)
		}
		curDepth = nextDepth
		nextDepth = []*order.Order{}
	}

	return orders, nil
}

func (r *LocalRepositoryOrder) Write(ctx context.Context, order *order.Order) (*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:Write")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Write")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.OrderHierarchy == nil {
		r.OrderHierarchy = order
		return order, nil
	}
	parent, err := r.ReadByID(ctx, order.GetParentOrderID())
	if err != nil {
		return nil, err
	}
	found := false
	for i, do := range parent.GetDelegatedOrders() {
		if do.GetID() == order.GetID() {
			found = true
			parent.GetDelegatedOrders()[i] = order
			break
		}
	}
	if !found {
		parent.SetDelegatedOrders(append(parent.GetDelegatedOrders(), order))
	}

	return order, nil
}

func (r *LocalRepositoryOrder) DeleteOrder(ctx context.Context, id meta.ID) errwrap.Error {
	middleware.SpanStart(ctx, "LocalStorageOrder:DeleteOrder")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:DeleteOrder")

	if r == nil {
		return errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	type opp struct {
		Parent *order.Order
		Order  *order.Order
	}
	var curDepth []*opp
	var nextDepth []*opp

	for _, do := range r.OrderHierarchy.GetDelegatedOrders() {
		curDepth = append(curDepth, &opp{Order: do, Parent: do})
	}

	var op *opp
	for len(curDepth) != 0 && op == nil {
		for _, do := range curDepth {
			if do.Order.GetID() == id {
				op = do
				break
			}
			for _, doi := range do.Order.GetDelegatedOrders() {
				nextDepth = append(nextDepth, &opp{Parent: do.Order, Order: doi})
			}
		}
		curDepth = nextDepth
		nextDepth = []*opp{}
	}

	if op == nil { // NOTE: not found
		return nil
	}

	if len(op.Order.GetDelegatedOrders()) > 0 {
		return errwrap.NewError(http.StatusPreconditionFailed, "order %q has delegated orders", id)
	}

	for i, o := range op.Parent.GetDelegatedOrders() {
		if o.GetID() == op.Order.GetID() {
			pdel := op.Parent.GetDelegatedOrders()
			op.Parent.SetDelegatedOrders(append(pdel[:i], pdel[i+1:]...))
		}
	}

	return nil
}

func (r *LocalRepositoryOrder) DeleteOrders(ctx context.Context, ids []meta.ID) (bool, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:DeleteOrder")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:DeleteOrder")

	if r == nil {
		return false, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	didDelete := false
	// for _, id := range ids {
	// 	_, ok := r.Orders[id]
	// 	didDelete = ok || didDelete
	// 	delete(r.Orders, id)
	// }

	return didDelete, nil
}

func (r *LocalRepositoryOrder) DeleteSitReps(ctx context.Context, ids []meta.ID) (bool, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:DeleteSitRep")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:DeleteSitRep")

	if r == nil {
		return false, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	didDelete := false
	// for _, id := range ids {
	// 	_, ok := r.SitReps[id]
	// 	didDelete = ok || didDelete
	// 	delete(r.SitReps, id)
	// }

	return didDelete, nil
}
