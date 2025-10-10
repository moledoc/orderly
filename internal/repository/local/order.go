package local

import (
	"context"
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/internal/repository"
	"github.com/moledoc/orderly/pkg/utils"
)

type orderInfo struct {
	TaskID           meta.ID
	ParentOrderID    meta.ID
	DelegatedTaskIDs []meta.ID
	SitRepIDs        []meta.ID
	Meta             meta.Meta
}

type LocalRepositoryOrder struct {
	Orders  map[meta.ID]orderInfo
	Tasks   map[meta.ID]order.Task
	SitReps map[meta.ID]order.SitRep
}

var (
	_ repository.RepositoryOrderAPI = (*LocalRepositoryOrder)(nil)
)

func NewLocalRepositoryOrder() *LocalRepositoryOrder {
	return &LocalRepositoryOrder{
		Orders:  make(map[meta.ID]orderInfo),
		Tasks:   make(map[meta.ID]order.Task),
		SitReps: make(map[meta.ID]order.SitRep),
	}
}

func (r *LocalRepositoryOrder) composeOrder(storedOrder orderInfo) *order.Order {
	task := r.Tasks[storedOrder.TaskID]

	var delegatedTasks []*order.Task
	for _, delegatedID := range storedOrder.DelegatedTaskIDs {
		d, ok := r.Tasks[delegatedID]
		if !ok {
			// TODO: log warning
			continue
		}
		delegatedTasks = append(delegatedTasks, &d)
	}

	var sitreps []*order.SitRep
	for _, sitrepID := range storedOrder.SitRepIDs {
		sr, ok := r.SitReps[sitrepID]
		if !ok {
			// TODO: log warning
			continue
		}
		sitreps = append(sitreps, &sr)
	}

	if len(delegatedTasks) == 0 {
		delegatedTasks = nil
	}
	if len(sitreps) == 0 {
		sitreps = nil
	}
	resp := &order.Order{
		Task:           &task,
		ParentOrderID:  storedOrder.ParentOrderID,
		DelegatedTasks: delegatedTasks,
		SitReps:        sitreps,
		Meta:           &storedOrder.Meta,
	}
	return resp
}

func (r *LocalRepositoryOrder) storeOrder(o *order.Order) orderInfo {

	r.Tasks[o.GetTask().GetID()] = utils.Deref(o.GetTask())

	var delegatedTaskIDs []meta.ID
	var sitrepIDs []meta.ID
	for _, delegated := range o.GetDelegatedTasks() {
		delegatedTaskIDs = append(delegatedTaskIDs, delegated.GetID())
		r.Tasks[delegated.GetID()] = utils.Deref(delegated)
	}
	for _, sitrep := range o.GetSitReps() {
		sitrepIDs = append(sitrepIDs, sitrep.GetID())
		r.SitReps[sitrep.GetID()] = utils.Deref(sitrep)
	}
	info := orderInfo{
		TaskID:           o.GetID(),
		ParentOrderID:    o.GetParentOrderID(),
		DelegatedTaskIDs: delegatedTaskIDs,
		SitRepIDs:        sitrepIDs,
		Meta:             utils.Deref(o.GetMeta()),
	}
	r.Orders[o.GetID()] = info
	return info
}

func (r *LocalRepositoryOrder) deleteOrder(storedOrder orderInfo) {
	delete(r.Tasks, storedOrder.TaskID)
	for _, id := range storedOrder.DelegatedTaskIDs {
		delete(r.Tasks, id)
	}
	for _, id := range storedOrder.SitRepIDs {
		delete(r.SitReps, id)
	}
	delete(r.Orders, storedOrder.TaskID)
	return
}

func (r *LocalRepositoryOrder) Close(ctx context.Context) errwrap.Error {
	middleware.SpanStart(ctx, "LocalStorageOrder:Close")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Close")

	if r == nil {
		return errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	r = nil
	return nil
}

func (r *LocalRepositoryOrder) ReadByID(ctx context.Context, ID meta.ID) (*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:ReadByID")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:ReadByID")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	storedOrder, ok := r.Orders[ID]
	if !ok {
		return nil, errwrap.NewError(http.StatusNotFound, "not found")
	}

	resp := r.composeOrder(storedOrder)

	return resp, nil
}

func (r *LocalRepositoryOrder) ReadSubOrders(ctx context.Context, ID meta.ID) ([]*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:ReadSubOrders")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:ReadSubOrders")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	parentOrder, ok := r.Orders[ID]
	if !ok {
		return nil, errwrap.NewError(http.StatusNotFound, "not found")
	}

	var storedOrders []orderInfo
	for _, o := range r.Orders {
		if o.ParentOrderID == parentOrder.TaskID {
			storedOrders = append(storedOrders, o)
		}
	}

	var subOrders []*order.Order
	for _, storedOrder := range storedOrders {
		subOrders = append(subOrders, r.composeOrder(storedOrder))
	}

	return subOrders, nil
}

func (r *LocalRepositoryOrder) ReadAll(ctx context.Context) ([]*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:ReadAll")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:ReadAll")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	var orders []*order.Order
	for _, storedOrder := range r.Orders {
		orders = append(orders, r.composeOrder(storedOrder))
	}

	return orders, nil
}

func (r *LocalRepositoryOrder) Write(ctx context.Context, order *order.Order) (*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:Write")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Write")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	info := r.storeOrder(order)
	order = r.composeOrder(info)

	return order, nil
}

func (r *LocalRepositoryOrder) DeleteOrder(ctx context.Context, ID meta.ID) errwrap.Error {
	middleware.SpanStart(ctx, "LocalStorageOrder:DeleteOrder")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:DeleteOrder")

	if r == nil {
		return errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	storedOrder, ok := r.Orders[ID]
	if !ok {
		return nil
	}
	r.deleteOrder(storedOrder)

	return nil
}

func (r *LocalRepositoryOrder) DeleteTasks(ctx context.Context, IDs []meta.ID) (bool, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:DeleteTask")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:DeleteTask")

	if r == nil {
		return false, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	didDelete := false
	for _, id := range IDs {
		_, ok := r.Tasks[id]
		didDelete = ok || didDelete
		delete(r.Tasks, id)
	}

	return didDelete, nil
}

func (r *LocalRepositoryOrder) DeleteSitReps(ctx context.Context, IDs []meta.ID) (bool, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:DeleteSitRep")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:DeleteSitRep")

	if r == nil {
		return false, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	didDelete := false
	for _, id := range IDs {
		_, ok := r.SitReps[id]
		didDelete = ok || didDelete
		delete(r.SitReps, id)
	}

	return didDelete, nil
}
