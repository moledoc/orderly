package storage

import (
	"context"

	"github.com/moledoc/orderly/actions"
	"github.com/moledoc/orderly/models"
)

type StorageUserAPI interface {
	Close(ctx context.Context)
	Read(ctx context.Context, action actions.Action, id uint) ([]*models.User, models.IError)
	Write(ctx context.Context, action actions.Action, user *models.User) (*models.User, models.IError)
}
