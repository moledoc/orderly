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

func (tt *Task) SetAccountable(accountable *user.User) {
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

func (sr *SitRep) SetDateTime(dateTime time.Time) {
	if sr == nil {
		return
	}
	sr.DateTime = dateTime
}

func (sr *SitRep) SetBy(by *user.User) {
	if sr == nil {
		return
	}
	sr.By = by
}

func (sr *SitRep) SetPing(ping []*user.User) {
	if sr == nil {
		return
	}
	sr.Ping = ping
}

func (sr *SitRep) SetSituation(situation string) {
	if sr == nil {
		return
	}
	sr.Situation = situation
}

func (sr *SitRep) SetActions(actions string) {
	if sr == nil {
		return
	}
	sr.Actions = actions
}

func (sr *SitRep) SetTBD(tBD string) {
	if sr == nil {
		return
	}
	sr.TBD = tBD
}

func (sr *SitRep) SetIssues(issues string) {
	if sr == nil {
		return
	}
	sr.Issues = issues
}

////////////

func (o *Order) SetID(id meta.ID) {
	if o == nil {
		return
	}
	o.GetTask().SetID(id)
}

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
