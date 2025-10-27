package order

import (
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
)

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

func (sr *SitRep) GetTODO() string {
	if sr == nil {
		return ""
	}
	return sr.TODO
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
		return meta.EmptyID()
	}
	return o.ID
}

func (o *Order) GetAccountableID() meta.ID {
	if o == nil {
		return ""
	}
	return o.AccountableID
}

func (o *Order) GetParentOrderID() meta.ID {
	if o == nil {
		return meta.EmptyID()
	}
	return o.ParentOrderID
}

func (o *Order) GetObjective() string {
	if o == nil {
		return ""
	}
	return o.Objective
}

func (o *Order) GetState() State {
	if o == nil || o.State == nil {
		return NotStarted
	}
	return *(o.State)
}

func (o *Order) GetDeadline() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.Deadline
}

func (o *Order) GetDelegatedOrders() []*Order {
	if o == nil {
		return nil
	}
	return o.DelegatedOrders
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
