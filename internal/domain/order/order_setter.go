package order

import (
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
)

func (tt *Task) SetID(id meta.ID) {
	if tt == nil || tt.ID == nil {
		return
	}
	*tt.ID = id
}

func (tt *Task) SetState(state State) {
	if tt == nil || tt.State == nil {
		return
	}
	*tt.State = state
}

func (tt *Task) SetAccountable(accountable user.Email) {
	if tt == nil || tt.Accountable == nil {
		return
	}
	*tt.Accountable = accountable
}

func (tt *Task) SetObjective(objective string) {
	if tt == nil || tt.Objective == nil {
		return
	}
	*tt.Objective = objective
}

func (tt *Task) SetMeta(meta *meta.Meta) {
	if tt == nil {
		return
	}
	tt.Meta = meta
}

////////////

func (sr *SitRep) SetID(id meta.ID) {
	if sr == nil || sr.ID == nil {
		return
	}
	*sr.ID = id
}

func (sr *SitRep) SetState(state State) {
	if sr == nil || sr.State == nil {
		return
	}
	*sr.State = state
}

func (sr *SitRep) SetWorkCompleted(workcompleted uint) {
	if sr == nil || sr.WorkCompleted == nil {
		return
	}
	*sr.WorkCompleted = workcompleted
}

func (sr *SitRep) SetSummary(summary string) {
	if sr == nil || sr.Summary == nil {
		return
	}
	*sr.Summary = summary
}

func (sr *SitRep) SetMeta(meta *meta.Meta) {
	if sr == nil {
		return
	}
	sr.Meta = meta
}

////////////

func (o *Order) SetTask(task *Task) {
	if o == nil {
		return
	}
	o.Task = task
}

func (o *Order) SetSubTasks(subtasks []*Task) {
	if o == nil {
		return
	}
	o.SubTasks = subtasks
}

func (o *Order) SetParentOrderID(parentorderid meta.ID) {
	if o == nil || o.ParentOrderID == nil {
		return
	}
	*o.ParentOrderID = parentorderid
}

func (o *Order) SetDeadline(deadline time.Time) {
	if o == nil || o.Deadline == nil {
		return
	}
	*o.Deadline = deadline
}

func (o *Order) SetSitReps(sitreps []*SitRep) {
	if o == nil || o.SitReps == nil {
		return
	}
	o.SitReps = sitreps
}
