package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/response"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/cleanup"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/stretchr/testify/require"
)

func (s *OrderSuite) TestPostOrder_InputValidation() {
	tt := s.T()

	deadline, _ := time.Parse("2006-01-02", "2025-10-06")
	task := order.Task{
		State:       order.NotStarted,
		Accountable: user.Email("example@example.com"),
		Objective:   "this is the main objective",
		Deadline:    deadline,
	}
	parentOrderID := meta.NewID()
	delegatedTasks := []*order.Task{
		{
			State:       order.NotStarted,
			Accountable: user.Email("example@example.com"),
			Objective:   "this is delegated objective 1",
			Deadline:    deadline,
		},
		{
			State:       order.InProgress,
			Accountable: user.Email("example@example.com"),
			Objective:   "this is delegated objective 2",
			Deadline:    deadline,
		},
		{
			State:       order.Blocked,
			Accountable: user.Email("example@example.com"),
			Objective:   "this is delegated objective 3",
			Deadline:    deadline,
		},
		{
			State:       order.HavingIssues,
			Accountable: user.Email("example@example.com"),
			Objective:   "this is delegated objective 4",
			Deadline:    deadline,
		},
		{
			State:       order.Completed,
			Accountable: user.Email("example@example.com"),
			Objective:   "this is delegated objective 5",
			Deadline:    deadline,
		},
	}
	sitreps := []*order.SitRep{
		{
			WorkCompleted: 0,
			State:         order.NotStarted,
			Summary:       "summary of progress 1",
		},
		{
			WorkCompleted: 20,
			State:         order.InProgress,
			Summary:       "summary of progress 2",
		},
		{
			WorkCompleted: 80,
			State:         order.HavingIssues,
			Summary:       "summary of progress 3",
		},
		{
			WorkCompleted: 100,
			State:         order.Completed,
			Summary:       "summary of progress 4",
		},
	}

	tt.Run("EmptyRequest", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), nil)
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "empty request", err.GetStatusMessage())
		})
		t.Run("empty", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "empty order", err.GetStatusMessage())
		})
	})

	tt.Run("InvalidFieldProvided", func(t *testing.T) {
		t.Run("order.task.id", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &order.Task{
						ID:          meta.NewID(),
						State:       order.NotStarted,
						Accountable: user.Email("example@example.com"),
						Objective:   "this is the main objective",
						Deadline:    deadline,
					},
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "order.task.id disallowed", err.GetStatusMessage())
		})
		t.Run("order.meta", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
					Meta: &meta.Meta{
						Version: 2,
						Created: time.Now().UTC(),
						Updated: time.Now().UTC(),
					},
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "order.meta disallowed", err.GetStatusMessage())
		})

		/////////

		t.Run("order.delegatedTask.id", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &task,
					DelegatedTasks: []*order.Task{
						{
							ID:          meta.NewID(),
							State:       order.NotStarted,
							Accountable: user.Email("example@example.com"),
							Objective:   "this is the main objective",
							Deadline:    deadline,
						},
					},
					ParentOrderID: parentOrderID,
					SitReps:       sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "order.delegated.0.id disallowed", err.GetStatusMessage())
		})

		//////////

		t.Run("order.sitrep.id", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps: []*order.SitRep{
						{
							ID:            meta.NewID(),
							WorkCompleted: 50,
							State:         order.NotStarted,
							Summary:       "summary of progress 1",
						},
					},
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "order.sitrep.0.id disallowed", err.GetStatusMessage())
		})
	})

	tt.Run("MissingRequiredField", func(t *testing.T) {
		t.Run("order.task", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           nil,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "empty task", err.GetStatusMessage())
		})

		//////////

		t.Run("order.task.state", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &order.Task{
						Accountable: user.Email("example@example.com"),
						Objective:   "this is the main objective",
						Deadline:    deadline,
					},
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "task.state missing", err.GetStatusMessage())
		})
		t.Run("order.task.accountable", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &order.Task{
						State:     order.NotStarted,
						Objective: "this is the main objective",
						Deadline:  deadline,
					},
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "invalid email length", err.GetStatusMessage())
		})
		t.Run("order.task.objective", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &order.Task{
						State:       order.NotStarted,
						Accountable: user.Email("example@example.com"),
						Deadline:    deadline,
					},
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "invalid task.objective", err.GetStatusMessage())
		})
		t.Run("order.task.deadline", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &order.Task{
						State:       order.NotStarted,
						Accountable: user.Email("example@example.com"),
						Objective:   "this is the main objective",
					},
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "invalid task.deadline", err.GetStatusMessage())
		})

		//////////

		t.Run("order.parentOrderID", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
			require.Equal(t, "invalid id length", err.GetStatusMessage())
		})

		//////////

		t.Run("order.delegatedTasks.state", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &task,
					DelegatedTasks: []*order.Task{
						{
							Accountable: user.Email("example@example.com"),
							Objective:   "this is the main objective",
							Deadline:    deadline,
						},
					},
					ParentOrderID: parentOrderID,
					SitReps:       sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.delegatedTasks.accountable", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &task,
					DelegatedTasks: []*order.Task{
						{
							State:     order.NotStarted,
							Objective: "this is the main objective",
							Deadline:  deadline,
						},
					},
					ParentOrderID: parentOrderID,
					SitReps:       sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.delegatedTasks.objective", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &task,
					DelegatedTasks: []*order.Task{
						{
							State:       order.NotStarted,
							Accountable: user.Email("example@example.com"),
							Deadline:    deadline,
						},
					},
					ParentOrderID: parentOrderID,
					SitReps:       sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.delegatedTasks.deadline", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &task,
					DelegatedTasks: []*order.Task{
						{
							State:       order.NotStarted,
							Accountable: user.Email("example@example.com"),
							Objective:   "this is the main objective",
						},
					},
					ParentOrderID: parentOrderID,
					SitReps:       sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})

		////////////////

		t.Run("order.sitrep.workCompleted", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps: []*order.SitRep{
						{
							State:   order.NotStarted,
							Summary: "summary of progress 1",
						},
					},
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.sitrep.state", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps: []*order.SitRep{
						{
							WorkCompleted: 0,
							Summary:       "summary of progress 1",
						},
					},
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})

		t.Run("order.sitrep.summary", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps: []*order.SitRep{
						{
							WorkCompleted: 0,
							State:         order.NotStarted,
						},
					},
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})

	tt.Run("InvalidRequiredField", func(t *testing.T) {
		t.Run("order.task.state.low", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &order.Task{
						State:       order.NotStarted - 1,
						Accountable: user.Email("example@example.com"),
						Objective:   "this is the main objective",
						Deadline:    deadline,
					},
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.task.state.high", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &order.Task{
						State:       order.Completed + 1,
						Accountable: user.Email("example@example.com"),
						Objective:   "this is the main objective",
						Deadline:    deadline,
					},
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.task.accountable", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &order.Task{
						State:       order.Completed,
						Accountable: user.Email("incorrect email"),
						Objective:   "this is the main objective",
						Deadline:    deadline,
					},
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.task.deadline", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &order.Task{
						State:       order.Completed,
						Accountable: user.Email("example@example.com"),
						Objective:   "this is the main objective",
						Deadline:    time.Time{},
					},
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})

		//////////

		t.Run("order.delegatedTask.state.low", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &task,
					DelegatedTasks: []*order.Task{
						{
							State:       order.NotStarted - 1,
							Accountable: user.Email("example@example.com"),
							Objective:   "this is the main objective",
							Deadline:    deadline,
						},
					},
					ParentOrderID: parentOrderID,
					SitReps:       sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.delegatedTask.state.high", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &task,
					DelegatedTasks: []*order.Task{
						{
							State:       order.Completed + 1,
							Accountable: user.Email("example@example.com"),
							Objective:   "this is the main objective",
							Deadline:    deadline,
						},
					},
					ParentOrderID: parentOrderID,
					SitReps:       sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.delegatedTask.accountable", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &task,
					DelegatedTasks: []*order.Task{
						{
							State:       order.Completed,
							Accountable: user.Email("incorrect email"),
							Objective:   "this is the main objective",
							Deadline:    deadline,
						},
					},
					ParentOrderID: parentOrderID,
					SitReps:       sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.delegatedTask.deadline", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task: &task,
					DelegatedTasks: []*order.Task{
						{
							State:       order.Completed,
							Accountable: user.Email("example@example.com"),
							Objective:   "this is the main objective",
							Deadline:    time.Time{},
						},
					},
					ParentOrderID: parentOrderID,
					SitReps:       sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})

		//////////

		t.Run("order.parentOrderID.short", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID[:10],
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.parentOrderID.long", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID + parentOrderID,
					SitReps:        sitreps,
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})

		//////////

		t.Run("order.sitrep.workCompleted.high", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps: []*order.SitRep{
						{
							WorkCompleted: 150,
							State:         order.NotStarted,
							Summary:       "progress summary",
						},
					},
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.sitrep.state.low", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps: []*order.SitRep{
						{
							WorkCompleted: 50,
							State:         order.NotStarted - 1,
							Summary:       "progress summary",
						},
					},
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
		t.Run("order.sitrep.state.high", func(t *testing.T) {
			resp, err := s.API.PostOrder(t, context.Background(), &request.PostOrderRequest{
				Order: &order.Order{
					Task:           &task,
					DelegatedTasks: delegatedTasks,
					ParentOrderID:  parentOrderID,
					SitReps: []*order.SitRep{
						{
							WorkCompleted: 50,
							State:         order.Completed + 1,
							Summary:       "progress summary",
						},
					},
				},
			})
			defer cleanup.Order(t, s.API, resp.GetOrder())
			require.Error(t, err)
			require.Empty(t, resp)
			require.Equal(t, http.StatusBadRequest, err.GetStatusCode(), err)
		})
	})
}

