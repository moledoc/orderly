package setup

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/api"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/stretchr/testify/require"
)

func UserObj(extra ...string) *user.User {
	ee := strings.Join(append([]string{""}, extra...), ".")
	return &user.User{
		Name:         fmt.Sprintf("name%v", ee),
		Email:        user.Email(fmt.Sprintf("example%v@example.com", ee)),
		SupervisorID: user.Email(fmt.Sprintf("example.supervisor%v@example.com", ee)),
		Meta: &meta.Meta{
			Version: 1,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		},
	}
}

func UserObjWithID(extra ...string) *user.User {
	u := UserObj(extra...)
	u.ID = meta.NewID()
	return u
}

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
