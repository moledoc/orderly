package request

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
)

type PostUserRequest struct {
	User *user.User `json:"user,omitempty"`
}

type GetUserByIDRequest struct {
	ID meta.ID `json:"id,omitempty"`
}

type GetUserByRequest struct {
	ID         meta.ID    `json:"id,omitempty"`
	Email      user.Email `json:"email,omitempty"`
	Supervisor user.Email `json:"supervisor,omitempty"`
}

type GetUsersRequest struct{}

type GetUserSubOrdinatesRequest struct {
	ID meta.ID `json:"id,omitempty"`
}

type PatchUserRequest struct {
	User *user.User `json:"user,omitempty"`
}

type DeleteUserRequest struct {
	ID meta.ID `json:"id,omitempty"`
}
