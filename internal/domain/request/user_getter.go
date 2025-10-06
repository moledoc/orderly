package request

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
)

func (r *PostUserRequest) GetUser() *user.User {
	if r == nil {
		return nil
	}
	return r.User
}

////////////////

func (r *GetUserByIDRequest) GetID() meta.ID {
	if r == nil || r.ID == nil {
		return ""
	}
	return *r.ID
}

////////////////

func (r *GetUserSubOrdinatesRequest) GetID() meta.ID {
	if r == nil || r.ID == nil {
		return ""
	}
	return *r.ID
}

////////////////

func (r *PatchUserRequest) GetUser() *user.User {
	if r == nil {
		return nil
	}
	return r.User
}

////////////////

func (r *DeleteUserRequest) GetID() meta.ID {
	if r == nil || r.ID == nil {
		return ""
	}
	return *r.ID
}
