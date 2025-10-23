package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/performance"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *UserPerformanceSuite) TestPerformance_PostUser() {
	setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
		return context.Background, &request.PostUserRequest{
			User: setup.UserObj(utils.RandAlphanum()),
		}, nil
	}
	tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
		if errReq != nil {
			return nil, errReq
		}
		return s.API.PostUser(s.T(), ctx, req.(*request.PostUserRequest))
	}
	plan := performance.Plan{
		T:               s.T(),
		RPS:             10,
		DurationSec:     10,
		RampDurationSec: 10,
		Setup:           setup,
		Test:            tst,
		NFR: performance.NFRs{
			P50: 50 * time.Millisecond,
			P90: 90 * time.Millisecond,
			P95: 95 * time.Millisecond,
			P99: 99 * time.Millisecond,
		},
		Notes: []string{"test", "test", "test"},
	}
	report, _, _ := plan.Run()
	fmt.Printf("%+v\n", report)
}

func (s *UserPerformanceSuite) TestPerformance_GetUserByID() {
	setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
		u := setup.MustCreateUserWithCleanup(s.T(), context.Background(), s.API, setup.UserObj(utils.RandAlphanum()))
		return context.Background, &request.GetUserByIDRequest{
			ID: u.GetID(),
		}, nil
	}
	tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
		if errReq != nil {
			return nil, errReq
		}
		return s.API.GetUserByID(s.T(), ctx, req.(*request.GetUserByIDRequest))
	}
	plan := performance.Plan{
		T:               s.T(),
		RPS:             10,
		DurationSec:     10,
		RampDurationSec: 10,
		Setup:           setup,
		Test:            tst,
		NFR: performance.NFRs{
			P50: 50 * time.Millisecond,
			P90: 90 * time.Millisecond,
			P95: 95 * time.Millisecond,
			P99: 99 * time.Millisecond,
		},
		Notes: []string{"test", "test", "test"},
	}
	report, _, _ := plan.Run()
	fmt.Printf("%+v\n", report)
}

