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
	if r == nil {
		return ""
	}
	return r.ID
}

////////////////

func (r *GetUsersRequest) GetEmails() []user.Email {
	if r == nil {
		return []user.Email{}
	}
	return r.Emails
}

func (r *GetUsersRequest) GetSupervisor() user.Email {
	if r == nil {
		return ""
	}
	return r.Supervisor
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
	if r == nil {
		return ""
	}
	return r.ID
}
