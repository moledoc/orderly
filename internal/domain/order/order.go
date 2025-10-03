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
	Meta        *meta.Meta  `json:"meta,omitempty"`
}

type SitRep struct {
	ID            *meta.ID   `json:"id,omitempty"`
	WorkCompleted *uint      `json:"work_completed,omitempty"`
	State         *State     `json:"state,omitempty"`
	Summary       *string    `json:"summary,omitempty"`
	Meta          *meta.Meta `json:"meta,omitempty"`
}

type Order struct {
	Task          *Task      `json:"task,omitempty"`
	SubTasks      []*Task    `json:"subtasks,omitempty"`
	ParentOrderID *meta.ID   `json:"parent_order_id,omitempty"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	SitReps       []*SitRep  `json:"sitreps,omitempty"`
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
			Meta:        o.Task.Meta.Clone(),
		},
		Deadline:      utils.Ptr(o.GetDeadline()),
		ParentOrderID: utils.Ptr(o.GetParentOrderID()),
		SubTasks:      make([]*Task, len(o.GetSubTasks())),
		SitReps:       make([]*SitRep, len(o.GetSitReps())),
	}

	for i, subtask := range o.GetSubTasks() {
		clone.SubTasks[i] = &Task{
			ID:          utils.Ptr(subtask.GetID()),
			State:       utils.Ptr(subtask.GetState()),
			Accountable: utils.Ptr(subtask.GetAccountable()),
			Objective:   utils.Ptr(subtask.GetObjective()),
			Meta:        subtask.GetMeta().Clone(),
		}
	}

	for i, sitrep := range o.GetSitReps() {
		clone.SitReps[i] = &SitRep{
			ID:            utils.Ptr(sitrep.GetID()),
			WorkCompleted: utils.Ptr(sitrep.GetWorkCompleted()),
			State:         utils.Ptr(sitrep.GetState()),
			Summary:       utils.Ptr(sitrep.GetSummary()),
			Meta:          sitrep.GetMeta().Clone(),
		}
	}
	return &clone
}
