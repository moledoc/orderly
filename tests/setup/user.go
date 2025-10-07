package setup

import (
	"context"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/api"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/stretchr/testify/require"
)

func CreateUserWithCleanup(t *testing.T, ctx context.Context, api api.User, userObj *user.User) (*user.User, errwrap.Error) {
	resp, err := api.PostUser(t, ctx, &request.PostUserRequest{
		User: userObj,
	})
	if err != nil {
		return nil, err
	}
	cleanup.User(t, api, resp.GetUser())
	return resp.GetUser(), nil
}

func MustCreateUserWithCleanup(t *testing.T, ctx context.Context, api api.User, userObj *user.User) *user.User {
	user, err := CreateUserWithCleanup(t, ctx, api, userObj)
	require.NoError(t, err)
	return user
}
