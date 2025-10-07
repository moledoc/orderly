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
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestGetUserByID_InputValidation() {
	tt := s.T()

	tt.Run("EmptyRequest", func(t *testing.T) {
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

	tt.Run("InvalidRequiredField", func(t *testing.T) {
		t.Run("user.id.empty", func(t *testing.T) {
			resp, err := s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{
				ID: meta.ID(""),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.id.shorter", func(t *testing.T) {
			resp, err := s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{
				ID: meta.ID(utils.RandAlphanum()[:10]),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.id.longer", func(t *testing.T) {
			resp, err := s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{
				ID: meta.ID(utils.RandAlphanum() + utils.RandAlphanum()),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})
}

func (s *UserSuite) TestGetUserByID() {
	tt := s.T()

	userObj := &user.User{
		Name:       "name",
		Email:      user.Email("example@example.com"),
		Supervisor: user.Email("example.supervisor@example.com"),
	}

	user := setup.MustCreateUserWithCleanup(tt, context.Background(), s.API, userObj)

	resp, err := s.API.GetUserByID(tt, context.Background(), &request.GetUserByIDRequest{
		ID: user.GetID(),
	})
	require.NoError(tt, err)

	opts := []cmp.Option{}

	expected := &response.GetUserByIDResponse{
		User: user,
	}
	compare.RequireEqual(tt, expected, resp, opts...)
}

func (s *UserSuite) TestGetUserByID_Failed() {
	tt := s.T()

	tt.Run("NotFound", func(t *testing.T) {
		_, err := s.API.GetUserByID(tt, context.Background(), &request.GetUserByIDRequest{
			ID: meta.ID(utils.RandAlphanum()),
		})
		require.Error(tt, err)
		require.Equal(tt, http.StatusNotFound, err.GetStatusCode(), err)
	})
}
