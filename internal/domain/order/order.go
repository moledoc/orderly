package order

import (
	"encoding/json"
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
)

type State int

func (s State) String() string {

	switch s {
	case NotStarted:
		return "Not Started"
	case InProgress:
		return "In Progress"
	case HavingIssues:
		return "Having Issues"
	case Blocked:
		return "Blocked"
	case Completed:
		return "Completed"
	default:
		return "Unknown"
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
	case "Not Started":
		*s = NotStarted
	case "In Progress":
		*s = InProgress
	case "Having Issues":
		*s = HavingIssues
	case "Blocked":
		*s = Blocked
	case "Completed":
		*s = Completed
	default:
		*s = NotStarted
	}

	return nil
}

func (s *State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

const (
	NotStarted State = iota + 1
	InProgress
	HavingIssues
	Blocked
	Completed
)

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

type Task struct {
	ID          meta.ID    `json:"id,omitempty"`
	State       *State     `json:"state,omitempty"`
	Accountable *user.User `json:"accountable,omitempty"`
	Objective   string     `json:"objective,omitempty"`
	Deadline    time.Time  `json:"deadline,omitempty"`
}

type SitRep struct {
	ID meta.ID `json:"id,omitempty"`

	DateTime time.Time  `json:"datetime,omitempty"`
	By       *user.User `json:"email,omitempty"`

	Situation string `json:"situation,omitempty"`
	Actions   string `json:"actions,omitempty"`
	TODO      string `json:"todo,omitempty"`
	Issues    string `json:"issues,omitempty"`
}

type Order struct {
	Task           *Task      `json:"task,omitempty"`
	DelegatedTasks []*Task    `json:"delegated_tasks,omitempty"`
	ParentOrderID  meta.ID    `json:"parent_order_id,omitempty"`
	SitReps        []*SitRep  `json:"sitreps,omitempty"`
	Meta           *meta.Meta `json:"meta,omitempty"`
}

func Empty() *Order {
	return &Order{}
}

func (o *Order) Clone() *Order {
	if o == nil {
		return nil
	}
	var clone Order = Order{
		Task: &Task{
			ID:          o.GetTask().GetID(),
			State:       utils.Ptr(o.GetTask().GetState()),
			Accountable: o.GetTask().GetAccountable(),
			Objective:   o.GetTask().GetObjective(),
			Deadline:    o.GetTask().GetDeadline(),
		},
		ParentOrderID:  o.GetParentOrderID(),
		DelegatedTasks: make([]*Task, len(o.GetDelegatedTasks())),
		SitReps:        make([]*SitRep, len(o.GetSitReps())),
		Meta:           o.GetMeta().Clone(),
	}

	for i, delegatedTask := range o.GetDelegatedTasks() {
		clone.DelegatedTasks[i] = &Task{
			ID:          delegatedTask.GetID(),
			State:       utils.Ptr(delegatedTask.GetState()),
			Accountable: delegatedTask.GetAccountable(),
			Objective:   delegatedTask.GetObjective(),
			Deadline:    delegatedTask.GetDeadline(),
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
