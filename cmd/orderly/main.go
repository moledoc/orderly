package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/moledoc/orderly/tests/setup"
)

var (
	htmlBase = `<!DOCTYPE html>
	<html lang="en">
	<head>
	</head>
	
	<body>
	%v
	</body>
	</html>`
)

var (
	orderTempl = `
ID: {{.Task.ID}}<br>
<div>
<label for="parent-order-id"><a href="/order/{{.ParentOrderID}}">Parent Order ID:</a></label>
<input type="text" id="parent-order-id" style="border: none;" name="order-task" value="{{.ParentOrderID}}" required>
</div>

<div>
	<label for="task-objective">Objective:</label>
	<input type="text" id="task-objective" style="border: none;" name="order-task" value="{{.Task.Objective}}" required>
</div>

<div>
	<label for="task-accountable">Accountable:</label>
	<input list="users" id="task-accountable" name="order-task" value="{{.Task.Accountable.Email}}"/>
	<datalist id="users">
	{{range UserEmails}}
		<option value="{{.}}">{{.}}</option>
    {{end}}
	</datalist>
</div>

<div>
	<label for="task-deadline">Deadline:</label>
	<input type="date" id="task-deadline" style="border: none;"  name="order-task" value="{{formatToDate .Task.Deadline}}" min="1970-01-01" max="3000-01-01" required />
</div>

<div>
	<label for="task-state">State:</label>
	<select id="task-state" name="order-task" style="border: none;" required>
	{{range States}}
        <option value="{{.}}" {{if eq . $.Task.State}}selected{{end}}>{{.}}</option>
    {{end}}
	</select>
</div>

<div>
<details>
  <summary>Delegated Tasks</summary>
</caption>
 <table style="width: 100%; border-collapse: collapse;">
  <thead>
    <tr>
      <th>ID</th>
      <th>Objective</th>
      <th>Accountable</th>
      <th>Deadline</th>
    </tr>
  </thead>
  <tbody>
	{{range .DelegatedTasks}}
		<tr>
			<td style="width: 25%; border: 1px solid #ccc;"><a href="/order/{{.ID}}">{{.ID}}</a></td>
			<td style="width: 25%; border: 1px solid #ccc;">{{.Objective}}</td>
			<td style="width: 25%; border: 1px solid #ccc;"><a href="/user/{{.Accountable.ID}}">{{.Accountable.Email}}</a></td>
			<td style="width: 25%; border: 1px solid #ccc;"><input type="date" style="border: none;" value="{{formatToDate .Deadline}}" min="1970-01-01" max="3000-01-01" required /></td>
		</tr>
	{{end}}
  </tbody>
</table>
</details>
</div>

<details>
  <summary>SitReps</summary>
{{range .SitReps}}
<div class="card-body" id="sitreps" style="border: 1px solid #ccc; padding: 16px; margin-bottom: 10px; border-radius: 4px;">
ID: {{.ID}}<br>
By: <a href="/user/{{.By.ID}}">{{.By.Email}}</a><br>
DateTime: {{formatToDate .DateTime}}<br>
Ping: {{range .Ping}}
<a href="/user/{{.ID}}">{{.Email}}</a>
{{end}}<br>

<label for="sitrep-situation">Situation:</label>
<input type="text" id="sitrep-situation" style="border: none;" name="order-sitrep" value="{{.Situation}}"><br>

<label for="sitrep-actions">Actions:</label>
<input type="text" id="sitrep-actions" style="border: none;" name="order-sitrep" value="{{.Actions}}"><br>


<label for="sitrep-tbd">TBD:</label>
<input type="text" id="sitrep-tbd" style="border: none;" name="order-sitrep" value="{{.TBD}}"><br>


<label for="sitrep-issues">Issues:</label>
<input type="text" id="sitrep-issues" style="border: none;" name="order-sitrep" value="{{.Issues}}"><br>
</div>
{{end}}
</details>
`

	htmlOrder = fmt.Sprintf(htmlBase, orderTempl)
)

