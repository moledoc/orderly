package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestGetUserByID_InputValidation() {
	t := s.T()

	t.Run("EmptyRequest", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			resp, err := s.API.GetUserByID(t, context.Background(), nil)
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("empty", func(t *testing.T) {
			resp, err := s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})

	t.Run("InvalidRequiredField", func(t *testing.T) {
		t.Run("user.id.empty", func(t *testing.T) {
			resp, err := s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{
				ID: utils.Ptr(meta.ID("")),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.id.shorter", func(t *testing.T) {
			resp, err := s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{
				ID: utils.Ptr(meta.ID(utils.RandAlphanum()[:10])),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.id.longer", func(t *testing.T) {
			resp, err := s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{
				ID: utils.Ptr(meta.ID(utils.RandAlphanum() + utils.RandAlphanum())),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})
}

func (s *UserSuite) TestGetUserByID() {
	t := s.T()

	userObj := &user.User{
		Name:       utils.Ptr("name"),
		Email:      utils.Ptr(user.Email("example@example.com")),
		Supervisor: utils.Ptr(user.Email("example.supervisor@example.com")),
	}

	respPost, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
		User: userObj,
	})
	require.NoError(t, err)
	expectedUser := respPost.GetUser()

	resp, err := s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{
		ID: utils.Ptr(expectedUser.GetID()),
	})
	require.NoError(t, err)

	opts := []cmp.Option{
		compare.IgnorePath("User.Meta"),
		compare.ComparerUser(),
	}

	expected := &response.GetUserByIDResponse{
		User: expectedUser,
	}
	compare.RequireEqual(t, expected, resp, opts...)
}
