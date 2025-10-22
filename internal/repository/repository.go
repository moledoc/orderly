package repository

import (
	"context"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/user"
)

type RepositoryOrderAPI interface {
	Close(ctx context.Context) errwrap.Error
	ReadByID(ctx context.Context, id meta.ID) (*order.Order, errwrap.Error)
	ReadSubOrders(ctx context.Context, id meta.ID) ([]*order.Order, errwrap.Error)
	ReadAll(ctx context.Context) ([]*order.Order, errwrap.Error)
	// TODO: split Write to specific funcs
	Write(ctx context.Context, order *order.Order) (*order.Order, errwrap.Error)
	DeleteOrder(ctx context.Context, id meta.ID) errwrap.Error
	DeleteTasks(ctx context.Context, IDs []meta.ID) (bool, errwrap.Error)
	DeleteSitReps(ctx context.Context, IDs []meta.ID) (bool, errwrap.Error)
	ReadUserOrders(ctx context.Context, userID meta.ID) ([]*order.Order, errwrap.Error)
}

type RepositoryUserAPI interface {
	Close(ctx context.Context) errwrap.Error
	ReadByID(ctx context.Context, id meta.ID) (*user.User, errwrap.Error)
	ReadSubOrdinates(ctx context.Context, id meta.ID) ([]*user.User, errwrap.Error)
	ReadAll(ctx context.Context) ([]*user.User, errwrap.Error)
	// TODO: split Write to specific funcs
	Write(ctx context.Context, user *user.User) (*user.User, errwrap.Error)
	Delete(ctx context.Context, id meta.ID) errwrap.Error
}