var (
	ordersTempl = `
<div>
 <table style="width: 100%; border-collapse: collapse;">
  <thead>
    <tr>
      <th>ID</th>
      <th>Objective</th>
      <th>Deadline</th>
      <th>Accountable</th>
    </tr>
  </thead>
  <tbody>
	{{range .}}
		<tr>
			<td style="width: 25%; border: 1px solid #ccc;"><a href="/order/{{.Task.ID}}">{{.Task.ID}}</a></td>
			<td style="width: 25%; border: 1px solid #ccc;">{{.Task.Objective}}</td>
			<td style="width: 25%; border: 1px solid #ccc;"><input type="date" style="border: none;" value="{{formatToDate .Task.Deadline}}" min="1970-01-01" max="3000-01-01" required /></td>
			<td style="width: 25%; border: 1px solid #ccc;"><a href="/user/{{.Task.Accountable.ID}}">{{.Task.Accountable.Email}}</a></td>
		</tr>
	{{end}}
  </tbody>
</table>
</div>`

	htmlOrders = fmt.Sprintf(htmlBase, ordersTempl)
)

var (
	usersTempl = `
<div>
 <table style="width: 100%; border-collapse: collapse;">
  <thead>
    <tr>
      <th>ID</th>
      <th>Name</th>
      <th>Email</th>
    </tr>
  </thead>
  <tbody>
	{{range .}}
		<tr>
			<td style="width: 25%; border: 1px solid #ccc;"><a href="/user/{{.ID}}">{{.ID}}</a></td>
			<td style="width: 25%; border: 1px solid #ccc;">{{.Name}}</td>
			<td style="width: 25%; border: 1px solid #ccc;">{{.Email}}</td>
		</tr>
	{{end}}
  </tbody>
</table>
</div>`

	htmlUsers = fmt.Sprintf(htmlBase, usersTempl)
)

var (
	userTempl = `
ID: {{.User.ID}}<br>

<div>
	<label for="user-name">Name:</label>
	<input type="text" id="user-name" style="border: none;" name="user" value="{{.User.Name}}" required>
</div>

Email: {{.User.Email}}<br>

<div>
	<label for="user-supervisor">Supervisor:</label>
	<input type="text" id="user-supervisor" style="border: none;" name="user" value="{{.User.Supervisor}}" required>
</div>

<div>
<details>
  <summary>Accountable Orders</summary>
</caption>
 <table style="width: 100%; border-collapse: collapse;">
  <thead>
    <tr>
      <th>ID</th>
      <th>Objective</th>
      <th>Deadline</th>
      <th>State</th>
    </tr>
  </thead>
  <tbody>
	{{range .AccountableOrders}}
		<tr>
			<td style="width: 25%; border: 1px solid #ccc;"><a href="/order/{{.Task.ID}}">{{.Task.ID}}</a></td>
			<td style="width: 25%; border: 1px solid #ccc;">{{.Task.Objective}}</td>
			<td style="width: 25%; border: 1px solid #ccc;"><input type="date" style="border: none;" value="{{formatToDate .Task.Deadline}}" min="1970-01-01" max="3000-01-01" required /></td>
			<td style="width: 25%; border: 1px solid #ccc;">{{.Task.State}}</td>
		</tr>
	{{end}}
  </tbody>
</table>
</details>
</div>

<div>
<details>
  <summary>Sub-Ordinates</summary>
</caption>
 <table style="width: 100%; border-collapse: collapse;">
  <thead>
    <tr>
      <th>ID</th>
      <th>Name</th>
      <th>Email</th>
    </tr>
  </thead>
  <tbody>
	{{range .SubOrdinates}}
		<tr>
			<td style="width: 25%; border: 1px solid #ccc;"><a href="/user/{{.ID}}">{{.ID}}</a></td>
			<td style="width: 25%; border: 1px solid #ccc;">{{.Name}}</td>
			<td style="width: 25%; border: 1px solid #ccc;">{{.Email}}</td>
		</tr>
	{{end}}
  </tbody>
</table>
</details>
</div>
`

	htmlUser = fmt.Sprintf(htmlBase, userTempl)
)

