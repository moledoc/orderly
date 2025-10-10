package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestPatchUser() {
	tt := s.T()

	u := setup.MustCreateUserWithCleanup(tt, context.Background(), s.API, &user.User{
		Name:       "name",
		Email:      user.Email("example@example.com"),
		Supervisor: user.Email("example.supervisor@example.com"),
	})

	changes := []struct {
		Name string
		f    func() *request.PatchUserRequest
	}{
		{
			Name: "name",
			f: func() *request.PatchUserRequest {
				u.SetName(u.GetName() + "-updated")
				u.GetMeta().VersionIncr()

				return &request.PatchUserRequest{
					User: &user.User{
						ID:   u.GetID(),
						Name: u.GetName(),
					},
				}
			},
		},
		{
			Name: "email",
			f: func() *request.PatchUserRequest {
				u.SetEmail("example.updated@example.com")
				u.GetMeta().VersionIncr()

				return &request.PatchUserRequest{
					User: &user.User{
						ID:    u.GetID(),
						Email: u.GetEmail(),
					},
				}
			},
		},
		{
			Name: "supervisor",
			f: func() *request.PatchUserRequest {
				u.SetSupervisor("example.supervisor.updated@example.com")
				u.GetMeta().VersionIncr()

				return &request.PatchUserRequest{
					User: &user.User{
						ID:         u.GetID(),
						Supervisor: u.GetSupervisor(),
					},
				}
			},
		},
		{
			Name: "meta",
			f: func() *request.PatchUserRequest {

				return &request.PatchUserRequest{
					User: &user.User{
						ID: u.GetID(),
						Meta: &meta.Meta{
							Version: 2,
							Created: time.Now().UTC(),
							Updated: time.Now().UTC(),
						},
					},
				}
			},
		},
	}

	for _, change := range changes {
		tt.Run(change.Name, func(t *testing.T) {
			req := change.f()

			respPatch, err := s.API.PatchUser(t, context.Background(), req)
			require.NoError(t, err)

			opts := []cmp.Option{
				compare.IgnorePaths("User.Meta.Updated"),
			}
			expectedPatch := &response.PatchUserResponse{
				User: u,
			}
			compare.RequireEqual(t, expectedPatch, respPatch, opts...)

			respGet, err := s.API.GetUserByID(t, context.Background(), &request.GetUserByIDRequest{
				ID: u.GetID(),
			})
			require.NoError(t, err)

			expectedGet := &response.GetUserByIDResponse{
				User: u,
			}
			compare.RequireEqual(t, expectedGet, respGet, opts...)

		})
	}
}

func (s *UserSuite) TestPatchUser_Failed() {
	tt := s.T()

	tt.Run("NotFound", func(t *testing.T) {
		_, err := s.API.PatchUser(t, context.Background(), &request.PatchUserRequest{

			User: &user.User{
				ID:         meta.NewID(),
				Name:       "name",
				Email:      user.Email("example@example.com"),
				Supervisor: user.Email("example.supervisor@example.com"),
			},
		})
		require.Error(t, err)
		require.Equal(t, http.StatusNotFound, err.GetStatusCode(), err)
	})

}
