package mgmtuser

import (
	"context"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/middleware"
	"github.com/moledoc/orderly/internal/repository"
)

type ServiceMgmtUserAPI interface {
	GetRootUser(context.Context) *user.User
	////
	PostUser(ctx context.Context, req *request.PostUserRequest) (*response.PostUserResponse, errwrap.Error)
	GetUserByID(ctx context.Context, req *request.GetUserByIDRequest) (*response.GetUserByIDResponse, errwrap.Error)
	GetUsers(ctx context.Context, req *request.GetUsersRequest) (*response.GetUsersResponse, errwrap.Error)
	PatchUser(ctx context.Context, req *request.PatchUserRequest) (*response.PatchUserResponse, errwrap.Error)
	DeleteUser(ctx context.Context, req *request.DeleteUserRequest) (*response.DeleteUserResponse, errwrap.Error)
}

type serviceMgmtUser struct {
	RootUser   *user.User
	Repository repository.RepositoryUserAPI
}

var (
	_   ServiceMgmtUserAPI = (*serviceMgmtUser)(nil)
	svc ServiceMgmtUserAPI = nil
)

func postRootUser(ctx context.Context, repo repository.RepositoryUserAPI) (*user.User, errwrap.Error) {
	now := time.Now().UTC()
	u := &user.User{
		ID:         meta.NewID(),
		Name:       "Root",
		Email:      "root@root.com",
		Supervisor: "root@root.com",
		Meta: &meta.Meta{
			Version: 1,
			Created: now,
			Updated: now,
		},
	}

	u, err := repo.Write(ctx, u)
	if err != nil {
		return nil, middleware.AddTraceToErrFromCtx(err, ctx)
	}
	return u, nil
}

func GetServiceMgmtUser() ServiceMgmtUserAPI {
	return svc
}

func NewServiceMgmtUser(repo repository.RepositoryUserAPI) ServiceMgmtUserAPI {
	u, err := postRootUser(context.Background(), repo)
	if err != nil {
		panic(err)
	}
	svc = &serviceMgmtUser{
		RootUser:   u,
		Repository: repo,
	}
	return svc
}