func (s *UserPerformanceSuite) TestPerformance_GetUsers() {
	for _, userCount := range []int{10, 100, 1000} {
		s.T().Run(fmt.Sprintf("%v", userCount), func(t *testing.T) {
			for i := 0; i < userCount; i++ {
				setup.MustCreateUserWithCleanup(t, context.Background(), s.API, setup.UserObj(fmt.Sprintf("%v-%v", userCount, i)))
			}
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				return context.Background, &request.GetUsersRequest{}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.GetUsers(t, ctx, req.(*request.GetUsersRequest))
			}
			checkLen := func(ctx context.Context, resp any, err errwrap.Error) {
				require.Len(t, resp.(*response.GetUsersResponse).GetUsers(), userCount)
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(userCount),
				DurationSec:     10,
				RampDurationSec: 10,
				Setup:           setup,
				Test:            tst,
				Assert:          checkLen,
				NFR: performance.NFRs{
					P50: 50 * time.Millisecond,
					P90: 90 * time.Millisecond,
					P95: 95 * time.Millisecond,
					P99: 99 * time.Millisecond,
				},
				Notes: []string{"test", "test", "test"},
			}
			report, _, _ := plan.Run()
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}

func (s *UserPerformanceSuite) TestPerformance_GetUsers_ByEmails() {
	for _, userCount := range []int{0, 10, 100, 1000} {
		s.T().Run(fmt.Sprintf("%v", userCount), func(t *testing.T) {
			for i := 0; i < userCount; i++ {
				setup.MustCreateUserWithCleanup(t, context.Background(), s.API, setup.UserObj(fmt.Sprintf("%v-%v", userCount, i)))
			}
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				return context.Background, &request.GetUsersRequest{
					Emails: []user.Email{setup.UserObj(fmt.Sprintf("%v-%v", userCount, 0)).GetEmail()},
				}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.GetUsers(t, ctx, req.(*request.GetUsersRequest))
			}
			checkLen := func(ctx context.Context, resp any, err errwrap.Error) {
				require.Len(t, resp.(*response.GetUsersResponse).GetUsers(), min(userCount, 1))
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(userCount),
				DurationSec:     10,
				RampDurationSec: 10,
				Setup:           setup,
				Test:            tst,
				Assert:          checkLen,
				NFR: performance.NFRs{
					P50: 50 * time.Millisecond,
					P90: 90 * time.Millisecond,
					P95: 95 * time.Millisecond,
					P99: 99 * time.Millisecond,
				},
				Notes: []string{"test", "test", "test"},
			}
			report, _, _ := plan.Run()
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}

func (s *UserPerformanceSuite) TestPerformance_GetUsers_BySupervisor() {
	for _, userCount := range []int{0, 10, 100, 1000} {
		s.T().Run(fmt.Sprintf("%v", userCount), func(t *testing.T) {
			supervisor := setup.MustCreateUserWithCleanup(t, context.Background(), s.API, setup.UserObj("supervisor"))
			for i := 0; i < userCount; i++ {
				obj := setup.UserObj(fmt.Sprintf("%v-%v", userCount, i))
				obj.SetSupervisor(supervisor.GetEmail())
				setup.MustCreateUserWithCleanup(t, context.Background(), s.API, obj)
			}
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				return context.Background, &request.GetUsersRequest{
					Supervisor: supervisor.GetEmail(),
				}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.GetUsers(t, ctx, req.(*request.GetUsersRequest))
			}
			checkLen := func(ctx context.Context, resp any, err errwrap.Error) {
				require.Len(t, resp.(*response.GetUsersResponse).GetUsers(), userCount)
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(userCount),
				DurationSec:     10,
				RampDurationSec: 10,
				Setup:           setup,
				Test:            tst,
				Assert:          checkLen,
				NFR: performance.NFRs{
					P50: 50 * time.Millisecond,
					P90: 90 * time.Millisecond,
					P95: 95 * time.Millisecond,
					P99: 99 * time.Millisecond,
				},
				Notes: []string{"test", "test", "test"},
			}
			report, _, _ := plan.Run()
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}

func (s *UserPerformanceSuite) TestPerformance_PatchUser() {
	setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
		u := setup.MustCreateUserWithCleanup(s.T(), context.Background(), s.API, setup.UserObj(utils.RandAlphanum()))
		u.SetName(u.GetName() + "-patched")
		return context.Background, &request.PatchUserRequest{
			User: u,
		}, nil
	}
	tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
		if errReq != nil {
			return nil, errReq
		}
		return s.API.PatchUser(s.T(), ctx, req.(*request.PatchUserRequest))
	}
	plan := performance.Plan{
		T:               s.T(),
		RPS:             10,
		DurationSec:     10,
		RampDurationSec: 10,
		Setup:           setup,
		Test:            tst,
		NFR: performance.NFRs{
			P50: 50 * time.Millisecond,
			P90: 90 * time.Millisecond,
			P95: 95 * time.Millisecond,
			P99: 99 * time.Millisecond,
		},
		Notes: []string{"test", "test", "test"},
	}
	report, _, _ := plan.Run()
	fmt.Printf("%+v\n", report)
}

func (s *UserPerformanceSuite) TestPerformance_DeleteUser() {
	setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
		u := setup.MustCreateUserWithCleanup(s.T(), context.Background(), s.API, setup.UserObj(utils.RandAlphanum()))
		return context.Background, &request.DeleteUserRequest{
			ID: u.GetID(),
		}, nil
	}
	tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
		if errReq != nil {
			return nil, errReq
		}
		return s.API.DeleteUser(s.T(), ctx, req.(*request.DeleteUserRequest))
	}
	plan := performance.Plan{
		T:               s.T(),
		RPS:             10,
		DurationSec:     10,
		RampDurationSec: 10,
		Setup:           setup,
		Test:            tst,
		NFR: performance.NFRs{
			P50: 50 * time.Millisecond,
			P90: 90 * time.Millisecond,
			P95: 95 * time.Millisecond,
			P99: 99 * time.Millisecond,
		},
		Notes: []string{"test", "test", "test"},
	}
	report, _, _ := plan.Run()
	fmt.Printf("%+v\n", report)
}
