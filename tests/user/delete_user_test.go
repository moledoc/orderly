package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestDeleteUser_InputValidation() {
	tt := s.T()

	tt.Run("EmptyRequest", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			resp, err := s.API.DeleteUser(t, context.Background(), nil)
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("empty", func(t *testing.T) {
			resp, err := s.API.DeleteUser(t, context.Background(), &request.DeleteUserRequest{})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})

	tt.Run("InvalidRequiredField", func(t *testing.T) {
		t.Run("user.id.empty", func(t *testing.T) {
			resp, err := s.API.DeleteUser(t, context.Background(), &request.DeleteUserRequest{
				ID: meta.EmptyID(),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.id.shorter", func(t *testing.T) {
			resp, err := s.API.DeleteUser(t, context.Background(), &request.DeleteUserRequest{
				ID: meta.NewID()[:10],
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.id.longer", func(t *testing.T) {
			resp, err := s.API.DeleteUser(t, context.Background(), &request.DeleteUserRequest{
				ID: meta.NewID() + meta.NewID(),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})
}

func (s *UserSuite) TestDeleteUser() {
	tt := s.T()

	tt.Run("Existing", func(t *testing.T) {
		userObj := &user.User{
			Name:       "name",
			Email:      user.Email("example@example.com"),
			Supervisor: user.Email("example.supervisor@example.com"),
		}

		user := setup.MustCreateUserWithCleanup(tt, context.Background(), s.API, userObj)

		resp, err := s.API.DeleteUser(tt, context.Background(), &request.DeleteUserRequest{
			ID: user.GetID(),
		})
		require.NoError(tt, err)
		require.Empty(tt, resp)

		_, err = s.API.GetUserByID(tt, context.Background(), &request.GetUserByIDRequest{
			ID: user.GetID(),
		})
		require.Error(tt, err)
		require.Equal(tt, http.StatusNotFound, err.GetStatusCode(), err)
	})
	tt.Run("Non-Existing", func(t *testing.T) {
		id := meta.NewID()
		resp, err := s.API.DeleteUser(tt, context.Background(), &request.DeleteUserRequest{
			ID: id,
		})
		require.NoError(tt, err)
		require.Empty(tt, resp)

		_, err = s.API.GetUserByID(tt, context.Background(), &request.GetUserByIDRequest{
			ID: id,
		})
		require.Error(tt, err)
		require.Equal(tt, http.StatusNotFound, err.GetStatusCode(), err)
	})
}
