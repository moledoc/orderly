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

func (s *UserSuite) TestDeleteUser() {
	tt := s.T()

	tt.Run("Existing", func(t *testing.T) {
		userObj := &user.User{
			Name:         "name",
			Email:        user.Email("example@example.com"),
			SupervisorID: user.Email("example.supervisor@example.com"),
		}

		user := setup.MustCreateUserWithCleanup(t, context.Background(), s.API, userObj)

		resp, err := s.API.DeleteUser(t, context.Background(), &request.DeleteUserRequest{
			ID: user.GetID(),
		})
		require.NoError(t, err)
		require.Empty(t, resp)

		_, err = s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{
			ID: user.GetID(),
		})
		require.Error(t, err)
		require.Equal(t, http.StatusNotFound, err.GetStatusCode(), err)
	})
	tt.Run("Non-Existing", func(t *testing.T) {
		id := meta.NewID()
		resp, err := s.API.DeleteUser(t, context.Background(), &request.DeleteUserRequest{
			ID: id,
		})
		require.NoError(t, err)
		require.Empty(t, resp)

		_, err = s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{
			ID: id,
		})
		require.Error(t, err)
		require.Equal(t, http.StatusNotFound, err.GetStatusCode(), err)
	})
}
