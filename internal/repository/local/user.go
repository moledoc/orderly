package local

import (
	"context"
	"net/http"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/internal/repository"
	"github.com/moledoc/orderly/pkg/utils"
)

type userInfo struct {
	ID         meta.ID
	Name       string
	Email      user.Email
	Supervisor user.Email
	Meta       meta.Meta
}

type LocalRepositoryUser map[meta.ID]userInfo

var (
	_ repository.RepositoryUserAPI = (LocalRepositoryUser)(nil)
)

func NewLocalRepositoryUser() LocalRepositoryUser {
	return make(LocalRepositoryUser)
}

func (r LocalRepositoryUser) Close(ctx context.Context) errwrap.Error {
	middleware.SpanStart(ctx, "LocalRepositoryUser:Close")
	defer middleware.SpanStop(ctx, "LocalRepositoryUser:Close")

	if r == nil {
		return errwrap.NewError(http.StatusInternalServerError, "local repository user uninitialized")
	}
	r = nil
	return nil
}

func (r LocalRepositoryUser) ReadByID(ctx context.Context, ID meta.ID) (*user.User, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalRepositoryUser:ReadByID")
	defer middleware.SpanStop(ctx, "LocalRepositoryUser:ReadByID")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository user uninitialized")
	}

	u, ok := r[ID]
	if !ok {
		return nil, errwrap.NewError(http.StatusNotFound, "not found")
	}
	return &user.User{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Supervisor: u.Supervisor,
		Meta:       &u.Meta,
	}, nil
}

func (r LocalRepositoryUser) ReadSubOrdinates(ctx context.Context, ID meta.ID) ([]*user.User, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalRepositoryUser:ReadSubOrdinates")
	defer middleware.SpanStop(ctx, "LocalRepositoryUser:ReadSubOrdinates")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository user uninitialized")
	}

	supervisor, ok := r[ID]
	if !ok {
		return nil, errwrap.NewError(http.StatusNotFound, "not found")
	}

	// MAYBE: TODO: optimize subordinate finding
	var subOrdinates []*user.User
	for _, u := range r {
		if u.Supervisor == supervisor.Email {
			subOrdinates = append(subOrdinates, &user.User{
				ID:         u.ID,
				Name:       u.Name,
				Email:      u.Email,
				Supervisor: u.Supervisor,
				Meta:       &u.Meta,
			})
		}
	}

	return subOrdinates, nil
}

func (r LocalRepositoryUser) ReadAll(ctx context.Context) ([]*user.User, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalRepositoryUser:ReadAll")
	defer middleware.SpanStop(ctx, "LocalRepositoryUser:ReadAll")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository user uninitialized")
	}

	var users []*user.User
	for _, u := range r {
		users = append(users, &user.User{
			ID:         u.ID,
			Name:       u.Name,
			Email:      u.Email,
			Supervisor: u.Supervisor,
			Meta:       &u.Meta,
		})
	}
	return users, nil
}

func (r LocalRepositoryUser) Write(ctx context.Context, user *user.User) (*user.User, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalRepositoryUser:Write")
	defer middleware.SpanStop(ctx, "LocalRepositoryUser:Write")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository user uninitialized")
	}

	id := user.GetID()
	r[id] = userInfo{
		ID:         user.GetID(),
		Name:       user.GetName(),
		Email:      user.GetEmail(),
		Supervisor: user.GetSupervisor(),
		Meta:       utils.Deref(user.GetMeta()),
	}

	return user, nil
}

func (r LocalRepositoryUser) Delete(ctx context.Context, ID meta.ID) errwrap.Error {
	middleware.SpanStart(ctx, "LocalStorageUser:Delete")
	defer middleware.SpanStop(ctx, "LocalStorageUser:Delete")

	if r == nil {
		return errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}

	delete(r, ID)

	return nil
}
