package tests

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/stretchr/testify/require"
)

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
		compare.IgnoreMeta,
	}

	expected := &response.PostUserResponse{
		User: user,
	}

	compare.RequireEqual(tt, expected, resp, opts...)
}
