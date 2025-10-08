package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestPatchUser() {
	tt := s.T()

	user := setup.MustCreateUserWithCleanup(tt, context.Background(), s.API, &user.User{
		Name:       "name",
		Email:      user.Email("example@example.com"),
		Supervisor: user.Email("example.supervisor@example.com"),
	})

	tt.Run("Name", func(t *testing.T) {
		pathcedUser := *user
		pathcedUser.SetName(pathcedUser.GetName() + "-updated")

		resp, err := s.API.PatchUser(tt, context.Background(), &request.PatchUserRequest{
			User: &pathcedUser,
		})
		require.NoError(t, err)

		opts := []cmp.Option{}

		expected := &response.PatchUserResponse{
			User: &pathcedUser,
		}

		compare.RequireEqual(tt, expected, resp, opts...)
		require.NotEmpty(tt, resp.GetUser().GetMeta())
	})
	tt.Run("Email", func(t *testing.T) {
		pathcedUser := *user
		pathcedUser.SetEmail("example.updated@example.com")

		resp, err := s.API.PatchUser(tt, context.Background(), &request.PatchUserRequest{
			User: &pathcedUser,
		})
		require.NoError(t, err)

		opts := []cmp.Option{}

		expected := &response.PatchUserResponse{
			User: &pathcedUser,
		}

		compare.RequireEqual(tt, expected, resp, opts...)
		require.NotEmpty(tt, resp.GetUser().GetMeta())
	})
	tt.Run("Supervisor", func(t *testing.T) {
		pathcedUser := *user
		pathcedUser.SetSupervisor("example.supervisor.updated@example.com")

		resp, err := s.API.PatchUser(tt, context.Background(), &request.PatchUserRequest{
			User: &pathcedUser,
		})
		require.NoError(t, err)

		opts := []cmp.Option{}

		expected := &response.PatchUserResponse{
			User: &pathcedUser,
		}

		compare.RequireEqual(tt, expected, resp, opts...)
		require.NotEmpty(tt, resp.GetUser().GetMeta())
	})
	tt.Run("Meta", func(t *testing.T) {
		// NOTE: meta is ignored in PATCH call
		pathcedUser := *user
		pathcedUser.SetMeta(&meta.Meta{
			Version: 2,
			Created: time.Now().UTC(),
			Updated: time.Now().UTC(),
		})

		resp, err := s.API.PatchUser(tt, context.Background(), &request.PatchUserRequest{
			User: &pathcedUser,
		})
		require.NoError(t, err)

		opts := []cmp.Option{}

		expected := &response.PatchUserResponse{
			User: user,
		}

		compare.RequireEqual(tt, expected, resp, opts...)
		require.NotEmpty(tt, resp.GetUser().GetMeta())
	})
}

func (s *UserSuite) TestPatchUser_Failed() {
	tt := s.T()

	tt.Run("NotFound", func(t *testing.T) {
		_, err := s.API.PatchUser(tt, context.Background(), &request.PatchUserRequest{

			User: &user.User{
				ID:         meta.NewID(),
				Name:       "name",
				Email:      user.Email("example@example.com"),
				Supervisor: user.Email("example.supervisor@example.com"),
			},
		})
		require.Error(tt, err)
		require.Equal(tt, http.StatusNotFound, err.GetStatusCode(), err)
	})

}
