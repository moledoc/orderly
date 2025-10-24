package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestPostUser() {
	tt := s.T()

	user := setup.UserObj()

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

func (s *UserSuite) TestPostUser_Failed() {
	tt := s.T()

	tt.Run("email.already.exists", func(t *testing.T) {

		userObj := setup.UserObj()
		setup.MustCreateUserWithCleanup(t, context.Background(), s.API, userObj)

		resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
			User: userObj,
		})
		defer cleanup.User(t, s.API, resp.GetUser())
		require.Error(t, err)
		require.Empty(t, resp)
		require.Equal(t, http.StatusConflict, err.GetStatusCode())
		require.Equal(t, fmt.Sprintf("user with email '%v' already exists", userObj.GetEmail()), err.GetStatusMessage())
	})
}

func (s *UserSuite) TestXxx() {
	rootUser := mgmtuser.GetServiceMgmtUser().GetRootUser(context.Background())
	for i := 0; i < 10; i++ {
		obj := setup.UserObj(fmt.Sprintf("%v", i))
		obj.SetSupervisor(rootUser.GetEmail())
		setup.MustCreateUserWithCleanup(s.T(), context.Background(), s.API, obj)
	}
	fmt.Printf("\n")
}
