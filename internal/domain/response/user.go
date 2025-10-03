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

type GetUserVersionsResponse struct {
	UserVersions []*user.User `json:"user_versions"`
}

type GetUserSubOrdinatesResponse struct {
	SubOrdinates []*user.User `json:"sub_ordinates"`
}

type PatchUserResponse struct {
	User *user.User `json:"user"`
}

type DeleteUserResponse struct{}