func getStates() []order.State {
	states := []order.State{}
	for i := order.NotStarted; i <= order.Completed; i++ {
		states = append(states, i)
	}
	return states
}

func getUsers() []*user.User {
	// TODO: get all users
	userCount := 5
	us := make([]*user.User, userCount)
	for i := 0; i < userCount; i++ {
		us[i] = setup.UserObjWithID(fmt.Sprintf("%v", i))
	}
	return us
}

func getUserEmails() []string {
	us := getUsers()

	emails := make([]string, len(us))
	for i, u := range us {
		emails[i] = string(u.GetEmail())
	}
	return emails
}

func formatToDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func serveOrders(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("page").Funcs(template.FuncMap{
		"formatToDate": formatToDate,
		"States":       getStates,
		"UserEmails":   getUserEmails,
	}).Parse(htmlOrders))

	// TODO: get orders

	// REMOVEME: START: when getting order by ID
	orders := []*order.Order{}
	for i := 0; i < 10; i++ {
		orders = append(orders, setup.OrderObjWithIDs(fmt.Sprintf("%v", i)))
	}
	// REMOVEME: END: when getting order by ID

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(w, orders)
	if err != nil {
		log.Printf("[ERROR]: executing html tmpl failed: %s\n", err)
	}
}

func serveOrder(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("page").Funcs(template.FuncMap{
		"formatToDate": formatToDate,
		"States":       getStates,
		"UserEmails":   getUserEmails,
	}).Parse(htmlOrder))

	// TODO: get order by ID

	// REMOVEME: START: when getting order by ID
	order := setup.OrderObjWithIDs("1")
	// REMOVEME: END: when getting order by ID

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(w, order)
	if err != nil {
		log.Printf("[ERROR]: executing html tmpl failed: %s\n", err)
	}
}

func serveUsers(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("page").Funcs(template.FuncMap{
		"formatToDate": formatToDate,
	}).Parse(htmlUsers))

	users := getUsers()

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(w, users)
	if err != nil {
		log.Printf("[ERROR]: executing html tmpl failed: %s\n", err)
	}
}

func serveUser(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("page").Funcs(template.FuncMap{
		"formatToDate": formatToDate,
	}).Parse(htmlUser))

	// TODO: get user by ID
	// REMOVEME: START: when getting user by ID
	u := setup.UserObjWithID()
	// REMOVEME: END: when getting user by ID

	// TODO: get user subordinates
	// REMOVEME: START: when getting user subordinates
	subOrdinatesCount := 5
	subOrdinates := make([]*user.User, subOrdinatesCount)
	for i := 0; i < subOrdinatesCount; i++ {
		subOrdinates[i] = setup.UserObjWithID(fmt.Sprintf("%v", i))
	}
	// REMOVEME: END: when getting user subordinates

	// TODO: get user accountable orders
	// REMOVEME: START: when getting user accountable orders
	accountableOrdersCount := 5
	accountableOrders := make([]*order.Order, accountableOrdersCount)
	for i := 0; i < accountableOrdersCount; i++ {
		accountableOrders[i] = setup.OrderObjWithIDs(fmt.Sprintf("%v", i))
	}
	// REMOVEME: END: when getting user accountable orders

	type userExtended struct {
		*user.User
		SubOrdinates      []*user.User
		AccountableOrders []*order.Order
	}

	ue := &userExtended{
		User:              u,
		SubOrdinates:      subOrdinates,
		AccountableOrders: accountableOrders,
	}

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(w, ue)
	if err != nil {
		log.Printf("[ERROR]: executing html tmpl failed: %s\n", err)
	}
}

func main() {

	router.Route(&router.Service{
		MgmtOrder: mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()),
		MgmtUser:  mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()),
	})

	// http.HandleFunc("GET /", serveLogin) // NOTE: login
	// http.HandleFunc("GET /", serveLogin) // NOTE: login

	http.HandleFunc("GET /orders", serveOrders)
	http.HandleFunc("GET /order/{id}", serveOrder)

	http.HandleFunc("GET /users", serveUsers)
	http.HandleFunc("GET /user/{id}", serveUser)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
