package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/setup"
)

func getStates() []order.State {
	states := []order.State{}
	for i := order.NotStarted; i <= order.Completed; i++ {
		states = append(states, i)
	}
	return states
}

// MAYBE: cache result
func getUsers() []*user.User {
	resp, _ := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{}) // TODO: handle error
	return resp.GetUsers()
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
	w.Header().Set("Content-Type", "text/html")
	err := templUsers.Execute(w, getUsers())
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
