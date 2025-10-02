package models

import (
	"time"

	"github.com/moledoc/orderly/utils"
)

type State int

const (
	NotStarted State = iota
	InProgress
	HavingIssues
	Blocked
	Completed
	StateCount
)

type Task struct {
	ID          *uint   `json:"id"`
	State       *State  `json:"state"`
	Accountable *string `json:"accountable"`
	Objective   *string `json:"objective"`
	Meta        *Meta   `json:"meta"`
}

type SitRep struct {
	ID            *uint   `json:"id"`
	Cron          *string `json:"cron"`
	WorkCompleted *uint   `json:"work_completed"`
	State         *State  `json:"state"`
	Summary       *string `json:"summary"`
	Meta          *Meta   `json:"meta"`
}

type Order struct {
	Task          *Task      `json:"task"`
	SubTasks      []*Task    `json:"subtasks"`
	ParentOrderID *uint      `json:"parent_order_id"`
	Deadline      *time.Time `json:"deadline"`
	SitReps       []*SitRep  `json:"sitreps"`
}

func (o *Order) Clone() *Order {
	var clone Order
	clone.Task = &Task{
		ID:          utils.RePtr(o.Task.ID),
		State:       utils.RePtr(o.Task.State),
		Accountable: utils.RePtr(o.Task.Accountable),
		Objective:   utils.RePtr(o.Task.Objective),
		Meta:        o.Task.Meta.Clone(),
	}
	clone.Deadline = utils.RePtr(o.Deadline)
	clone.ParentOrderID = utils.RePtr(o.ParentOrderID)
	clone.SubTasks = make([]*Task, len(o.SubTasks))
	for i, subtask := range o.SubTasks {
		clone.SubTasks[i] = &Task{
			ID:          utils.RePtr(subtask.ID),
			State:       utils.RePtr(subtask.State),
			Accountable: utils.RePtr(subtask.Accountable),
			Objective:   utils.RePtr(subtask.Objective),
			Meta:        subtask.Meta.Clone(),
		}
	}
	clone.SitReps = make([]*SitRep, len(o.SitReps))
	for i, sitrep := range o.SitReps {
		clone.SitReps[i] = &SitRep{
			ID:            utils.RePtr(sitrep.ID),
			Cron:          utils.RePtr(sitrep.Cron),
			WorkCompleted: utils.RePtr(sitrep.WorkCompleted),
			State:         utils.RePtr(sitrep.State),
			Summary:       utils.RePtr(sitrep.Summary),
			Meta:          sitrep.Meta.Clone(),
		}
	}
	return &clone
}
