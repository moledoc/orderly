package main

import (
	"fmt"
	"html/template"
	"log"
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

var (
	strToState = map[string]State{
		"Not Started": NotStarted,
		"In Progress": InProgress,
		"Blocked":     Blocked,
		"Completed":   Completed,
	}
	stateToStr = map[State]string{
		NotStarted: "Not Started",
		InProgress: "In Progress",
		Blocked:    "Blocked",
		Completed:  "Completed",
	}
)

func (s State) String() string {
	val, ok := stateToStr[s]
	if !ok {
		return "Undefined"
	}
	return val
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

func validatePageData(pd PageData) error {
	zeroTime := time.Time{}
	if len(pd.MainTask) <= 0 {
		return fmt.Errorf("main-task is empty")
	}
	if len(pd.Responsible) <= 0 {
		return fmt.Errorf("responsible is empty")
	}
	if pd.State < 0 || StateCount <= pd.State {
		return fmt.Errorf("state is incorrect")
	}
	if pd.ParentOrder <= 0 {
		return fmt.Errorf("parent-order is empty")
	}
	if zeroTime.Equal(pd.Deadline) {
		return fmt.Errorf("deadline is empty")
	}
	return nil
}

func handleSave(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Printf("[ERROR]: parsing form: %v\n", err)
		return
	}

	formData := PageData{
		MainTask:    r.FormValue("order-info-main-task"),
		Responsible: r.FormValue("order-info-responsible"),
		State:       strToState[r.FormValue("order-info-state")],
	}
	parentOrder, _ := strconv.Atoi(r.FormValue("order-info-parent-order"))
	deadline, _ := time.Parse("2006-01-02", r.FormValue("order-info-deadline"))
	formData.ParentOrder = parentOrder
	formData.Deadline = deadline

	if err := validatePageData(formData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("[INFO]: validation failed: %s\n", err)
		w.Write([]byte(fmt.Sprintf(`{"error": "order-info-state is %s"}`, err)))
		return
	}

	changes := []string{}
	if data.MainTask != formData.MainTask {
		data.MainTask = formData.MainTask
		changes = append(changes, "main-task")
	}

	if data.Responsible != formData.Responsible {
		data.Responsible = formData.Responsible
		changes = append(changes, "responsible")
	}

	if data.ParentOrder != formData.ParentOrder {
		data.ParentOrder = formData.ParentOrder
		changes = append(changes, "parent-order")
	}

	if data.Deadline != formData.Deadline {
		data.Deadline = formData.Deadline
		changes = append(changes, "deadline")
	}

	if data.State != formData.State {
		data.State = formData.State
		changes = append(changes, "state")
	}

	if len(changes) > 0 {
		fmt.Printf("[INFO]: changes to: %v\n", changes)
	}
	w.WriteHeader(http.StatusOK)
}
