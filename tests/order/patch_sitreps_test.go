package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestPatchSitReps() {
	tt := s.T()

	expected := setup.MustCreateOrderWithCleanup(tt, context.Background(), s.API, setup.OrderObj())

	changes := []struct {
		name string
		f    func() *request.PatchSitRepsRequest
	}{
		{
			name: "sitrep.0.datetime",
			f: func() *request.PatchSitRepsRequest {
				if len(expected.GetSitReps()) == 0 {
					return nil
				}
				expected.GetSitReps()[0].SetDateTime(time.Now().UTC())
				return &request.PatchSitRepsRequest{
					OrderID: expected.GetID(),
					SitReps: []*order.SitRep{
						{
							ID:       expected.GetSitReps()[0].GetID(),
							DateTime: expected.GetSitReps()[0].GetDateTime(),
						},
					},
				}
			},
		},
		{
			name: "sitrep.0.by",
			f: func() *request.PatchSitRepsRequest {
				if len(expected.GetSitReps()) == 0 {
					return nil
				}
				expected.GetSitReps()[0].SetBy("example.by.updated@email.com")
				return &request.PatchSitRepsRequest{
					OrderID: expected.GetID(),
					SitReps: []*order.SitRep{
						{
							ID: expected.GetSitReps()[0].GetID(),
							By: expected.GetSitReps()[0].GetBy(),
						},
					},
				}
			},
		},
		{
			name: "sitrep.0.ping",
			f: func() *request.PatchSitRepsRequest {
				if len(expected.GetSitReps()) == 0 {
					return nil
				}
				expected.GetSitReps()[0].SetPing([]user.Email{"user1@email.com", "user2@email.com"})
				return &request.PatchSitRepsRequest{
					OrderID: expected.GetID(),
					SitReps: []*order.SitRep{
						{
							ID:   expected.GetSitReps()[0].GetID(),
							Ping: expected.GetSitReps()[0].GetPing(),
						},
					},
				}
			},
		},
		{
			name: "sitrep.0.situation",
			f: func() *request.PatchSitRepsRequest {
				if len(expected.GetSitReps()) == 0 {
					return nil
				}
				expected.GetSitReps()[0].SetSituation("updated situation")
				return &request.PatchSitRepsRequest{
					OrderID: expected.GetID(),
					SitReps: []*order.SitRep{
						{
							ID:        expected.GetSitReps()[0].GetID(),
							Situation: expected.GetSitReps()[0].GetSituation(),
						},
					},
				}
			},
		},
		{
			name: "sitrep.0.actions",
			f: func() *request.PatchSitRepsRequest {
				if len(expected.GetSitReps()) == 0 {
					return nil
				}
				expected.GetSitReps()[0].SetActions("updated actions")
				return &request.PatchSitRepsRequest{
					OrderID: expected.GetID(),
					SitReps: []*order.SitRep{
						{
							ID:      expected.GetSitReps()[0].GetID(),
							Actions: expected.GetSitReps()[0].GetActions(),
						},
					},
				}
			},
		},
		{
			name: "sitrep.0.tbd",
			f: func() *request.PatchSitRepsRequest {
				if len(expected.GetSitReps()) == 0 {
					return nil
				}
				expected.GetSitReps()[0].SetTBD("updated tbd")
				return &request.PatchSitRepsRequest{
					OrderID: expected.GetID(),
					SitReps: []*order.SitRep{
						{
							ID:  expected.GetSitReps()[0].GetID(),
							TBD: expected.GetSitReps()[0].GetTBD(),
						},
					},
				}
			},
		},
		{
			name: "sitrep.0.issues",
			f: func() *request.PatchSitRepsRequest {
				if len(expected.GetSitReps()) == 0 {
					return nil
				}
				expected.GetSitReps()[0].SetIssues("updated issues")
				return &request.PatchSitRepsRequest{
					OrderID: expected.GetID(),
					SitReps: []*order.SitRep{
						{
							ID:     expected.GetSitReps()[0].GetID(),
							Issues: expected.GetSitReps()[0].GetIssues(),
						},
					},
				}
			},
		},
		{
			name: "sitrep.1",
			f: func() *request.PatchSitRepsRequest {
				if len(expected.GetSitReps()) < 2 {
					return nil
				}

				expected.GetSitReps()[1].SetDateTime(time.Now().UTC())
				expected.GetSitReps()[1].SetBy("example.by.updated@email.com")
				expected.GetSitReps()[1].SetPing([]user.Email{"user1@email.com", "user2@email.com"})
				expected.GetSitReps()[1].SetSituation("updated situation")
				expected.GetSitReps()[1].SetActions("updated actions")
				expected.GetSitReps()[1].SetTBD("updated tbd")
				expected.GetSitReps()[1].SetIssues("updated issues")

				return &request.PatchSitRepsRequest{
					OrderID: expected.GetID(),
					SitReps: []*order.SitRep{
						{
							ID:        expected.GetSitReps()[1].GetID(),
							DateTime:  expected.GetSitReps()[1].GetDateTime(),
							By:        expected.GetSitReps()[1].GetBy(),
							Ping:      expected.GetSitReps()[1].GetPing(),
							Situation: expected.GetSitReps()[1].GetSituation(),
							Actions:   expected.GetSitReps()[1].GetActions(),
							TBD:       expected.GetSitReps()[1].GetTBD(),
							Issues:    expected.GetSitReps()[1].GetIssues(),
						},
					},
				}
			},
		},
	}

	for _, change := range changes {
		tt.Run(change.name, func(t *testing.T) {
			req := change.f()
			respPatch, err := s.API.PatchSitReps(t, context.Background(), req)
			require.NoError(t, err)

			expected.GetMeta().VersionIncr()
			expectedPatch := &response.PatchSitRepsResponse{
				Order: expected,
			}

			opts := []cmp.Option{
				compare.IgnorePaths("Order.Meta.Updated"),
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
