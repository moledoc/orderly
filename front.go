package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type State int

const (
	NotStarted State = iota
	InProgress
	Blocked
	Completed
	StateCount
)

func (s State) String() string {
	switch s {
	case NotStarted:
		return "Not Started"
	case InProgress:
		return "In Progress"
	case Blocked:
		return "Blocked"
	case Completed:
		return "Completed"
	default:
		return "Undefined"
	}
}

func StateValues() []State {
	states := make([]State, StateCount)
	for i := 0; i < int(StateCount); i++ {
		states[i] = State(i)
	}
	return states
}

type PageData struct {
	OrderNr     int
	MainTask    string
	Responsible string
	ParentOrder int
	Deadline    time.Time
	State       State
	States      []State
}

func main() {
	http.HandleFunc("/", servePage)
	http.HandleFunc("POST /save", handleSave)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func formatToDate(t time.Time) string {
	return t.Format("2006-01-02")
}

var (
	data = PageData{
		OrderNr:     2,
		MainTask:    "Go html templates + htmx",
		Responsible: "Lala",
		ParentOrder: 1,
		Deadline:    time.Date(2025, 12, 12, 0, 0, 0, 0, time.UTC),
		State:       InProgress,
		States:      StateValues(),
	}
)

func servePage(w http.ResponseWriter, r *http.Request) {

	tmplFilename := "order.tmpl.html"
	tmpl, err := template.New(tmplFilename).Funcs(template.FuncMap{
		"formatToDate": formatToDate,
	}).ParseFiles(tmplFilename)

	if err != nil {
		fmt.Printf("[ERROR]: parse template: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

func handleSave(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Printf("[ERROR]: parsing form: %v\n", err)
		return
	}

	changes := []string{}
	mainTask := r.FormValue("order-info-main-task")
	if data.MainTask != mainTask && len(mainTask) > 0 {
		data.MainTask = mainTask
		changes = append(changes, "main-task")
	}

	responsible := r.FormValue("order-info-responsible")
	if data.Responsible != responsible && len(responsible) > 0 {
		data.Responsible = responsible
		changes = append(changes, "responsible")
	}

	parentOrder, err := strconv.Atoi(r.FormValue("order-info-parent-order"))
	if err == nil && data.ParentOrder != parentOrder && parentOrder > 0 {
		data.ParentOrder = parentOrder
		changes = append(changes, "parent-order")
	}

	deadline, err := time.Parse("2006-01-02", r.FormValue("order-info-deadline"))
	if err == nil && data.Deadline != deadline {
		data.Deadline = deadline
		changes = append(changes, "deadline")
	}

	state, err := strconv.Atoi(r.FormValue("order-info-state"))
	if err == nil && data.State != State(state) {
		data.State = State(state)
		changes = append(changes, "state")
	}

	if len(changes) > 0 {
		fmt.Printf("[INFO]: changes to: %v\n", changes)
	}
}
