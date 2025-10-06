package order

import (
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
)

type State int

const (
	NotStarted State = iota
	InProgress
	HavingIssues
	Blocked
	Completed
)

type Task struct {
	ID          *meta.ID    `json:"id,omitempty"`
	State       *State      `json:"state,omitempty"`
	Accountable *user.Email `json:"accountable,omitempty"`
	Objective   *string     `json:"objective,omitempty"`
	Deadline    *time.Time  `json:"deadline,omitempty"`
}

type SitRep struct {
	ID            *meta.ID `json:"id,omitempty"`
	WorkCompleted *uint    `json:"work_completed,omitempty"`
	State         *State   `json:"state,omitempty"`
	Summary       *string  `json:"summary,omitempty"`
}

type Order struct {
	Task           *Task      `json:"task,omitempty"`
	DelegatedTasks []*Task    `json:"delegated_tasks,omitempty"`
	ParentOrderID  *meta.ID   `json:"parent_order_id,omitempty"`
	SitReps        []*SitRep  `json:"sitreps,omitempty"`
	Meta           *meta.Meta `json:"meta,omitempty"`
}

func Empty() *Order {
	return &Order{}
}

func (o *Order) Clone() *Order {
	if o == nil {
		return nil
	}
	var clone Order = Order{
		Task: &Task{
			ID:          utils.Ptr(o.GetTask().GetID()),
			State:       utils.Ptr(o.GetTask().GetState()),
			Accountable: utils.Ptr(o.GetTask().GetAccountable()),
			Objective:   utils.Ptr(o.GetTask().GetObjective()),
			Deadline:    utils.Ptr(o.GetTask().GetDeadline()),
		},
		ParentOrderID:  utils.Ptr(o.GetParentOrderID()),
		DelegatedTasks: make([]*Task, len(o.GetDelegatedTasks())),
		SitReps:        make([]*SitRep, len(o.GetSitReps())),
		Meta:           o.GetMeta().Clone(),
	}

	for i, delegatedTask := range o.GetDelegatedTasks() {
		clone.DelegatedTasks[i] = &Task{
			ID:          utils.Ptr(delegatedTask.GetID()),
			State:       utils.Ptr(delegatedTask.GetState()),
			Accountable: utils.Ptr(delegatedTask.GetAccountable()),
			Objective:   utils.Ptr(delegatedTask.GetObjective()),
		}
	}

	for i, sitrep := range o.GetSitReps() {
		clone.SitReps[i] = &SitRep{
			ID:            utils.Ptr(sitrep.GetID()),
			WorkCompleted: utils.Ptr(sitrep.GetWorkCompleted()),
			State:         utils.Ptr(sitrep.GetState()),
			Summary:       utils.Ptr(sitrep.GetSummary()),
		}
	}
	return &clone
}
