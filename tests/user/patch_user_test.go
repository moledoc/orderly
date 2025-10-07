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
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestPatchUser_InputValidation() {
	tt := s.T()

	name := "name"
	email := user.Email("example@example.com")
	supervisor := user.Email("example.supervisor@example.com")

	tt.Run("EmptyRequest", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), nil)
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("empty", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{})
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})

	tt.Run("InvalidFieldProvided", func(t *testing.T) {
		t.Run("user.meta", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
				User: &user.User{
					Name:       name,
					Email:      email,
					Supervisor: supervisor,
					Meta: &meta.Meta{
						Created: time.Now().UTC(),
						Updated: time.Now().UTC(),
						Version: 2,
					},
				},
			})
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})

	tt.Run("MissingRequiredField", func(t *testing.T) {
		t.Run("user.id", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
				User: &user.User{
					Name:       name,
					Email:      email,
					Supervisor: supervisor,
				},
			})
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.name", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
				User: &user.User{
					ID:         meta.NewID(),
					Email:      email,
					Supervisor: supervisor,
				},
			})
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.email", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
				User: &user.User{
					ID:         meta.NewID(),
					Name:       name,
					Supervisor: supervisor,
				},
			})
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.supervisor", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
				User: &user.User{
					ID:    meta.NewID(),
					Name:  name,
					Email: email,
				},
			})
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})

	tt.Run("InvalidRequiredField", func(t *testing.T) {
		t.Run("user.id.empty", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
				User: &user.User{
					ID:         meta.EmptyID(),
					Name:       name,
					Email:      email,
					Supervisor: supervisor,
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.id.shorter", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
				User: &user.User{
					ID:         meta.NewID()[:10],
					Name:       name,
					Email:      email,
					Supervisor: supervisor,
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.id.longer", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
				User: &user.User{
					ID:         meta.NewID() + meta.NewID(),
					Name:       name,
					Email:      email,
					Supervisor: supervisor,
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.email", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
				User: &user.User{
					ID:         meta.NewID(),
					Name:       name,
					Email:      user.Email("this is not an email"),
					Supervisor: supervisor,
				},
			})
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.supervisor", func(t *testing.T) {
			resp, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{
				User: &user.User{
					ID:         meta.NewID(),
					Name:       name,
					Email:      email,
					Supervisor: user.Email("this is not an email"),
				},
			})
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})
}

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
