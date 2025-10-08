package order

import (
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
)

type State int

const (
	NotStarted State = iota + 1
	InProgress
	HavingIssues
	Blocked
	Completed
)

type Task struct {
	ID          meta.ID    `json:"id,omitempty"`
	State       State      `json:"state,omitempty"`
	Accountable user.Email `json:"accountable,omitempty"`
	Objective   string     `json:"objective,omitempty"`
	Deadline    time.Time  `json:"deadline,omitempty"`
}

type SitRep struct {
	ID            meta.ID `json:"id,omitempty"`
	WorkCompleted uint    `json:"work_completed,omitempty"`
	State         State   `json:"state,omitempty"`
	Summary       string  `json:"summary,omitempty"`
}

type Order struct {
	Task           *Task      `json:"task,omitempty"`
	DelegatedTasks []*Task    `json:"delegated_tasks,omitempty"`
	ParentOrderID  meta.ID    `json:"parent_order_id,omitempty"`
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
			ID:          o.GetTask().GetID(),
			State:       o.GetTask().GetState(),
			Accountable: o.GetTask().GetAccountable(),
			Objective:   o.GetTask().GetObjective(),
			Deadline:    o.GetTask().GetDeadline(),
		},
		ParentOrderID:  o.GetParentOrderID(),
		DelegatedTasks: make([]*Task, len(o.GetDelegatedTasks())),
		SitReps:        make([]*SitRep, len(o.GetSitReps())),
		Meta:           o.GetMeta().Clone(),
	}

	for i, delegatedTask := range o.GetDelegatedTasks() {
		clone.DelegatedTasks[i] = &Task{
			ID:          delegatedTask.GetID(),
			State:       delegatedTask.GetState(),
			Accountable: delegatedTask.GetAccountable(),
			Objective:   delegatedTask.GetObjective(),
		}
	}

	for i, sitrep := range o.GetSitReps() {
		clone.SitReps[i] = &SitRep{
			ID:            sitrep.GetID(),
			WorkCompleted: sitrep.GetWorkCompleted(),
			State:         sitrep.GetState(),
			Summary:       sitrep.GetSummary(),
		}
	}
	return &clone
}
