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
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestPostUser_InputValidation() {
	tt := s.T()

	name := utils.Ptr("name")
	email := utils.Ptr(user.Email("example@example.com"))
	supervisor := utils.Ptr(user.Email("example.supervisor@example.com"))

	tt.Run("EmptyRequest", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), nil)
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("empty", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})

	tt.Run("InvalidFieldProvided", func(t *testing.T) {
		t.Run("user.id", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
					ID: utils.Ptr(meta.ID(utils.RandAlphanum())),
				},
			})
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
					Email:      utils.Ptr(user.Email("this is not an email")),
					Supervisor: supervisor,
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.supervisor", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
					Name:       name,
					Email:      email,
					Supervisor: utils.Ptr(user.Email("this is not an email")),
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})
}

func (s *UserSuite) TestPostUser_CreateUser() {
	t := s.T()

	user := &user.User{
		Name:       utils.Ptr("name"),
		Email:      utils.Ptr(user.Email("example@example.com")),
		Supervisor: utils.Ptr(user.Email("example.supervisor@example.com")),
	}

	resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
		User: user,
	})
	require.NoError(t, err)

	opts := []cmp.Option{
		compare.IgnorePath("User.ID", "User.Meta"),
		compare.ComparerUser(),
	}
	expected := &response.PostUserResponse{
		User: user,
	}
	compare.RequireEqual(t, expected, resp, opts...)

	cleanup.User(t, s.API, resp.GetUser())
}