func (s *OrderSuite) TestPostOrder() {
	tt := s.T()

	deadline, _ := time.Parse("2006-01-02", "2025-10-06")
	task := order.Task{
		State:       order.NotStarted,
		Accountable: user.Email("example@example.com"),
		Objective:   "this is the main objective",
		Deadline:    deadline,
	}
	parentOrderID := meta.NewID()
	delegatedTasks := []*order.Task{
		{
			State:       order.NotStarted,
			Accountable: user.Email("example@example.com"),
			Objective:   "this is delegated objective 1",
			Deadline:    deadline,
		},
		{
			State:       order.InProgress,
			Accountable: user.Email("example@example.com"),
			Objective:   "this is delegated objective 2",
			Deadline:    deadline,
		},
		{
			State:       order.Blocked,
			Accountable: user.Email("example@example.com"),
			Objective:   "this is delegated objective 3",
			Deadline:    deadline,
		},
		{
			State:       order.HavingIssues,
			Accountable: user.Email("example@example.com"),
			Objective:   "this is delegated objective 4",
			Deadline:    deadline,
		},
		{
			State:       order.Completed,
			Accountable: user.Email("example@example.com"),
			Objective:   "this is delegated objective 5",
			Deadline:    deadline,
		},
	}
	sitreps := []*order.SitRep{
		{
			WorkCompleted: 0,
			State:         order.NotStarted,
			Summary:       "summary of progress 1",
		},
		{
			WorkCompleted: 20,
			State:         order.InProgress,
			Summary:       "summary of progress 2",
		},
		{
			WorkCompleted: 80,
			State:         order.HavingIssues,
			Summary:       "summary of progress 3",
		},
		{
			WorkCompleted: 100,
			State:         order.Completed,
			Summary:       "summary of progress 4",
		},
	}
	o := &order.Order{
		Task:           &task,
		ParentOrderID:  parentOrderID,
		DelegatedTasks: delegatedTasks,
		SitReps:        sitreps,
	}

	resp, err := s.API.PostOrder(tt, context.Background(), &request.PostOrderRequest{
		Order: o,
	})
	defer cleanup.Order(tt, s.API, resp.GetOrder())
	require.NoError(tt, err)

	opts := []cmp.Option{
		compare.IgnoreID,
	}

	expected := &response.PostOrderResponse{
		Order: o,
	}

	compare.RequireEqual(tt, expected, resp, opts...)
}
