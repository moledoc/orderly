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

type GetUsersRequest struct {
	Emails       []user.Email `json:"emails,omitempty"`
	SupervisorID meta.ID      `json:"supervisor_id,omitempty"`
}

type PatchUserRequest struct {
	User *user.User `json:"user,omitempty"`
}

type DeleteUserRequest struct {
	ID meta.ID `json:"id,omitempty"`
}
