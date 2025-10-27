package order

import (
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
)

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

func (sr *SitRep) SetBy(by user.Email) {
	if sr == nil {
		return
	}
	sr.By = by
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

func (sr *SitRep) SetTODO(todo string) {
	if sr == nil {
		return
	}
	sr.TODO = todo
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
	o.ID = id
}

func (o *Order) SetAccountableID(id meta.ID) {
	if o == nil {
		return
	}
	o.AccountableID = id
}

func (o *Order) SetParentOrderID(id meta.ID) {
	if o == nil {
		return
	}
	o.ParentOrderID = id
}

func (o *Order) SetObjective(objective string) {
	if o == nil {
		return
	}
	o.Objective = objective
}

func (o *Order) SetState(state State) {
	if o == nil || o.State == nil {
		return
	}
	*(o.State) = state
}

func (o *Order) SetDeadline(time time.Time) {
	if o == nil {
		return
	}
	o.Deadline = time
}

func (o *Order) SetDelegatedOrders(orders []*Order) {
	if o == nil {
		return
	}
	o.DelegatedOrders = orders
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
