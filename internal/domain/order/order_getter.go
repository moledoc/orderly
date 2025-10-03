package order

import (
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
)

func (tt *Task) GetID() meta.ID {
	if tt == nil || tt.ID == nil {
		return meta.EmptyID()
	}
	return *tt.ID
}

func (tt *Task) GetState() State {
	if tt == nil || tt.State == nil {
		return NotStarted
	}
	return *tt.State
}

func (tt *Task) GetAccountable() user.Email {
	if tt == nil || tt.Accountable == nil {
		return ""
	}
	return *tt.Accountable
}

func (tt *Task) GetObjective() string {
	if tt == nil || tt.Objective == nil {
		return ""
	}
	return *tt.Objective
}

func (tt *Task) GetMeta() *meta.Meta {
	if tt == nil {
		return nil
	}
	return tt.Meta
}

////////////

func (sr *SitRep) GetID() meta.ID {
	if sr == nil || sr.ID == nil {
		return meta.EmptyID()
	}
	return *sr.ID
}

func (sr *SitRep) GetState() State {
	if sr == nil || sr.State == nil {
		return NotStarted
	}
	return *sr.State
}

func (sr *SitRep) GetWorkCompleted() uint {
	if sr == nil || sr.WorkCompleted == nil {
		return 0
	}
	return *sr.WorkCompleted
}

func (sr *SitRep) GetSummary() string {
	if sr == nil || sr.Summary == nil {
		return ""
	}
	return *sr.Summary
}

func (sr *SitRep) GetMeta() *meta.Meta {
	if sr == nil {
		return nil
	}
	return sr.Meta
}

////////////

func (o *Order) GetTask() *Task {
	if o == nil {
		return nil
	}
	return o.Task
}

func (o *Order) GetSubTasks() []*Task {
	if o == nil {
		return nil
	}
	return o.SubTasks
}

func (o *Order) GetParentOrderID() meta.ID {
	if o == nil || o.ParentOrderID == nil {
		return meta.EmptyID()
	}
	return *o.ParentOrderID
}

func (o *Order) GetDeadline() time.Time {
	if o == nil || o.Deadline == nil {
		return time.Time{}
	}
	return *o.Deadline
}

func (o *Order) GetSitReps() []*SitRep {
	if o == nil || o.SitReps == nil {
		return nil
	}
	return o.SitReps
}
