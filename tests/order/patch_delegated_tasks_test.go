package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestPatchDelegatedTasks() {
	tt := s.T()

	expected := setup.MustCreateOrderWithCleanup(tt, context.Background(), s.API, setup.OrderObj())

	changes := []struct {
		name string
		f    func() *request.PatchDelegatedTasksRequest
	}{
		{
			name: "task.0.state",
			f: func() *request.PatchDelegatedTasksRequest {
				if len(expected.GetDelegatedTasks()) == 0 {
					return nil
				}
				expected.GetDelegatedTasks()[0].SetState(order.HavingIssues)
				return &request.PatchDelegatedTasksRequest{
					OrderID: expected.GetID(),
					Tasks: []*order.Task{
						{
							ID:    expected.GetDelegatedTasks()[0].GetID(),
							State: utils.Ptr(expected.GetDelegatedTasks()[0].GetState()),
						},
					},
				}
			},
		},
		{
			name: "task.0.accountable",
			f: func() *request.PatchDelegatedTasksRequest {
				if len(expected.GetDelegatedTasks()) == 0 {
					return nil
				}
				expected.GetDelegatedTasks()[0].SetAccountable("example.accountable.updated@email.com")
				return &request.PatchDelegatedTasksRequest{
					OrderID: expected.GetID(),
					Tasks: []*order.Task{
						{
							ID:          expected.GetDelegatedTasks()[0].GetID(),
							Accountable: expected.GetDelegatedTasks()[0].GetAccountable(),
						},
					},
				}
			},
		},
		{
			name: "task.0.objective",
			f: func() *request.PatchDelegatedTasksRequest {
				if len(expected.GetDelegatedTasks()) == 0 {
					return nil
				}
				expected.GetDelegatedTasks()[0].SetObjective("updated objective")
				return &request.PatchDelegatedTasksRequest{
					OrderID: expected.GetID(),
					Tasks: []*order.Task{
						{
							ID:        expected.GetDelegatedTasks()[0].GetID(),
							Objective: expected.GetDelegatedTasks()[0].GetObjective(),
						},
					},
				}
			},
		},
		{
			name: "task.0.deadline",
			f: func() *request.PatchDelegatedTasksRequest {
				if len(expected.GetDelegatedTasks()) == 0 {
					return nil
				}
				expected.GetDelegatedTasks()[0].SetDeadline(time.Now().UTC())
				return &request.PatchDelegatedTasksRequest{
					OrderID: expected.GetID(),
					Tasks: []*order.Task{
						{
							ID:       expected.GetDelegatedTasks()[0].GetID(),
							Deadline: expected.GetDelegatedTasks()[0].GetDeadline(),
						},
					},
				}
			},
		},
		{
			name: "task.1",
			f: func() *request.PatchDelegatedTasksRequest {
				if len(expected.GetDelegatedTasks()) < 2 {
					return nil
				}
				expected.GetDelegatedTasks()[1].SetState(order.HavingIssues)
				expected.GetDelegatedTasks()[1].SetAccountable("example.accountable.updated@email.com")
				expected.GetDelegatedTasks()[1].SetObjective("updated objective")
				expected.GetDelegatedTasks()[1].SetDeadline(time.Now().UTC())
				return &request.PatchDelegatedTasksRequest{
					OrderID: expected.GetID(),
					Tasks: []*order.Task{
						{
							ID:          expected.GetDelegatedTasks()[1].GetID(),
							State:       utils.Ptr(expected.GetDelegatedTasks()[1].GetState()),
							Accountable: expected.GetDelegatedTasks()[1].GetAccountable(),
							Objective:   expected.GetDelegatedTasks()[1].GetObjective(),
							Deadline:    expected.GetDelegatedTasks()[1].GetDeadline(),
						},
					},
				}
			},
		},
	}

	for _, change := range changes {
		tt.Run(change.name, func(t *testing.T) {
			req := change.f()
			respPatch, err := s.API.PatchDelegatedTasks(t, context.Background(), req)
			require.NoError(t, err)

			expected.GetMeta().VersionIncr()
			expectedPatch := &response.PatchDelegatedTasksResponse{
				Order: expected,
			}

			opts := []cmp.Option{
				compare.IgnorePaths("Order.Meta.Updated"),
				compare.ComparerState(),
				compare.SorterOrder(compare.SortOrderByID),
				compare.SorterTask(compare.SortTaskByID),
				compare.SorterSitRep(compare.SortSitRepByID),
			}
			compare.RequireEqual(t, expectedPatch, respPatch, opts...)

			respGet, err := s.API.GetOrderByID(t, context.Background(), &request.GetOrderByIDRequest{
				ID: expected.GetID(),
			})
			require.NoError(t, err)

			expectedGet := &response.GetOrderByIDResponse{
				Order: expected,
			}
			compare.RequireEqual(t, expectedGet, respGet, opts...)
		})
	}
}
