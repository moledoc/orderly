package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestGetUserSubOrdinates() {
	tt := s.T()

	supervisor := setup.MustCreateUserWithCleanup(tt, context.Background(), s.API, &user.User{
		Name:       "name",
		Email:      user.Email("example.supervisor.1@example.com"),
		Supervisor: user.Email("example.supervisor.0@example.com"),
	})

	createSubOrdinates := func(t *testing.T, count int) []*user.User {
		if count == 0 {
			return nil
		}
		users := make([]*user.User, count)
		for i := 1; i <= count; i++ {
			userObj := &user.User{
				Name:       fmt.Sprintf("name-%d", count),
				Email:      user.Email(fmt.Sprintf("example.%s.%d.%d@example.com", t.Name(), count, i)),
				Supervisor: supervisor.GetEmail(),
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
				ID: supervisor.GetID(),
			})
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.SorterUser(compare.SortUserByID),
			}

			expected := &response.GetUserSubOrdinatesResponse{
				SubOrdinates: users,
			}
			compare.RequireEqual(t, expected, resp, opts...)
		})
	}
}
