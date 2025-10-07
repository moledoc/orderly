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
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestPostUser_InputValidation() {
	tt := s.T()

	name := "name"
	email := user.Email("example@example.com")
	supervisor := user.Email("example.supervisor@example.com")

	tt.Run("EmptyRequest", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), nil)
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("empty", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{})
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})

	tt.Run("InvalidFieldProvided", func(t *testing.T) {
		t.Run("user.id", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
					ID: meta.NewID(),
				},
			})
			defer cleanup.User(t, s.API, resp.GetUser())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.meta", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
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
		t.Run("user.name", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
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
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
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
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
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
		t.Run("user.email", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
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
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
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

func (s *UserSuite) TestPostUser() {
	tt := s.T()

	user := &user.User{
		Name:       "name",
		Email:      user.Email("example@example.com"),
		Supervisor: user.Email("example.supervisor@example.com"),
	}

	resp, err := s.API.PostUser(tt, context.Background(), &request.PostUserRequest{
		User: user,
	})
	defer cleanup.User(tt, s.API, resp.GetUser())
	require.NoError(tt, err)

	opts := []cmp.Option{
		compare.IgnorePaths("User.ID"),
	}

	expected := &response.PostUserResponse{
		User: user,
	}

	compare.RequireEqual(tt, expected, resp, opts...)
}
