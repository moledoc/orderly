package order

import (
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
)

func (tt *Task) SetID(id meta.ID) {
	if tt == nil {
		return
	}
	tt.ID = id
}

func (tt *Task) SetState(state State) {
	if tt == nil {
		return
	}
	tt.State = state
}

func (tt *Task) SetAccountable(accountable user.Email) {
	if tt == nil {
		return
	}
	tt.Accountable = accountable
}

func (tt *Task) SetObjective(objective string) {
	if tt == nil {
		return
	}
	tt.Objective = objective
}

func (tt *Task) SetDeadline(deadline time.Time) {
	if tt == nil {
		return
	}
	tt.Deadline = deadline
}

////////////

func (sr *SitRep) SetID(id meta.ID) {
	if sr == nil {
		return
	}
	sr.ID = id
}

func (sr *SitRep) SetState(state State) {
	if sr == nil {
		return
	}
	sr.State = state
}

func (sr *SitRep) SetWorkCompleted(workcompleted uint) {
	if sr == nil {
		return
	}
	sr.WorkCompleted = workcompleted
}

func (sr *SitRep) SetSummary(summary string) {
	if sr == nil {
		return
	}
	sr.Summary = summary
}

////////////

func (o *Order) SetTask(task *Task) {
	if o == nil {
		return
	}
	o.Task = task
}

func (o *Order) SetDelegatedTasks(delegatedTask []*Task) {
	if o == nil {
		return
	}
	o.DelegatedTasks = delegatedTask
}

func (o *Order) SetParentOrderID(parentorderid meta.ID) {
	if o == nil {
		return
	}
	o.ParentOrderID = parentorderid
}

func (o *Order) SetSitReps(sitreps []*SitRep) {
	if o == nil {
		return
	}
	o.SitReps = sitreps
}

func (o *Order) SetMeta(meta *meta.Meta) {
	if o == nil {
		return
	}
	o.Meta = meta
}
