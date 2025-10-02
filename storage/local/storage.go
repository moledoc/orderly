package local

import (
	"context"
	"net/http"
	"time"

	"github.com/moledoc/orderly/actions"
	"github.com/moledoc/orderly/middleware"
	"github.com/moledoc/orderly/models"
	"github.com/moledoc/orderly/utils"
)

type Storage map[uint][]*models.User

func New() Storage {
	return make(Storage)
}

func (s Storage) Close(ctx context.Context) {

	middleware.SpanStart(ctx, "LocalStorage:Close")
	defer middleware.SpanStop(ctx, "LocalStorage:Close")

	s = nil
}

func (s Storage) Read(ctx context.Context, action actions.Action, id uint) ([]*models.User, models.IError) {
	middleware.SpanStart(ctx, "LocalStorage:Read")
	defer middleware.SpanStop(ctx, "LocalStorage:Read")

	if s == nil {
		return nil, models.NewError(http.StatusInternalServerError, "localstorage not initialized for read")
	}

	switch action {
	case actions.READ:
		middleware.SpanStart(ctx, "LocalStorage:Read:READ")
		defer middleware.SpanStop(ctx, "LocalStorage:Read:READ")
		us, ok := s[id]
		if !ok || len(us) == 0 {
			return nil, models.NewError(http.StatusNotFound, "not found during read")
		}
		return []*models.User{us[len(us)-1]}, nil
	case actions.READVERSIONS:
		middleware.SpanStart(ctx, "LocalStorage:Read:READVERSIONS")
		defer middleware.SpanStop(ctx, "LocalStorage:Read:READVERSIONS")
		us, ok := s[id]
		if !ok || len(us) == 0 {
			return nil, models.NewError(http.StatusNotFound, "not found during read")
		}
		return us, nil
	case actions.READALL:
		middleware.SpanStart(ctx, "LocalStorage:Read:READALL")
		defer middleware.SpanStop(ctx, "LocalStorage:Read:READALL")
		uss := make([]*models.User, len(s))
		i := 0
		for _, us := range s {
			if len(us) == 0 {
				continue
			}
			uss[i] = us[len(us)-1]
			i += 1
		}
		return uss, nil
	case actions.READSUBORDINATES:
		middleware.SpanStart(ctx, "LocalStorage:Read:READSUBORDINATES")
		defer middleware.SpanStop(ctx, "LocalStorage:Read:READSUBORDINATES")
		ssupervisor, ok := s[id]
		if !ok || len(ssupervisor) == 0 {
			return nil, models.NewError(http.StatusNotFound, "not found during read")
		}
		supervisor := ssupervisor[len(ssupervisor)-1]

		var uss []*models.User
		for _, us := range s {
			if len(us) == 0 {
				continue
			}
			if u := us[len(us)-1]; u.Supervisor != nil && *u.Supervisor == *supervisor.Email {
				uss = append(uss, u)
			}
		}
		return uss, nil
	default:
		return nil, models.NewError(http.StatusInternalServerError, "undefined read action")
	}
}

func (s Storage) Write(ctx context.Context, action actions.Action, user *models.User) (*models.User, models.IError) {

	middleware.SpanStart(ctx, "LocalStorage:Write")
	defer middleware.SpanStop(ctx, "LocalStorage:Write")

	if s == nil {
		return nil, models.NewError(http.StatusInternalServerError, "localstorage not initialized for write")
	}
	if user == nil {
		return nil, models.NewError(http.StatusInternalServerError, "invalid user object in write")
	}

	var us []*models.User
	var ok bool
	if user.ID != nil {
		us, ok = s[*user.ID]
	}

	switch action {

	case actions.CREATE:
		middleware.SpanStart(ctx, "LocalStorage:Write:CREATE")
		defer middleware.SpanStop(ctx, "LocalStorage:Write:CREATE")
		if ok || len(us) > 0 {
			return nil, models.NewError(http.StatusConflict, "already exists during write")
		}
		id := uint(len(s) + 1)
		user.ID = &id
		now := time.Now().UTC()
		user.Meta = &models.Meta{
			Version: uint(1),
			Created: now,
			Updated: now,
		}
		s[id] = append(s[id], user)
		return user, nil

	case actions.UPDATE:
		middleware.SpanStart(ctx, "LocalStorage:Write:UPDATE")
		defer middleware.SpanStop(ctx, "LocalStorage:Write:UPDATE")
		if !ok || len(us) == 0 {
			return nil, models.NewError(http.StatusNotFound, "not found during write")
		}
		var updUser models.User = *(us[len(us)-1])
		updated := false
		if user.Name != nil && utils.Deref(updUser.Name) != utils.Deref(user.Name) {
			updUser.Name = user.Name
			updated = true
		}
		if user.Email != nil && utils.Deref(updUser.Email) != utils.Deref(user.Email) {
			updUser.Email = user.Email
			updated = true
		}
		if user.Supervisor != nil && utils.Deref(updUser.Supervisor) != utils.Deref(user.Supervisor) {
			updUser.Supervisor = user.Supervisor
			updated = true
		}
		if updated {
			now := time.Now().UTC()
			updUser.Meta = &models.Meta{
				Version: updUser.Meta.Version + 1,
				Created: updUser.Meta.Created,
				Updated: now,
			}
			s[*user.ID] = append(s[*user.ID], &updUser)
		}
		us = s[*user.ID]
		return us[len(us)-1], nil

	case actions.DELETESOFT:
		middleware.SpanStart(ctx, "LocalStorage:Write:SOFTDELETE")
		defer middleware.SpanStop(ctx, "LocalStorage:Write:SOFTDELETE")
		if ok {
			for _, u := range us {
				u.Meta.Deleted = true
			}
		}
		return nil, nil

	case actions.DELETEHARD:
		middleware.SpanStart(ctx, "LocalStorage:Write:HARDDELETE")
		defer middleware.SpanStop(ctx, "LocalStorage:Write:HARDDELETE")
		if ok {
			delete(s, *user.ID)
		}
		return nil, nil

	default:
		return nil, models.NewError(http.StatusInternalServerError, "undefined write action")
	}
}
