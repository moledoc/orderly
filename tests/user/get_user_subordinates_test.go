package tests

import (
	"context"
	"fmt"
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

func (s *UserSuite) TestGetUserSubOrdinates_InputValidation() {
	tt := s.T()

	tt.Run("EmptyRequest", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			resp, err := s.API.GetUserSubOrdinates(t, context.Background(), nil)
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("empty", func(t *testing.T) {
			resp, err := s.API.GetUserSubOrdinates(t, context.Background(), &request.GetUserSubOrdinatesRequest{})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})

	tt.Run("InvalidRequiredField", func(t *testing.T) {
		t.Run("user.id.empty", func(t *testing.T) {
			resp, err := s.API.GetUserSubOrdinates(t, context.Background(), &request.GetUserSubOrdinatesRequest{
				ID: utils.Ptr(meta.ID("")),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.id.shorter", func(t *testing.T) {
			resp, err := s.API.GetUserSubOrdinates(t, context.Background(), &request.GetUserSubOrdinatesRequest{
				ID: utils.Ptr(meta.ID(utils.RandAlphanum()[:10])),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("user.id.longer", func(t *testing.T) {
			resp, err := s.API.GetUserSubOrdinates(t, context.Background(), &request.GetUserSubOrdinatesRequest{
				ID: utils.Ptr(meta.ID(utils.RandAlphanum() + utils.RandAlphanum())),
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})
}

func (s *UserSuite) TestGetUserSubOrdinates() {
	tt := s.T()

	supervisor := setup.MustCreateUserWithCleanup(tt, context.Background(), s.API, &user.User{
		Name:       utils.Ptr("name"),
		Email:      utils.Ptr(user.Email("example.supervisor.1@example.com")),
		Supervisor: utils.Ptr(user.Email("example.supervisor.0@example.com")),
	})

	createSubOrdinates := func(t *testing.T, count int) []*user.User {
		if count == 0 {
			return nil
		}
		users := make([]*user.User, count)
		for i := 1; i <= count; i++ {
			userObj := &user.User{
				Name:       utils.Ptr(fmt.Sprintf("name-%d", count)),
				Email:      utils.Ptr(user.Email(fmt.Sprintf("example.%d@example.com", count))),
				Supervisor: utils.Ptr(supervisor.GetEmail()),
			}

			user := setup.MustCreateUserWithCleanup(t, context.Background(), s.API, userObj)

			users[i-1] = user
		}
		return users
	}

	for _, i := range []int{0, 1, 10} {
		tt.Run(fmt.Sprintf("count.%d", i), func(t *testing.T) {

			users := createSubOrdinates(t, i)

			resp, err := s.API.GetUserSubOrdinates(t, context.Background(), &request.GetUserSubOrdinatesRequest{
				ID: utils.Ptr(supervisor.GetID()),
			})
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.IgnorePath("User.Meta"),
				compare.SorterUser(compare.SortUserByID),
				compare.ComparerUser(),
			}

			expected := &response.GetUserSubOrdinatesResponse{
				SubOrdinates: users,
			}
			compare.RequireEqual(t, expected, resp, opts...)
		})
	}
}
