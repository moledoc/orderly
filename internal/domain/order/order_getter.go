package order

import (
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
)

func (tt *Task) GetID() meta.ID {
	if tt == nil {
		return meta.EmptyID()
	}
	return tt.ID
}

func (tt *Task) GetState() State {
	if tt == nil {
		return NotStarted
	}
	return tt.State
}

func (tt *Task) GetAccountable() user.Email {
	if tt == nil {
		return ""
	}
	return tt.Accountable
}

func (tt *Task) GetObjective() string {
	if tt == nil {
		return ""
	}
	return tt.Objective
}

func (tt *Task) GetDeadline() time.Time {
	if tt == nil {
		return time.Time{}
	}
	return tt.Deadline
}

////////////

func (sr *SitRep) GetID() meta.ID {
	if sr == nil {
		return meta.EmptyID()
	}
	return sr.ID
}

func (sr *SitRep) GetDateTime() time.Time {
	if sr == nil {
		return time.Time{}
	}
	return sr.DateTime
}

func (sr *SitRep) GetBy() user.Email {
	if sr == nil {
		return ""
	}
	return sr.By
}

func (sr *SitRep) GetPing() []user.Email {
	if sr == nil {
		return []user.Email{}
	}
	return sr.Ping
}

func (sr *SitRep) GetSituation() string {
	if sr == nil {
		return ""
	}
	return sr.Situation
}

func (sr *SitRep) GetActions() string {
	if sr == nil {
		return ""
	}
	return sr.Actions
}

func (sr *SitRep) GetTBD() string {
	if sr == nil {
		return ""
	}
	return sr.TBD
}

func (sr *SitRep) GetIssues() string {
	if sr == nil {
		return ""
	}
	return sr.Issues
}

////////////

func (o *Order) GetID() meta.ID {
	if o == nil {
		return ""
	}
	return o.GetTask().GetID()
}

func (o *Order) GetTask() *Task {
	if o == nil {
		return nil
	}
	return o.Task
}

func (o *Order) GetDelegatedTasks() []*Task {
	if o == nil {
		return nil
	}
	return o.DelegatedTasks
}

func (o *Order) GetParentOrderID() meta.ID {
	if o == nil {
		return meta.EmptyID()
	}
	return o.ParentOrderID
}

func (o *Order) GetSitReps() []*SitRep {
	if o == nil {
		return nil
	}
	return o.SitReps
}

func (o *Order) GetMeta() *meta.Meta {
	if o == nil {
		return nil
	}
	return o.Meta
}
