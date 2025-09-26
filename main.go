package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
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

type Info struct {
	OrderNr     int
	MainTask    string
	Responsible string
	ParentOrder int
	Deadline    time.Time
	State       State
	States      []State
}

type Task struct {
	Id          int
	Checked     bool
	Responsible string
	Description string
}

type PageData struct {
	Info  *Info
	Tasks *[]*Task
}

func main() {
	http.HandleFunc("/", servePage)
	http.HandleFunc("POST /save-info", handleSaveInfo)
	http.HandleFunc("POST /save-tasks", handleSaveTasks)
	http.HandleFunc("POST /add-task", handleAddTask)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func formatToDate(t time.Time) string {
	return t.Format("2006-01-02")
}

var (
	info = &Info{
		OrderNr:     2,
		MainTask:    "Go html templates + htmx",
		Responsible: "Lala",
		ParentOrder: 1,
		Deadline:    time.Date(2025, 12, 12, 0, 0, 0, 0, time.UTC),
		State:       InProgress,
		States:      StateValues(),
	}

	tasks = &[]*Task{
		{
			Id:          1,
			Checked:     false,
			Responsible: "lala",
			Description: "eat pancake",
		},
		{
			Id:          2,
			Checked:     true,
			Responsible: "lala",
			Description: "eat sandwich",
		},
	}

	pageData = PageData{
		Info:  info,
		Tasks: tasks,
	}
)

func servePage(w http.ResponseWriter, r *http.Request) {

	tmplFilename := "order.tmpl.html"
	tmpl, err := template.New(tmplFilename).Funcs(template.FuncMap{
		"formatToDate": formatToDate,
	}).ParseFiles(tmplFilename)

	if err != nil {
		fmt.Printf("[ERROR]: parse template: %s\n", err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Printf("[ERROR]: executing html tmpl failed: %s\n", err)
	}
}

func validatePageData(pd Info) error {
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

func handleSaveInfo(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Printf("[ERROR]: parsing form: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "parsing form failed: %s"}`, err)))
		return
	}

	formData := Info{
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
	if info.MainTask != formData.MainTask {
		info.MainTask = formData.MainTask
		changes = append(changes, "main-task")
	}

	if info.Responsible != formData.Responsible {
		info.Responsible = formData.Responsible
		changes = append(changes, "responsible")
	}

	if info.ParentOrder != formData.ParentOrder {
		info.ParentOrder = formData.ParentOrder
		changes = append(changes, "parent-order")
	}

	if info.Deadline != formData.Deadline {
		info.Deadline = formData.Deadline
		changes = append(changes, "deadline")
	}

	if info.State != formData.State {
		info.State = formData.State
		changes = append(changes, "state")
	}

	if len(changes) > 0 {
		fmt.Printf("[INFO]: changes to: %v\n", changes)
	}
	w.WriteHeader(http.StatusOK)
}

func handleSaveTasks(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Printf("[ERROR]: parsing form: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "parsing form failed: %s"}`, err)))
		return
	}

	fmt.Printf("HERE1: %+v\n", r.Form)

	collected := make(map[string]*Task)
	for k, v := range r.Form {
		elems := strings.Split(k, "-")
		i := elems[len(elems)-1]
		_, ok := collected[i]
		if !ok {
			collected[i] = &Task{}
		}
		ii, err := strconv.Atoi(i)
		if err != nil {
			log.Printf("[WARNING]: idx str to int conversion failed: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "parsing idx failed: %s"}`, err)))
			return
		}
		collected[i].Id = ii
		kk := strings.Replace(k, fmt.Sprintf("-%v", ii), "", 1)
		switch kk {
		case "task-checkbox":
			collected[i].Checked = v[0] == "on"
		case "task-responsible":
			collected[i].Responsible = v[0]
		case "task-description":
			collected[i].Description = v[0]
		}
	}
	*tasks = []*Task{}
	for _, task := range collected {
		*tasks = append(*tasks, task)
	}
	slices.SortFunc(*tasks, func(a *Task, b *Task) int {
		if a.Id < b.Id {
			return -1
		} else if a.Id > b.Id {
			return 1
		} else {
			return 0
		}
	})
	w.WriteHeader(http.StatusOK)
}

func handleAddTask(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Printf("[ERROR]: parsing form: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "parsing form failed: %s"}`, err)))
		return
	}

	fmt.Printf("HERE: %+v\n", r.Form)
	if len(r.Form) == 0 {
		return
	}

	l := len(*tasks) - 1
	if l < 0 {
		l = 0
	}
	task := (*tasks)[l]
	idx := task.Id + 1
	newTask := &Task{
		Id:          idx,
		Checked:     false,
		Responsible: r.FormValue("new-task-responsible"),
		Description: r.FormValue("new-task-description"),
	}
	*tasks = append(*tasks, newTask)

	w.WriteHeader(http.StatusOK)
}
