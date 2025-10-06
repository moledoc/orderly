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
	ReadByID(ctx context.Context, ID meta.ID) (*order.Order, errwrap.Error)
	ReadSubOrders(ctx context.Context, ID meta.ID) ([]*order.Order, errwrap.Error)
	ReadAll(ctx context.Context) ([]*order.Order, errwrap.Error)
	Write(ctx context.Context, order *order.Order) (*order.Order, errwrap.Error)
	Delete(ctx context.Context, ID meta.ID) errwrap.Error
}

type RepositoryUserAPI interface {
	Close(ctx context.Context) errwrap.Error
	ReadByID(ctx context.Context, ID meta.ID) (*user.User, errwrap.Error)
	ReadSubOrdinates(ctx context.Context, ID meta.ID) ([]*user.User, errwrap.Error)
	ReadAll(ctx context.Context) ([]*user.User, errwrap.Error)
	Write(ctx context.Context, user *user.User) (*user.User, errwrap.Error)
	Delete(ctx context.Context, ID meta.ID) errwrap.Error
}
