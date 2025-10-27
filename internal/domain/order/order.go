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
	notStartedStr   string = "Not Started"
	inProgressStr   string = "In Progress"
	havingIssuesStr string = "Having Issues"
	blockedStr      string = "Blocked"
	cancelledStr    string = "Cancelled"
	completedStr    string = "Completed"
	unknownStr      string = "Unknown"
)

const (
	NotStarted State = iota + 1
	InProgress
	HavingIssues
	Blocked
	Cancelled
	Completed
)

func (s State) String() string {

	switch s {
	case NotStarted:
		return notStartedStr
	case InProgress:
		return inProgressStr
	case HavingIssues:
		return havingIssuesStr
	case Blocked:
		return blockedStr
	case Cancelled:
		return cancelledStr
	case Completed:
		return completedStr
	default:
		return unknownStr
	}
}

func (s *State) UnmarshalJSON(data []byte) error {

	var nr int
	if err := json.Unmarshal(data, &nr); err == nil && NotStarted <= State(nr) && State(nr) <= Completed {
		*s = State(nr)
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	case notStartedStr:
		*s = NotStarted
	case inProgressStr:
		*s = InProgress
	case havingIssuesStr:
		*s = HavingIssues
	case blockedStr:
		*s = Blocked
	case cancelledStr:
		*s = Cancelled
	case completedStr:
		*s = Completed
	default:
		*s = NotStarted
	}

	return nil
}

func (s *State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

var (
	ListStates func() []*State = func() func() []*State {
		states := []*State{}
		for i := NotStarted; i <= Completed; i++ {
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
	By       user.Email `json:"email,omitempty"`

	Situation string `json:"situation,omitempty"`
	Actions   string `json:"actions,omitempty"`
	TODO      string `json:"todo,omitempty"`
	Issues    string `json:"issues,omitempty"`
}

type Order struct {
	ID            meta.ID    `json:"id,omitempty"`
	Accountable   user.Email `json:"accountable,omitempty"`
	ParentOrderID meta.ID    `json:"parent_order_id,omitempty"`
	Objective     string     `json:"objective,omitempty"`
	Deadline      time.Time  `json:"deadline,omitempty"`
	State         *State     `json:"state,omitempty"`

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
		Accountable:     o.GetAccountable(),
		ParentOrderID:   o.GetParentOrderID(),
		Objective:       o.GetObjective(),
		Deadline:        o.GetDeadline(),
		DelegatedOrders: make([]*Order, len(o.GetDelegatedOrders())),
		SitReps:         make([]*SitRep, len(o.GetSitReps())),
		Meta:            o.GetMeta().Clone(),
	}

	for i, delegatedOrder := range o.GetDelegatedOrders() {
		clone.DelegatedOrders[i] = &Order{
			ID:          delegatedOrder.GetID(),
			State:       utils.Ptr(delegatedOrder.GetState()),
			Accountable: delegatedOrder.GetAccountable(),
			Objective:   delegatedOrder.GetObjective(),
			Deadline:    delegatedOrder.GetDeadline(),
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
