package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/tests/performance"
	"github.com/moledoc/orderly/tests/setup"
)

func (s *OrderPerformanceSuite) TestPerformance_PostOrder() {
	setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
		return context.Background, &request.PostOrderRequest{
			Order: setup.OrderObj(),
		}, nil
	}
	tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
		if errReq != nil {
			return nil, errReq
		}
		return s.API.PostOrder(s.T(), ctx, req.(*request.PostOrderRequest))
	}
	plan := performance.Plan{
		T:               s.T(),
		RPS:             1,
		DurationSec:     1,
		RampDurationSec: 0,
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

func (s *OrderPerformanceSuite) TestPerformance_GetOrderByID() {
	setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
		o := setup.MustCreateOrderWithCleanup(s.T(), context.Background(), s.API, setup.OrderObj())
		return context.Background, &request.GetOrderByIDRequest{
			ID: o.GetID(),
		}, nil
	}
	tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
		if errReq != nil {
			return nil, errReq
		}
		return s.API.GetOrderByID(s.T(), ctx, req.(*request.GetOrderByIDRequest))
	}
	plan := performance.Plan{
		T:               s.T(),
		RPS:             1,
		DurationSec:     1,
		RampDurationSec: 0,
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

func (s *OrderPerformanceSuite) TestPerformance_GetOrders() {
	for _, orderCount := range []int{10, 100 /*, 1000*/} {
		s.T().Run(fmt.Sprintf("%v", orderCount), func(t *testing.T) {
			for i := 0; i < orderCount; i++ {
				setup.MustCreateOrderWithCleanup(s.T(), context.Background(), s.API, setup.OrderObj())
			}
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				return context.Background, &request.GetOrdersRequest{}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.GetOrders(t, ctx, req.(*request.GetOrdersRequest))
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(orderCount),
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
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}

func (s *OrderPerformanceSuite) TestPerformance_GetOrderSubOrders() {
	for _, orderCount := range []int{10, 100 /*, 1000*/} {
		s.T().Run(fmt.Sprintf("%v", orderCount), func(t *testing.T) {
			o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
			for i := 0; i < orderCount; i++ {
				orderObj := setup.OrderObj()
				orderObj.SetParentOrderID(o.GetID())
				setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, orderObj)
			}
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				return context.Background, &request.GetOrderSubOrdersRequest{
					ID: o.GetID(),
				}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.GetOrderSubOrders(t, ctx, req.(*request.GetOrderSubOrdersRequest))
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(orderCount),
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
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}

func (s *OrderPerformanceSuite) TestPerformance_PatchOrder() {
	setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
		o := setup.MustCreateOrderWithCleanup(s.T(), context.Background(), s.API, setup.OrderObj())
		oPatched := setup.OrderObj()
		oPatched.SetID(o.GetID())
		for i, delegated := range o.GetDelegatedTasks() {
			oPatched.DelegatedTasks[i].SetID(delegated.GetID())
		}
		for i, sitrep := range o.GetSitReps() {
			oPatched.SitReps[i].SetID(sitrep.GetID())
		}
		return context.Background, &request.PatchOrderRequest{
			Order: oPatched,
		}, nil
	}
	tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
		if errReq != nil {
			return nil, errReq
		}
		return s.API.PatchOrder(s.T(), ctx, req.(*request.PatchOrderRequest))
	}
	plan := performance.Plan{
		T:               s.T(),
		RPS:             1,
		DurationSec:     1,
		RampDurationSec: 0,
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

func (s *OrderPerformanceSuite) TestPerformance_DeleteOrder() {
	setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
		o := setup.MustCreateOrderWithCleanup(s.T(), context.Background(), s.API, setup.OrderObj())
		return context.Background, &request.DeleteOrderRequest{
			ID: o.GetID(),
		}, nil
	}
	tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
		if errReq != nil {
			return nil, errReq
		}
		return s.API.DeleteOrder(s.T(), ctx, req.(*request.DeleteOrderRequest))
	}
	plan := performance.Plan{
		T:               s.T(),
		RPS:             1,
		DurationSec:     1,
		RampDurationSec: 0,
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

func (s *OrderPerformanceSuite) TestPerformance_PutDelegatedTasks() {
	for _, taskCount := range []int{10, 100 /*, 1000*/} {
		s.T().Run(fmt.Sprintf("%v", taskCount), func(t *testing.T) {
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
				newTasks := make([]*order.Task, taskCount)
				for i := 0; i < taskCount; i++ {
					newTasks[i] = setup.TaskObj()
				}
				return context.Background, &request.PutDelegatedTasksRequest{
					OrderID: o.GetID(),
					Tasks:   newTasks,
				}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.PutDelegatedTasks(t, ctx, req.(*request.PutDelegatedTasksRequest))
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(taskCount),
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
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}

func (s *OrderPerformanceSuite) TestPerformance_PatchDelegatedTasks() {
	for _, taskCount := range []int{10, 100 /*, 1000*/} {
		s.T().Run(fmt.Sprintf("%v", taskCount), func(t *testing.T) {
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
				patchedOrder := setup.OrderObj()
				patchedDelegatedTasks := patchedOrder.GetDelegatedTasks()
				for i, delegated := range patchedDelegatedTasks {
					delegated.SetID(o.GetDelegatedTasks()[i].GetID())
				}
				return context.Background, &request.PatchDelegatedTasksRequest{
					OrderID: o.GetID(),
					Tasks:   patchedDelegatedTasks,
				}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.PatchDelegatedTasks(t, ctx, req.(*request.PatchDelegatedTasksRequest))
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(taskCount),
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
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}

func (s *OrderPerformanceSuite) TestPerformance_DeleteDelegatedTasks() {
	for _, taskCount := range []int{10, 100 /*, 1000*/} {
		s.T().Run(fmt.Sprintf("%v", taskCount), func(t *testing.T) {
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
				taskIDs := make([]meta.ID, len(o.GetDelegatedTasks()))
				for i, delegated := range o.GetDelegatedTasks() {
					taskIDs[i] = delegated.GetID()
				}
				return context.Background, &request.DeleteDelegatedTasksRequest{
					OrderID:          o.GetID(),
					DelegatedTaskIDs: taskIDs,
				}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.DeleteDelegatedTasks(t, ctx, req.(*request.DeleteDelegatedTasksRequest))
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(taskCount),
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
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}

func (s *OrderPerformanceSuite) TestPerformance_PutSitReps() {
	for _, taskCount := range []int{10, 100 /*, 1000*/} {
		s.T().Run(fmt.Sprintf("%v", taskCount), func(t *testing.T) {
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
				newSitReps := make([]*order.SitRep, taskCount)
				for i := 0; i < taskCount; i++ {
					newSitReps[i] = setup.SitrepObj()
				}
				return context.Background, &request.PutSitRepsRequest{
					OrderID: o.GetID(),
					SitReps: newSitReps,
				}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.PutSitReps(t, ctx, req.(*request.PutSitRepsRequest))
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(taskCount),
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
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}

func (s *OrderPerformanceSuite) TestPerformance_PatchSitReps() {
	for _, taskCount := range []int{10, 100 /*, 1000*/} {
		s.T().Run(fmt.Sprintf("%v", taskCount), func(t *testing.T) {
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
				patchedOrder := setup.OrderObj()
				patchedSitReps := patchedOrder.GetSitReps()
				for i, sitrep := range patchedSitReps {
					sitrep.SetID(o.GetSitReps()[i].GetID())
				}
				return context.Background, &request.PatchSitRepsRequest{
					OrderID: o.GetID(),
					SitReps: patchedSitReps,
				}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.PatchSitReps(t, ctx, req.(*request.PatchSitRepsRequest))
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(taskCount),
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
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}

func (s *OrderPerformanceSuite) TestPerformance_DeleteSitReps() {
	for _, taskCount := range []int{10, 100 /*, 1000*/} {
		s.T().Run(fmt.Sprintf("%v", taskCount), func(t *testing.T) {
			setup := func() (ctxFunc func() context.Context, req any, err errwrap.Error) {
				o := setup.MustCreateOrderWithCleanup(t, context.Background(), s.API, setup.OrderObj())
				sitrepIDs := make([]meta.ID, len(o.GetSitReps()))
				for i, sitrep := range o.GetSitReps() {
					sitrepIDs[i] = sitrep.GetID()
				}
				return context.Background, &request.DeleteSitRepsRequest{
					OrderID:   o.GetID(),
					SitRepIDs: sitrepIDs,
				}, nil
			}
			tst := func(ctx context.Context, req any, errReq errwrap.Error) (response any, errResp errwrap.Error) {
				if errReq != nil {
					return nil, errReq
				}
				return s.API.DeleteSitReps(t, ctx, req.(*request.DeleteSitRepsRequest))
			}
			plan := performance.Plan{
				T:               t,
				RPS:             uint(taskCount),
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
			fmt.Printf("%s: %+v\n", t.Name(), report)
		})
	}
}
