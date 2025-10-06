package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserAPI interface {
	PostUser(t *testing.T, ctx context.Context, req *request.PostUserRequest) (*response.PostUserResponse, errwrap.Error)
	GetUserByID(t *testing.T, ctx context.Context, req *request.GetUserByIDRequest) (*response.GetUserByIDResponse, errwrap.Error)
	GetUsers(t *testing.T, ctx context.Context, req *request.GetUsersRequest) (*response.GetUsersResponse, errwrap.Error)
	GetUserSubOrdinates(t *testing.T, ctx context.Context, req *request.GetUserSubOrdinatesRequest) (*response.GetUserSubOrdinatesResponse, errwrap.Error)
	PatchUser(t *testing.T, ctx context.Context, req *request.PatchUserRequest) (*response.PatchUserResponse, errwrap.Error)
	DeleteUser(t *testing.T, ctx context.Context, req *request.DeleteUserRequest) (*response.DeleteUserResponse, errwrap.Error)
}

func (api *UserAPISvc) PostUser(t *testing.T, ctx context.Context, req *request.PostUserRequest) (*response.PostUserResponse, errwrap.Error) {
	t.Helper()
	return api.Svc.PostUser(ctx, req)
}

func (*UserAPISvc) GetUserByID(t *testing.T, ctx context.Context, req *request.GetUserByIDRequest) (*response.GetUserByIDResponse, errwrap.Error) {
	t.Helper()
	return nil, errwrap.NewError(http.StatusServiceUnavailable, "TODO: implement")
}
func (*UserAPISvc) GetUsers(t *testing.T, ctx context.Context, req *request.GetUsersRequest) (*response.GetUsersResponse, errwrap.Error) {
	t.Helper()
	return nil, errwrap.NewError(http.StatusServiceUnavailable, "TODO: implement")
}
func (*UserAPISvc) GetUserSubOrdinates(t *testing.T, ctx context.Context, req *request.GetUserSubOrdinatesRequest) (*response.GetUserSubOrdinatesResponse, errwrap.Error) {
	t.Helper()
	return nil, errwrap.NewError(http.StatusServiceUnavailable, "TODO: implement")
}
func (*UserAPISvc) PatchUser(t *testing.T, ctx context.Context, req *request.PatchUserRequest) (*response.PatchUserResponse, errwrap.Error) {
	t.Helper()
	return nil, errwrap.NewError(http.StatusServiceUnavailable, "TODO: implement")
}
func (*UserAPISvc) DeleteUser(t *testing.T, ctx context.Context, req *request.DeleteUserRequest) (*response.DeleteUserResponse, errwrap.Error) {
	t.Helper()
	return nil, errwrap.NewError(http.StatusServiceUnavailable, "TODO: implement")
}

func (api *UserAPIReq) PostUser(t *testing.T, ctx context.Context, req *request.PostUserRequest) (*response.PostUserResponse, errwrap.Error) {
	t.Helper()
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errwrap.NewError(http.StatusBadRequest, "marshaling request failed: %s", err)
	}
	respHttp, err := api.HttpClient.Post("http://localhost:8080/user", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "sending request failed: %s", err)
	}
	var resp response.PostUserResponse
	if err := json.NewDecoder(respHttp.Body).Decode(&resp); err != nil {
		return nil, errwrap.NewError(http.StatusInternalServerError, "unmarshaling response failed: %s", err)
	}
	return &resp, nil
}

func (*UserAPIReq) GetUserByID(t *testing.T, ctx context.Context, req *request.GetUserByIDRequest) (*response.GetUserByIDResponse, errwrap.Error) {
	t.Helper()
	return nil, errwrap.NewError(http.StatusServiceUnavailable, "TODO: implement")
}
func (*UserAPIReq) GetUsers(t *testing.T, ctx context.Context, req *request.GetUsersRequest) (*response.GetUsersResponse, errwrap.Error) {
	t.Helper()
	return nil, errwrap.NewError(http.StatusServiceUnavailable, "TODO: implement")
}
func (*UserAPIReq) GetUserSubOrdinates(t *testing.T, ctx context.Context, req *request.GetUserSubOrdinatesRequest) (*response.GetUserSubOrdinatesResponse, errwrap.Error) {
	t.Helper()
	return nil, errwrap.NewError(http.StatusServiceUnavailable, "TODO: implement")
}
func (*UserAPIReq) PatchUser(t *testing.T, ctx context.Context, req *request.PatchUserRequest) (*response.PatchUserResponse, errwrap.Error) {
	t.Helper()
	return nil, errwrap.NewError(http.StatusServiceUnavailable, "TODO: implement")
}
func (*UserAPIReq) DeleteUser(t *testing.T, ctx context.Context, req *request.DeleteUserRequest) (*response.DeleteUserResponse, errwrap.Error) {
	t.Helper()
	return nil, errwrap.NewError(http.StatusServiceUnavailable, "TODO: implement")
}

