package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/stretchr/testify/require"
)

func userObj(extra ...string) *user.User {
	ee := strings.Join(extra, ".")
	return &user.User{
		ID:         meta.NewID(),
		Name:       fmt.Sprintf("name%v", ee),
		Email:      user.Email(fmt.Sprintf("example%v@example.com", ee)),
		Supervisor: user.Email(fmt.Sprintf("example.supervisor%v@example.com", ee)),
		Meta: &meta.Meta{
			Version: 1,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		},
	}
}

func (s *UserSuite) TestValidation_User() {
	tt := s.T()
	tt.Run("user.id", func(t *testing.T) {
		u := userObj()
		u.SetID("")
		err := mgmtuser.ValidateUser(u)
		require.NoError(t, err) // NOTE: only being validated if len(user.id) > 0; it's to enable to use common user validation func across endpoints. user.ID checks are done in request validation
	})

	tt.Run("user.name", func(t *testing.T) {
		u := userObj()
		u.SetName("")
		err := mgmtuser.ValidateUser(u)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid user.name length", err.GetStatusMessage())
	})

	tt.Run("user.email", func(t *testing.T) {
		u := userObj()
		u.SetEmail("")
		err := mgmtuser.ValidateUser(u)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid user.email: invalid email length", err.GetStatusMessage())
	})

	tt.Run("user.supervisor", func(t *testing.T) {
		u := userObj()
		u.SetSupervisor("")
		err := mgmtuser.ValidateUser(u)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid user.supervisor: invalid email length", err.GetStatusMessage())
	})

	tt.Run("user.meta", func(t *testing.T) {
		u := userObj()
		u.SetMeta(&meta.Meta{
			Version: 2,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		})
		err := mgmtuser.ValidateUser(u)
		require.NoError(t, err) // NOTE: meta is not validated, as input.meta is ignored throughout the service
	})
}

func (s *UserSuite) TestValidation_PostUserRequest() {
	tt := s.T()

	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{})
		defer cleanup.User(t, s.API, resp.GetUser())
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty request", err.GetStatusMessage())
	})

	tt.Run("user.id.provided", func(t *testing.T) {
		resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
			User: userObj(),
		})
		defer cleanup.User(t, s.API, resp.GetUser())
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "user.id disallowed", err.GetStatusMessage())
	})
}

func (s *UserSuite) TestValidation_GetUserByIDRequest() {
	tt := s.T()

	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty request", err.GetStatusMessage())
	})
}

func (s *UserSuite) TestValidation_GetUsersRequest() {
	tt := s.T()

	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.GetUsers(t, context.Background(), &request.GetUsersRequest{})
		require.NoError(t, err)
		require.Empty(t, resp)
	})
}

func (s *UserSuite) TestValidation_GetUserSubOrdinatesRequest() {
	tt := s.T()

	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.GetUserSubOrdinates(t, context.Background(), &request.GetUserSubOrdinatesRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty request", err.GetStatusMessage())
	})
}

func (s *UserSuite) TestValidation_PatchUserRequest() {
	tt := s.T()

	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty user", err.GetStatusMessage())
	})
	tt.Run("empty.user.id", func(t *testing.T) {
		resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
			User: &user.User{
				Name: "patch",
			},
		})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "user.id missing", err.GetStatusMessage())
	})
}

func (s *UserSuite) TestValidation_DeleteUserRequest() {
	tt := s.T()

	tt.Run("empty.request", func(t *testing.T) {
		resp, err := s.API.DeleteUser(t, context.Background(), &request.DeleteUserRequest{})
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "empty request", err.GetStatusMessage())
	})
}
