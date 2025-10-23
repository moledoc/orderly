package local

import (
	"context"
	"net/http"
	"slices"
	"sync"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/internal/repository"
)

type userInfo struct {
	ID         meta.ID
	Name       string
	Email      user.Email
	Supervisor user.Email
	Meta       meta.Meta
}

type LocalRepositoryUser struct {
	mu sync.Mutex
	db map[meta.ID]*user.User
}

var (
	_ repository.RepositoryUserAPI = (*LocalRepositoryUser)(nil)
)

func NewLocalRepositoryUser() *LocalRepositoryUser {
	return &LocalRepositoryUser{
		mu: sync.Mutex{},
		db: make(map[meta.ID]*user.User),
	}
}

func (r *LocalRepositoryUser) Close(ctx context.Context) errwrap.Error {
	middleware.SpanStart(ctx, "LocalRepositoryUser:Close")
	defer middleware.SpanStop(ctx, "LocalRepositoryUser:Close")

	if r == nil {
		return errwrap.NewError(http.StatusInternalServerError, "local repository user uninitialized")
	}
	r = nil
	return nil
}

func (r *LocalRepositoryUser) ReadByID(ctx context.Context, id meta.ID) (*user.User, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalRepositoryUser:ReadByID")
	defer middleware.SpanStop(ctx, "LocalRepositoryUser:ReadByID")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository user uninitialized")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.db[id]
	if !ok {
		return nil, errwrap.NewError(http.StatusNotFound, "not found")
	}
	return u, nil
}

func (r *LocalRepositoryUser) ReadBy(ctx context.Context, req *request.GetUsersRequest) ([]*user.User, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalRepositoryUser:ReadBy")
	defer middleware.SpanStop(ctx, "LocalRepositoryUser:ReadBy")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository user uninitialized")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	var users []*user.User
	emails := req.GetEmails()
	supervisor := req.GetSupervisor()
	for _, u := range r.db {
		if (len(emails) == 0 || slices.Contains(emails, u.GetEmail())) &&
			(len(supervisor) == 0 || supervisor == u.GetSupervisor()) {
			users = append(users, u)
		}
	}

	return users, nil
}

func (r *LocalRepositoryUser) Write(ctx context.Context, user *user.User) (*user.User, errwrap.Error) {
	middleware.SpanStart(ctx, "LocalRepositoryUser:Write")
	defer middleware.SpanStop(ctx, "LocalRepositoryUser:Write")

	if r == nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "local repository user uninitialized")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	id := user.GetID()
	r.db[id] = user

	return user, nil
}

func (r *LocalRepositoryUser) Delete(ctx context.Context, id meta.ID) errwrap.Error {
	middleware.SpanStart(ctx, "LocalStorageUser:Delete")
	defer middleware.SpanStop(ctx, "LocalStorageUser:Delete")

	if r == nil {
		return errwrap.NewError(http.StatusInternalServerError, "local repository uninitialized")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.db, id)

	return nil
}