// DiffReporter is a simple custom reporter that only records differences
// detected during comparison.
type DiffReporter struct {
	path  cmp.Path
	diffs []string
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\t-: %+v\n\t+: %+v\n", r.path, vx, vy))
	}
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}

var (
	IgnoreID = cmp.FilterPath(func(pathsteps cmp.Path) bool {
		last := pathsteps.Last().String()
		return last == ".ID"
	}, cmp.Ignore())
	IgnoreMeta = cmp.FilterPath(func(pathsteps cmp.Path) bool {
		last := pathsteps.Last().String()
		return last == ".Meta"
	}, cmp.Ignore())
)

var (
	IgnorePath = func(paths ...string) cmp.Option {
		return cmp.FilterPath(func(ppaths cmp.Path) bool {
			for _, path := range paths {
				if ppaths[1:].String() == path {
					return true
				}
			}
			return false
		}, cmp.Ignore())
	}
)

var (
	SorterString = cmpopts.SortSlices(func(a string, b string) bool {
		return a < b
	})

	ComparerUser = func(comparers ...func(a *user.User, b *user.User) bool) cmp.Option {
		return cmp.Comparer(func(a *user.User, b *user.User) bool {
			for _, comparer := range comparers {
				if !comparer(a, b) {
					return false
				}
			}
			return true
		})
	}
)

func AssertEqual(t *testing.T, expected any, actual any, opts ...cmp.Option) {
	opts = append(opts, IgnoreMeta)
	var r DiffReporter
	opts = append(opts, cmp.Reporter(&r))
	assert.True(t, cmp.Equal(expected, actual, opts...), r.String())
}

func RequireEqual(t *testing.T, expected any, actual any, opts ...cmp.Option) {
	opts = append(opts, IgnoreMeta)
	var r DiffReporter
	opts = append(opts, cmp.Reporter(&r))
	require.True(t, cmp.Equal(expected, actual, opts...), r.String())
}

func (s *UserSuite) TestPostUser_InputValidation() {
	t := s.T()

	name := utils.Ptr("name")
	email := utils.Ptr(user.Email("example@example.com"))
	supervisor := utils.Ptr(user.Email("example.supervisor@example.com"))

	t.Run("EmptyRequest", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), nil)
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode())
		})
		t.Run("empty", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode())
		})
	})

	t.Run("InvalidFieldProvided", func(t *testing.T) {
		t.Run("user.id", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
					ID: utils.Ptr(meta.ID(utils.RandAlphanum())),
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode())
		})
		t.Run("user.meta", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
					Name:       name,
					Email:      email,
					Supervisor: supervisor,
					Meta: &meta.Meta{
						Created: time.Now().UTC(),
						Updated: time.Now().UTC(),
						Version: 2,
					},
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode())
		})
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		t.Run("user.name", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
					Email:      email,
					Supervisor: supervisor,
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode())
		})
		t.Run("user.email", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
					Name:       name,
					Supervisor: supervisor,
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode())
		})
		t.Run("user.supervisor", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
					Name:  name,
					Email: email,
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode())
		})
	})

	t.Run("InvalidRequiredField", func(t *testing.T) {
		t.Run("user.email", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
					Name:       name,
					Email:      utils.Ptr(user.Email("this is not an email")),
					Supervisor: supervisor,
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode())
		})
		t.Run("user.supervisor", func(t *testing.T) {
			resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
				User: &user.User{
					Name:       name,
					Email:      email,
					Supervisor: utils.Ptr(user.Email("this is not an email")),
				},
			})
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode())
		})
	})
}

func (s *UserSuite) TestPostUser_CreateUser() {
	t := s.T()

	user := &user.User{
		Name:       utils.Ptr("name"),
		Email:      utils.Ptr(user.Email("example@example.com")),
		Supervisor: utils.Ptr(user.Email("example.supervisor@example.com")),
	}

	resp, err := s.API.PostUser(t, context.Background(), &request.PostUserRequest{
		User: user,
	})
	require.NoError(t, err)

	opts := []cmp.Option{
		IgnorePath("User.ID", "User.Meta"),
		ComparerUser(),
	}
	expected := &response.PostUserResponse{
		User: user,
	}
	RequireEqual(t, expected, resp, opts...)
}
