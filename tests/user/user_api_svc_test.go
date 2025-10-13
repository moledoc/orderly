package tests

import (
	"context"
	"flag"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/moledoc/orderly/pkg/flags"
	"github.com/moledoc/orderly/tests/api"
	"github.com/stretchr/testify/suite"
)

type UserAPISvc struct { // NOTE: tests service layer methods
	Svc mgmtuser.ServiceMgmtUserAPI
}

func NewUserAPISvc() *UserAPISvc {
	// TODO: local vs db
	return &UserAPISvc{
		Svc: mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()),
	}
}

var (
	_ api.User = (*UserAPISvc)(nil)
)

func TestUserSvcSuite(t *testing.T) {
	flag.Parse()
	if flags.TestMode(*flags.ModeFlag) != flags.FuncTest {
		return
	}
	t.Run("UserAPISvc", func(t *testing.T) {
		suite.Run(t, &UserSuite{
			API: NewUserAPISvc(),
		})
	})
}

func TestUserSvcPerformanceSuite(t *testing.T) {
	flag.Parse()
	if flags.TestMode(*flags.ModeFlag) != flags.PerfTest {
		return
	}
	t.Run("UserAPISvcPerformance", func(t *testing.T) {
		suite.Run(t, &UserPerformanceSuite{
			API: NewUserAPISvc(),
		})
	})
}

func (api *UserAPISvc) PostUser(t *testing.T, ctx context.Context, req *request.PostUserRequest) (*response.PostUserResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PostUser(ctx, req)
}

func (api *UserAPISvc) GetUserByID(t *testing.T, ctx context.Context, req *request.GetUserByIDRequest) (*response.GetUserByIDResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.GetUserByID(ctx, req)
}
func (api *UserAPISvc) GetUsers(t *testing.T, ctx context.Context, req *request.GetUsersRequest) (*response.GetUsersResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.GetUsers(ctx, req)
}
func (api *UserAPISvc) GetUserSubOrdinates(t *testing.T, ctx context.Context, req *request.GetUserSubOrdinatesRequest) (*response.GetUserSubOrdinatesResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.GetUserSubOrdinates(ctx, req)
}
func (api *UserAPISvc) PatchUser(t *testing.T, ctx context.Context, req *request.PatchUserRequest) (*response.PatchUserResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PatchUser(ctx, req)
}
func (api *UserAPISvc) DeleteUser(t *testing.T, ctx context.Context, req *request.DeleteUserRequest) (*response.DeleteUserResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.DeleteUser(ctx, req)
}
