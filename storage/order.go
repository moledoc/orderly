package storage

import (
	"context"

	"github.com/moledoc/orderly/actions"
	"github.com/moledoc/orderly/models"
)

type StorageOrderAPI interface {
	Close(ctx context.Context)
	Read(ctx context.Context, action actions.Action, id uint) ([]*models.Order, models.IError)
	Write(ctx context.Context, action actions.Action, order *models.Order) (*models.Order, models.IError)
}
