package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/setup"
)

var (
	htmlBase = `<!DOCTYPE html>
	<html lang="en">
	
<head>
	<meta charset="UTF-8">
	<title>Task Details</title>
	<style>
		body {
			font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
			background-color: #f4f6f8;
			color: #333;
			padding: 2rem;
			line-height: 1.6;
		}
		.container {
			margin: 0 auto;
			background: #fff;
			padding: 2rem;
			border-radius: 8px;
			box-shadow: 0 2px 8px rgba(0,0,0,0.05);
			overflow-x: auto;
			max-width: 100%%;
		}

		h1 {
			font-size: 1.5rem;
			margin-bottom: 1rem;
		}
		.form-group {
			margin-bottom: 1.25rem;
		}
		label {
			display: block;
			font-weight: 600;
			margin-bottom: 0.5rem;
		}
		input[type="text"],
		input[type="date"],
		select {
			width: 100%%;
			padding: 0.5rem;
			border: 1px solid #ccc;
			border-radius: 4px;
			font-size: 1rem;
		}
		input:disabled {
			background-color: #f9f9f9;
		}
		table {
			width: 100%%;
			border-collapse: collapse;
			margin-top: 1rem;
		}
		th, td {
			padding: 0.75rem;
			border: 1px solid #ddd;
			text-align: left;
		}
		details summary {
			cursor: pointer;
			font-weight: bold;
			margin-top: 1rem;
			margin-bottom: 0.5rem;
		}
		.card {
			background-color: #f9f9f9;
			padding: 1rem;
			border: 1px solid #ccc;
			border-radius: 6px;
			margin-bottom: 1rem;
		}
		a {
			color: #007BFF;
			text-decoration: none;
		}
		a:hover {
			text-decoration: underline;
		}

		table {
			width: 100%%;
			border-collapse: collapse;
			margin-top: 1rem;
			background-color: #fff;
			border-radius: 6px;
			overflow: hidden;
		}

		th, td {
			padding: 0.75rem 1rem;
			border: 1px solid #ddd;
			text-align: left;
		}

		th {
			background-color: #f0f0f0;
			font-weight: 600;
		}

		input[type="date"] {
			border: none;
			background-color: transparent;
			font-size: 1rem;
		}
	</style>
</head>
	
	<body>
	<div class="container">
	%v
	</div>
	</body>
	</html>`
)

