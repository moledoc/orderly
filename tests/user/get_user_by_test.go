package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *UserSuite) TestGetUserBy() {
	tt := s.T()

	user := setup.MustCreateUserWithCleanup(tt, context.Background(), s.API, setup.UserObj())

	cases := []struct {
		name string
		req  request.GetUserByRequest
	}{
		{
			name: "by-id",
			req: request.GetUserByRequest{
				ID: user.GetID(),
			},
		},
		{
			name: "by-email",
			req: request.GetUserByRequest{
				Email: user.GetEmail(),
			},
		},
		{
			name: "by-supervisor",
			req: request.GetUserByRequest{
				Supervisor: user.GetSupervisor(),
			},
		},
		{
			name: "by-id-email",
			req: request.GetUserByRequest{
				ID:    user.GetID(),
				Email: user.GetEmail(),
			},
		},
		{
			name: "by-id-supervisor",
			req: request.GetUserByRequest{
				ID:         user.GetID(),
				Supervisor: user.GetSupervisor(),
			},
		},
		{
			name: "by-email-supervisor",
			req: request.GetUserByRequest{
				Email:      user.GetEmail(),
				Supervisor: user.GetSupervisor(),
			},
		},
		{
			name: "by-id-email-supervisor",
			req: request.GetUserByRequest{
				ID:         user.GetID(),
				Email:      user.GetEmail(),
				Supervisor: user.GetSupervisor(),
			},
		},
	}

	for _, cse := range cases {
		tt.Run(cse.name, func(t *testing.T) {

			resp, err := s.API.GetUserBy(t, context.Background(), &cse.req)
			require.NoError(t, err)

			opts := []cmp.Option{}

			expected := &response.GetUserByResponse{
				User: user,
			}
			compare.RequireEqual(t, expected, resp, opts...)
		})
	}
}

func (s *UserSuite) TestGetUserBy_Failed() {
	tt := s.T()

	user := setup.MustCreateUserWithCleanup(tt, context.Background(), s.API, setup.UserObj())

	cases := []struct {
		name string
		req  request.GetUserByRequest
	}{

		////////////////

		{
			name: "not-found-by-id.id",
			req: request.GetUserByRequest{
				ID: meta.NewID(),
			},
		},
		{
			name: "not-found-by-id.id-email",
			req: request.GetUserByRequest{
				ID:    meta.NewID(),
				Email: user.GetEmail(),
			},
		},
		{
			name: "not-found-by-id.id-supervisor",
			req: request.GetUserByRequest{
				ID:         meta.NewID(),
				Supervisor: user.GetSupervisor(),
			},
		},
		{
			name: "not-found-by-id.id-email-supervisor",
			req: request.GetUserByRequest{
				ID:         meta.NewID(),
				Email:      user.GetEmail(),
				Supervisor: user.GetSupervisor(),
			},
		},

		////////////////

		{
			name: "not-found-by-email.email",
			req: request.GetUserByRequest{
				Email: "not.found.email@email.com",
			},
		},
		{
			name: "not-found-by-email.id-email",
			req: request.GetUserByRequest{
				ID:    user.GetID(),
				Email: "not.found.email@email.com",
			},
		},
		{
			name: "not-found-by-email.id-supervisor",
			req: request.GetUserByRequest{
				Email:      "not.found.email@email.com",
				Supervisor: user.GetSupervisor(),
			},
		},
		{
			name: "not-found-by-email.id-email-supervisor",
			req: request.GetUserByRequest{
				ID:         user.GetID(),
				Email:      "not.found.email@email.com",
				Supervisor: user.GetSupervisor(),
			},
		},

		////////////////

		{
			name: "not-found-by-supervisor.supervisor",
			req: request.GetUserByRequest{
				Supervisor: "not.found.email@email.com",
			},
		},
		{
			name: "not-found-by-supervisor.id-supervisor",
			req: request.GetUserByRequest{
				ID:         user.GetID(),
				Supervisor: "not.found.email@email.com",
			},
		},
		{
			name: "not-found-by-supervisor.email-supervisor",
			req: request.GetUserByRequest{
				Email:      user.GetEmail(),
				Supervisor: "not.found.email@email.com",
			},
		},
		{
			name: "not-found-by-supervisor.id-email-supervisor",
			req: request.GetUserByRequest{
				ID:         user.GetID(),
				Supervisor: "not.found.email@email.com",
				Email:      user.GetEmail(),
			},
		},
	}

	for _, cse := range cases {
		tt.Run(cse.name, func(t *testing.T) {

			resp, err := s.API.GetUserBy(t, context.Background(), &cse.req)
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusNotFound, err.GetStatusCode())
			require.Equal(t, "not found", err.GetStatusMessage())

		})
	}
}
