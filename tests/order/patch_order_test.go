package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestPatchOrder() {
	tt := s.T()

	expected := setup.MustCreateOrderWithCleanup(tt, context.Background(), s.API, setup.OrderObj())

	changes := []struct {
		name string
		f    func() *request.PatchOrderRequest
	}{
		{
			name: "parentOrgId",
			f: func() *request.PatchOrderRequest {
				newID := meta.NewID()
				expected.SetParentOrderID(newID)

				return &request.PatchOrderRequest{
					Order: &order.Order{
						Task:          &order.Task{ID: expected.GetID()},
						ParentOrderID: newID,
					},
				}
			},
		},
		{
			name: "task",
			f: func() *request.PatchOrderRequest {
				patchedTask := utils.RePtr(expected.GetTask())

				patchedTask.SetState(order.Blocked)
				patchedTask.SetAccountable("changed@email.com")
				patchedTask.SetObjective("patched main objective")
				patchedTask.SetDeadline(time.Now().UTC())

				expected.SetTask(patchedTask)

				return &request.PatchOrderRequest{
					Order: &order.Order{
						Task: patchedTask,
					},
				}
			},
		},
		{
			name: "delegated_task",
			f: func() *request.PatchOrderRequest {
				if len(expected.GetDelegatedTasks()) == 0 {
					return nil
				}
				patchedDelegatedTask := &order.Task{
					ID:          expected.GetDelegatedTasks()[0].GetID(),
					State:       order.InProgress,
					Accountable: user.Email("new.accountable@email.com"),
					Objective:   "new objective",
					Deadline:    time.Now().UTC().Add(30 * 24 * time.Hour),
				}
				expected.GetDelegatedTasks()[0].SetState(patchedDelegatedTask.GetState())
				expected.GetDelegatedTasks()[0].SetAccountable(patchedDelegatedTask.GetAccountable())
				expected.GetDelegatedTasks()[0].SetObjective(patchedDelegatedTask.GetObjective())
				expected.GetDelegatedTasks()[0].SetDeadline(patchedDelegatedTask.GetDeadline())

				return &request.PatchOrderRequest{
					Order: &order.Order{
						Task:           &order.Task{ID: expected.GetID()},
						DelegatedTasks: []*order.Task{patchedDelegatedTask},
					},
				}
			},
		},
		{
			name: "sitrep",
			f: func() *request.PatchOrderRequest {
				if len(expected.GetSitReps()) == 0 {
					return nil
				}
				patchedSitRep := &order.SitRep{
					ID: expected.GetSitReps()[0].GetID(),

					DateTime: time.Now().UTC(),
					By:       user.Email("user1@email.com"),
					Ping: []user.Email{
						user.Email("user2@email.com"),
						user.Email("user3@email.com"),
					},
					// Situation: "New situation description",
					// Actions:   "<List of actions>",
					// TBD:       "<List of things to do>",
					Issues: "<List of issues>",
				}

				expected.GetSitReps()[0].SetDateTime(patchedSitRep.GetDateTime())
				expected.GetSitReps()[0].SetBy(patchedSitRep.GetBy())
				expected.GetSitReps()[0].SetPing(patchedSitRep.GetPing())
				// expected.GetSitReps()[0].SetSituation(patchedSitRep.GetSituation())
				// expected.GetSitReps()[0].SetActions(patchedSitRep.GetActions())
				// expected.GetSitReps()[0].SetTBD(patchedSitRep.GetTBD())
				expected.GetSitReps()[0].SetIssues(patchedSitRep.GetIssues())

				return &request.PatchOrderRequest{
					Order: &order.Order{
						Task:    &order.Task{ID: expected.GetID()},
						SitReps: []*order.SitRep{patchedSitRep},
					},
				}
			},
		},
		{
			name: "partial",
			f: func() *request.PatchOrderRequest {
				if len(expected.GetDelegatedTasks()) == 0 || len(expected.GetSitReps()) == 0 {
					return nil
				}
				patchedDelegatedTask := &order.Task{
					ID:        expected.GetDelegatedTasks()[0].GetID(),
					Objective: "partial new objective",
				}
				patchedSitRep := &order.SitRep{
					ID:        expected.GetSitReps()[0].GetID(),
					Situation: "partial New situation description",
				}
				expected.GetDelegatedTasks()[0].SetObjective(patchedDelegatedTask.GetObjective())
				expected.GetSitReps()[0].SetSituation(patchedSitRep.GetSituation())

				return &request.PatchOrderRequest{
					Order: &order.Order{
						Task:           &order.Task{ID: expected.GetID()},
						DelegatedTasks: []*order.Task{patchedDelegatedTask},
						SitReps:        []*order.SitRep{patchedSitRep},
					},
				}
			},
		},
	}

	for _, change := range changes {
		tt.Run(change.name, func(t *testing.T) {
			req := change.f()
			respPatch, err := s.API.PatchOrder(t, context.Background(), req)
			require.NoError(t, err)

			expectedPatch := &response.PatchOrderResponse{
				Order: expected,
			}

			compare.RequireEqual(t, expectedPatch, respPatch)
			require.NotEmpty(t, respPatch.GetOrder().GetMeta())

			respGet, err := s.API.GetOrderByID(t, context.Background(), &request.GetOrderByIDRequest{
				ID: expected.GetID(),
			})
			require.NoError(t, err)

			expectedGet := &response.GetOrderByIDResponse{
				Order: expected,
			}

			opts := []cmp.Option{
				compare.SorterOrder(compare.SortOrderByID),
				compare.SorterTask(compare.SortTaskByID),
				compare.SorterSitRep(compare.SortSitRepByID),
			}
			compare.RequireEqual(t, expectedGet, respGet, opts...)
			require.NotEmpty(t, respGet.GetOrder().GetMeta())
		})
	}
}