var (
	orderTempl0 = `
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
 <table style="width: 100%%; border-collapse: collapse;">
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
			<td style="width: 25%%; border: 1px solid #ccc;"><a href="/order/{{.ID}}">{{.ID}}</a></td>
			<td style="width: 25%%; border: 1px solid #ccc;">{{.Objective}}</td>
			<td style="width: 25%%; border: 1px solid #ccc;"><a href="/user/{{.Accountable.ID}}">{{.Accountable.Email}}</a></td>
			<td style="width: 25%%; border: 1px solid #ccc;"><input type="date" style="border: none;" value="{{formatToDate .Deadline}}" min="1970-01-01" max="3000-01-01" required /></td>
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

	orderTempl = `
	<h1>Task ID: {{.Task.ID}}</h1>

	<div class="form-group">
		<label for="parent-order-id">Parent Order ID:</label>
		<a href="/order/{{.ParentOrderID}}">{{.ParentOrderID}}</a>
		<input type="text" id="parent-order-id" name="order-task" value="{{.ParentOrderID}}" disabled>
	</div>

	<div class="form-group">
		<label for="task-objective">Objective:</label>
		<input type="text" id="task-objective" name="order-task" value="{{.Task.Objective}}" required>
	</div>

	<div class="form-group">
		<label for="task-accountable">Accountable:</label>
		<input list="users" id="task-accountable" name="order-task" value="{{.Task.Accountable.Email}}">
		<datalist id="users">
			{{range UserEmails}}
			<option value="{{.}}">{{.}}</option>
			{{end}}
		</datalist>
	</div>

	<div class="form-group">
		<label for="task-deadline">Deadline:</label>
		<input type="date" id="task-deadline" name="order-task" value="{{formatToDate .Task.Deadline}}" required>
	</div>

	<div class="form-group">
		<label for="task-state">State:</label>
		<select id="task-state" name="order-task" required>
			{{range States}}
			<option value="{{.}}" {{if eq . $.Task.State}}selected{{end}}>{{.}}</option>
			{{end}}
		</select>
	</div>

	<details>
		<summary>Delegated Tasks</summary>
		<table>
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
					<td><a href="/order/{{.ID}}">{{.ID}}</a></td>
					<td>{{.Objective}}</td>
					<td><a href="/user/{{.Accountable.ID}}">{{.Accountable.Email}}</a></td>
					<td><input type="date" value="{{formatToDate .Deadline}}" required disabled></td>
				</tr>
				{{end}}
			</tbody>
		</table>
	</details>

	<details>
		<summary>SitReps</summary>
		{{range .SitReps}}
		<div class="card">
			<div><strong>ID:</strong> {{.ID}}</div>
			<div><strong>By:</strong> <a href="/user/{{.By.ID}}">{{.By.Email}}</a></div>
			<div><strong>DateTime:</strong> {{formatToDate .DateTime}}</div>
			<div><strong>Ping:</strong>
				{{range $index, $ping := .Ping}}
					<a href="/user/{{.ID}}">{{.Email}}</a>
				{{end}}
			</div>

			<div class="form-group">
				<label for="sitrep-situation-{{.ID}}">Situation:</label>
				<input type="text" id="sitrep-situation-{{.ID}}" name="order-sitrep" value="{{.Situation}}">
			</div>

			<div class="form-group">
				<label for="sitrep-actions-{{.ID}}">Actions:</label>
				<input type="text" id="sitrep-actions-{{.ID}}" name="order-sitrep" value="{{.Actions}}">
			</div>

			<div class="form-group">
				<label for="sitrep-tbd-{{.ID}}">TBD:</label>
				<input type="text" id="sitrep-tbd-{{.ID}}" name="order-sitrep" value="{{.TBD}}">
			</div>

			<div class="form-group">
				<label for="sitrep-issues-{{.ID}}">Issues:</label>
				<input type="text" id="sitrep-issues-{{.ID}}" name="order-sitrep" value="{{.Issues}}">
			</div>
		</div>
		{{end}}
	</details>
`
	htmlOrder = fmt.Sprintf(htmlBase, orderTempl)
)

var (
	ordersTempl = `
<div>
 <table style="width: 100%%; border-collapse: collapse;">
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
			<td style="width: 25%%; border: 1px solid #ccc;"><a href="/order/{{.Task.ID}}">{{.Task.ID}}</a></td>
			<td style="width: 25%%; border: 1px solid #ccc;">{{.Task.Objective}}</td>
			<td style="width: 25%%; border: 1px solid #ccc;"><input type="date" style="border: none;" value="{{formatToDate .Task.Deadline}}" min="1970-01-01" max="3000-01-01" required /></td>
			<td style="width: 25%%; border: 1px solid #ccc;"><a href="/user/{{.Task.Accountable.ID}}">{{.Task.Accountable.Email}}</a></td>
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
 <table style="width: 100%%; border-collapse: collapse;">
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
			<td style="width: 25%%; border: 1px solid #ccc;"><a href="/user/{{.ID}}">{{.ID}}</a></td>
			<td style="width: 25%%; border: 1px solid #ccc;">{{.Name}}</td>
			<td style="width: 25%%; border: 1px solid #ccc;">{{.Email}}</td>
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
 <table style="width: 100%%; border-collapse: collapse;">
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
			<td style="width: 25%%; border: 1px solid #ccc;"><a href="/order/{{.Task.ID}}">{{.Task.ID}}</a></td>
			<td style="width: 25%%; border: 1px solid #ccc;">{{.Task.Objective}}</td>
			<td style="width: 25%%; border: 1px solid #ccc;"><input type="date" style="border: none;" value="{{formatToDate .Task.Deadline}}" min="1970-01-01" max="3000-01-01" required /></td>
			<td style="width: 25%%; border: 1px solid #ccc;">{{.Task.State}}</td>
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
 <table style="width: 100%%; border-collapse: collapse;">
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
			<td style="width: 25%%; border: 1px solid #ccc;"><a href="/user/{{.ID}}">{{.ID}}</a></td>
			<td style="width: 25%%; border: 1px solid #ccc;">{{.Name}}</td>
			<td style="width: 25%%; border: 1px solid #ccc;">{{.Email}}</td>
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

func getOrders() []*order.Order {
	// TODO: get all orders
	orderCount := 1000
	os := make([]*order.Order, orderCount)
	for i := 0; i < orderCount; i++ {
		os[i] = setup.OrderObjWithIDs(fmt.Sprintf("%v", i))
		os[i].Task.Objective += utils.RandAlphanum() + utils.RandAlphanum() + utils.RandAlphanum()
	}
	return os
}

func getParentOrder(orderID meta.ID) *order.Order {
	// TODO: get parent order by id
	if orderID == "" {
		return nil
	}
	o := setup.OrderObjWithIDs("parent")
	o.Task.Objective += "\nnewlined parent order objective"
	return o
}

func getSubOrdinates(userID meta.ID) []*user.User {
	// TODO: get user sub-ordinates
	userCount := 50
	us := make([]*user.User, userCount)
	for i := 0; i < userCount; i++ {
		us[i] = setup.UserObjWithID(fmt.Sprintf("%v", i))
	}
	return us
}

func getAccountableForOrders(userID meta.ID) []*order.Order {
	// TODO: get user accountable for orders
	orderCount := 100
	os := make([]*order.Order, orderCount)
	for i := 0; i < orderCount; i++ {
		os[i] = setup.OrderObjWithIDs(fmt.Sprintf("%v", i))
	}
	return os
}

func formatToDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func firstLine(lines string) string {
	elems := strings.Split(lines, "\n")
	if len(elems) == 0 {
		return ""
	}
	for _, el := range elems {
		if len(el) > 0 {
			return el
		}
	}
	return ""
}

var (
	templFuncMap = template.FuncMap{
		"formatToDate":   formatToDate,
		"firstLine":      firstLine,
		"States":         getStates,
		"UserEmails":     getUserEmails,
		"Orders":         getOrders,
		"ParentOrder":    getParentOrder,
		"AccountableFor": getAccountableForOrders,
		"SubOrdinates":   getSubOrdinates,
	}

	templOrders = template.Must(template.New("orders").Funcs(templFuncMap).ParseFiles("../../templates/orders.templ.html"))
	templOrder  = template.Must(template.New("order").Funcs(templFuncMap).ParseFiles("../../templates/order.templ.html"))

	templUsers = template.Must(template.New("users").Funcs(templFuncMap).ParseFiles("../../templates/users.templ.html"))
	templUser  = template.Must(template.New("user").Funcs(templFuncMap).ParseFiles("../../templates/user.templ.html"))
)

func serveOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := templOrders.Execute(w, getOrders())
	if err != nil {
		log.Printf("[ERROR]: executing orders html tmpl failed: %s\n", err)
	}
}

func serveOrder(w http.ResponseWriter, r *http.Request) {

	// TODO: get order by ID

	// REMOVEME: START: when getting order by ID
	order := setup.OrderObjWithIDs("1")
	order.Task.Objective += "\nnewlined objective\nsimulates textarea"
	// REMOVEME: END: when getting order by ID

	w.Header().Set("Content-Type", "text/html")
	err := templOrder.Execute(w, order)
	if err != nil {
		log.Printf("[ERROR]: executing html tmpl failed: %s\n", err)
	}
}

func serveUsers(w http.ResponseWriter, r *http.Request) {
	users := getUsers()

	w.Header().Set("Content-Type", "text/html")
	err := templUsers.Execute(w, users)
	if err != nil {
		log.Printf("[ERROR]: executing html tmpl failed: %s\n", err)
	}
}

func serveUser(w http.ResponseWriter, r *http.Request) {

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
	err := templUser.Execute(w, ue)
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))

	http.HandleFunc("GET /orders", serveOrders)
	http.HandleFunc("GET /order/{id}", serveOrder)

	http.HandleFunc("GET /users", serveUsers)
	http.HandleFunc("GET /user/{id}", serveUser)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
