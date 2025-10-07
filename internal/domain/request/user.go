package request

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
)

type PostUserRequest struct {
	User *user.User `json:"user"`
}

type GetUserByIDRequest struct {
	ID meta.ID `json:"id"`
}

type GetUsersRequest struct{}

type GetUserSubOrdinatesRequest struct {
	ID meta.ID `json:"id"`
}

type PatchUserRequest struct {
	User *user.User `json:"user"`
}

type DeleteUserRequest struct {
	ID meta.ID `json:"id"`
}
