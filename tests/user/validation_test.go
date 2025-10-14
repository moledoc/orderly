package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/service/common/validation"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestValidation_User() {
	tt := s.T()
	tt.Run("user.id", func(t *testing.T) {
		u := setup.UserObj()
		u.SetID("")
		err := mgmtuser.ValidateUser(u, validation.IgnoreNothing)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid user.id: invalid id length", err.GetStatusMessage())
	})

	tt.Run("user.name", func(t *testing.T) {
		u := setup.UserObjWithID()
		u.SetName("")
		err := mgmtuser.ValidateUser(u, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid user.name length", err.GetStatusMessage())
	})

	tt.Run("user.email", func(t *testing.T) {
		u := setup.UserObjWithID()
		u.SetEmail("")
		err := mgmtuser.ValidateUser(u, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid user.email: invalid email length", err.GetStatusMessage())
	})

	tt.Run("user.supervisor", func(t *testing.T) {
		u := setup.UserObjWithID()
		u.SetSupervisor("")
		err := mgmtuser.ValidateUser(u, validation.IgnoreNothing)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		require.Equal(t, "invalid user.supervisor: invalid email length", err.GetStatusMessage())
	})

	tt.Run("user.meta", func(t *testing.T) {
		u := setup.UserObjWithID()
		u.SetMeta(&meta.Meta{
			Version: 2,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		})
		err := mgmtuser.ValidateUser(u, validation.IgnoreNothing)
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
			User: setup.UserObjWithID(),
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
		require.Empty(t, resp.GetUsers())
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
