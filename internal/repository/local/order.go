package local

import (
	"context"
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/internal/repository"
)

type LocalRepositoryOrder map[meta.ID][]*order.Order

var (
	_ repository.RepositoryOrderAPI = (LocalRepositoryOrder)(nil)
)

func NewLocalRepositoryOrder() LocalRepositoryOrder {
	return make(LocalRepositoryOrder)
}

func (r LocalRepositoryOrder) Close(ctx context.Context) errwrap.Error {
	middleware.SpanStart(ctx, "LocalStorageOrder:Close")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Close")

	if r == nil {
		return errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	r = nil
	return nil
}

func (r LocalRepositoryOrder) ReadByID(ctx context.Context, ID meta.ID) (*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:Read")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Read")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	order, ok := r[ID]
	if !ok || len(order) == 0 {
		return nil, errwrap.NewError(http.StatusNotFound, "not found")
	}

	return order[len(order)-1], nil
}

func (r LocalRepositoryOrder) ReadSubOrders(ctx context.Context, ID meta.ID) ([]*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:ReadSubOrders")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:ReadSubOrders")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	parentOrder, ok := r[ID]
	if !ok || len(parentOrder) == 0 {
		return nil, errwrap.NewError(http.StatusNotFound, "not found")
	}
	po := parentOrder[len(parentOrder)-1]

	// MAYBE: TODO: optimize sub-order finding
	var subOrders []*order.Order
	for _, order := range r {
		if len(order) == 0 {
			continue
		}
		o := order[len(order)-1]
		if o.GetParentOrderID() == po.GetTask().GetID() {
			subOrders = append(subOrders, o)
		}

	}

	return subOrders, nil
}

func (r LocalRepositoryOrder) ReadVersions(ctx context.Context, ID meta.ID) ([]*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:ReadVersions")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:ReadVersions")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	order, ok := r[ID]
	if !ok || len(order) == 0 {
		return nil, errwrap.NewError(http.StatusNotFound, "not found")
	}

	return order, nil
}

func (r LocalRepositoryOrder) ReadAll(ctx context.Context) ([]*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:ReadAll")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:ReadAll")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	var orders []*order.Order
	for _, order := range r {
		if len(order) == 0 {
			continue
		}
		orders = append(orders, order[len(order)-1])
	}

	return orders, nil
}

func (r LocalRepositoryOrder) Write(ctx context.Context, order *order.Order) (*order.Order, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalStorageOrder:Write")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Write")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	id := order.GetTask().GetID()
	r[id] = append(r[id], order)

	return order, nil
}

func (r LocalRepositoryOrder) Delete(ctx context.Context, ID meta.ID) errwrap.Error {
	middleware.SpanStart(ctx, "LocalStorageOrder:Delete")
	defer middleware.SpanStop(ctx, "LocalStorageOrder:Delete")

	if r == nil {
		return errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	delete(r, ID)

	return nil
}
