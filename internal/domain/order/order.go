package order

import (
	"encoding/json"
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
)

type State int

const (
	delegatedStr  string = "Delegated"
	receivedStr   string = "Received"
	planningStr   string = "Planning"
	executingStr  string = "Executing"
	problemsStr   string = "Problems"
	cancelleddStr string = "Cancelledd"
	doneStr       string = "Done"
)

const (
	Delegated State = iota
	Received
	Planning
	Executing
	Problems
	Cancelled
	Done
)

func (s State) String() string {

	switch s {
	case Delegated:
		return delegatedStr
	case Received:
		return receivedStr
	case Planning:
		return planningStr
	case Executing:
		return executingStr
	case Problems:
		return problemsStr
	case Cancelled:
		return cancelleddStr
	case Done:
		return doneStr
	default:
		return delegatedStr
	}
}

func (s *State) UnmarshalJSON(data []byte) error {

	var nr int
	if err := json.Unmarshal(data, &nr); err == nil && 0 <= State(nr) && State(nr) <= Done {
		*s = State(nr)
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {

	case delegatedStr:
		*s = Delegated
	case receivedStr:
		*s = Received
	case planningStr:
		*s = Planning
	case executingStr:
		*s = Executing
	case problemsStr:
		*s = Problems
	case cancelleddStr:
		*s = Cancelled
	case doneStr:
		*s = Done
	default:
		*s = 0
	}

	return nil
}

func (s *State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

var (
	ListStates func() []*State = func() func() []*State {
		states := []*State{}
		for i := State(0); i <= Done; i++ {
			states = append(states, &i)
		}
		return func() []*State {
			return states
		}
	}()
)

type SitRep struct {
	ID meta.ID `json:"id,omitempty"`

	DateTime time.Time  `json:"datetime,omitempty"`
	By       user.Email `json:"by,omitempty"`

	Situation string `json:"situation,omitempty"`
	Actions   string `json:"actions,omitempty"`
	TODO      string `json:"todo,omitempty"`
	Issues    string `json:"issues,omitempty"`
}

type Order struct {
	ID            meta.ID   `json:"id,omitempty"`
	AccountableID meta.ID   `json:"accountable_id,omitempty"`
	ParentOrderID meta.ID   `json:"parent_order_id,omitempty"`
	Objective     string    `json:"objective,omitempty"`
	Deadline      time.Time `json:"deadline,omitempty"`
	State         *State    `json:"state,omitempty"`

	DelegatedOrders []*Order  `json:"delegated_orders,omitempty"`
	SitReps         []*SitRep `json:"sitreps,omitempty"`

	Meta *meta.Meta `json:"meta,omitempty"`
}

func Empty() *Order {
	return &Order{}
}

func (o *Order) Clone() *Order {
	if o == nil {
		return nil
	}
	var clone Order = Order{
		ID:              o.GetID(),
		State:           utils.Ptr(o.GetState()),
		AccountableID:   o.GetAccountableID(),
		ParentOrderID:   o.GetParentOrderID(),
		Objective:       o.GetObjective(),
		Deadline:        o.GetDeadline(),
		DelegatedOrders: make([]*Order, len(o.GetDelegatedOrders())),
		SitReps:         make([]*SitRep, len(o.GetSitReps())),
		Meta:            o.GetMeta().Clone(),
	}

	for i, delegatedOrder := range o.GetDelegatedOrders() {
		clone.DelegatedOrders[i] = &Order{
			ID:            delegatedOrder.GetID(),
			State:         utils.Ptr(delegatedOrder.GetState()),
			AccountableID: delegatedOrder.GetAccountableID(),
			Objective:     delegatedOrder.GetObjective(),
			Deadline:      delegatedOrder.GetDeadline(),
		}
	}

	for i, sitrep := range o.GetSitReps() {
		clone.SitReps[i] = &SitRep{
			ID: sitrep.GetID(),

			DateTime: sitrep.GetDateTime(),
			By:       sitrep.GetBy(),

			Situation: sitrep.GetSituation(),
			Actions:   sitrep.GetActions(),
			TODO:      sitrep.GetTODO(),
			Issues:    sitrep.GetIssues(),
		}
	}
	return &clone
}
