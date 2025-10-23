package mgmtuser

import (
	"context"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/repository"
)

type ServiceMgmtUserAPI interface {
	PostUser(ctx context.Context, req *request.PostUserRequest) (*response.PostUserResponse, errwrap.Error)
	GetUserByID(ctx context.Context, req *request.GetUserByIDRequest) (*response.GetUserByIDResponse, errwrap.Error)
	GetUsers(ctx context.Context, req *request.GetUsersRequest) (*response.GetUsersResponse, errwrap.Error)
	PatchUser(ctx context.Context, req *request.PatchUserRequest) (*response.PatchUserResponse, errwrap.Error)
	DeleteUser(ctx context.Context, req *request.DeleteUserRequest) (*response.DeleteUserResponse, errwrap.Error)
}

type serviceMgmtUser struct {
	Repository repository.RepositoryUserAPI
}

var (
	_   ServiceMgmtUserAPI = (*serviceMgmtUser)(nil)
	svc ServiceMgmtUserAPI = nil
)

func GetServiceMgmtUser() ServiceMgmtUserAPI {
	return svc
}

func NewServiceMgmtUser(repo repository.RepositoryUserAPI) ServiceMgmtUserAPI {
	svc = &serviceMgmtUser{
		Repository: repo,
	}
	return svc
}
