package response

import (
	"github.com/moledoc/orderly/internal/domain/user"
)

func (r *PostUserResponse) GetUser() *user.User {
	if r == nil {
		return nil
	}
	return r.User
}

////////////////

func (r *GetUserByIDResponse) GetUser() *user.User {
	if r == nil {
		return nil
	}
	return r.User
}

////////////////

func (r *GetUserSubOrdinatesResponse) GetSubOrdinates() []*user.User {
	if r == nil {
		return nil
	}
	return r.SubOrdinates
}

////////////////

func (r *PatchUserResponse) GetUser() *user.User {
	if r == nil {
		return nil
	}
	return r.User
}
