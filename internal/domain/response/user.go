package response

import "github.com/moledoc/orderly/internal/domain/user"

type PostUserResponse struct {
	User *user.User `json:"user"`
}

type GetUserByIDResponse struct {
	User *user.User `json:"user"`
}

type GetUsersResponse struct {
	Users []*user.User `json:"users"`
}

type PatchUserResponse struct {
	User *user.User `json:"user"`
}

type DeleteUserResponse struct{}
