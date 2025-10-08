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
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

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
	require.NotEmpty(tt, resp.GetUser().GetMeta())
}

func (s *UserSuite) TestGetUserByID_Failed() {
	tt := s.T()

	tt.Run("NotFound", func(t *testing.T) {
		_, err := s.API.GetUserByID(tt, context.Background(), &request.GetUserByIDRequest{
			ID: meta.NewID(),
		})
		require.Error(tt, err)
		require.Equal(tt, http.StatusNotFound, err.GetStatusCode(), err)
		require.Equal(t, "not found", err.GetStatusMessage())
	})
}
