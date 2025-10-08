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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestGetUsers() {

	createUsers := func(t *testing.T, count int) []*user.User {
		if count == 0 {
			return nil
		}
		users := make([]*user.User, count)
		for i := 1; i <= count; i++ {
			userObj := &user.User{
				Name:       fmt.Sprintf("name-%d", count),
				Email:      user.Email(fmt.Sprintf("example.%d@example.com", count)),
				Supervisor: user.Email(fmt.Sprintf("example.supervisor.%d@example.com", count)),
			}

			user := setup.MustCreateUserWithCleanup(t, context.Background(), s.API, userObj)

			users[i-1] = user
		}
		return users
	}

	for _, i := range []int{0, 1, 10} {
		s.T().Run(fmt.Sprintf("count.%v", i), func(t *testing.T) {
			users := createUsers(t, i)

			resp, err := s.API.GetUsers(t, context.Background(), &request.GetUsersRequest{})
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.SorterUser(compare.SortUserByID),
			}

			expected := &response.GetUsersResponse{
				Users: users,
			}
			compare.RequireEqual(t, expected, resp, opts...)
			for _, u := range resp.GetUsers() {
				assert.NotEmpty(t, u.GetMeta())
			}
		})
	}
}
