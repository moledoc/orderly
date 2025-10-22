package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/moledoc/orderly/tests/setup"
)

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

func getOrders() ([]*order.Order, errwrap.Error) {
	resp, err := mgmtorder.GetServiceMgmtOrder().GetOrders(context.Background(), &request.GetOrdersRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetOrders(), nil
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

func somethingWentWrong(w http.ResponseWriter, err errwrap.Error) {
	log.Printf("[ERROR]: %s\n", err)

	w.Header().Set("Content-Type", "text/html")
	e := templSomethingWrong.Execute(w, nil)
	if e != nil {
		log.Printf("[ERROR]: executing SomethingWrong html tmpl failed: %s\n", e)
	}
}

var (
	templFuncMap = template.FuncMap{
		"formatToDate": formatToDate,
		"firstLine":    firstLine,
		"States":       order.ListStates,
		"UserEmails":   getUserEmails,  // REMOVEME: move to handleFunc
		"ParentOrder":  getParentOrder, // REMOVEME: move to handleFunc
	}

	templOrders = template.Must(template.New("orders").Funcs(templFuncMap).ParseFiles(
		"./templates/header.templ.html",
		"./templates/footer.templ.html",
		"./templates/orders.templ.html",
	))
	templOrder = template.Must(template.New("order").Funcs(templFuncMap).ParseFiles(
		"./templates/header.templ.html",
		"./templates/footer.templ.html",
		"./templates/order.templ.html",
	))

	templUsers = template.Must(template.New("users").Funcs(templFuncMap).ParseFiles(
		"./templates/header.templ.html",
		"./templates/footer.templ.html",
		"./templates/users.templ.html",
	))
	templUser = template.Must(template.New("user").Funcs(templFuncMap).ParseFiles(
		"./templates/header.templ.html",
		"./templates/footer.templ.html",
		"./templates/user.templ.html",
	))
	templNewUser = template.Must(template.New("new_user").Funcs(templFuncMap).ParseFiles(
		"./templates/header.templ.html",
		"./templates/footer.templ.html",
		"./templates/new_user.templ.html",
	))
	templNewOrder = template.Must(template.New("new_order").Funcs(templFuncMap).ParseFiles(
		"./templates/header.templ.html",
		"./templates/footer.templ.html",
		"./templates/new_order.templ.html",
		"./templates/new_task.templ.html",
	))
	templNewTask = template.Must(template.New("new_task").Funcs(templFuncMap).ParseFiles(
		"./templates/new_task.templ.html",
	))
	templHome = template.Must(template.New("home").Funcs(templFuncMap).ParseFiles(
		"./templates/header.templ.html",
		"./templates/footer.templ.html",
		"./templates/home.templ.html",
	))

	templSomethingWrong = template.Must(template.New("something_wrong").Funcs(templFuncMap).ParseFiles(
		"./templates/header.templ.html",
		"./templates/footer.templ.html",
		"./templates/something_wrong.templ.html",
	))
)

func serveOrders(w http.ResponseWriter, r *http.Request) {
	os, errr := getOrders()
	if errr != nil {
		somethingWentWrong(w, errr)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err := templOrders.Execute(w, os)
	if err != nil {
		log.Printf("[ERROR]: executing orders html tmpl failed: %s\n", err)
	}
}

func serveHome(w http.ResponseWriter, _ *http.Request) {
	err := templHome.Execute(w, nil)
	if err != nil {
		log.Printf("[ERROR]: executing home html tmpl failed: %s\n", err)
	}
}

func serveOrder(w http.ResponseWriter, r *http.Request) {

	respGetOrderByID, errr := mgmtorder.GetServiceMgmtOrder().GetOrderByID(context.Background(), &request.GetOrderByIDRequest{
		ID: meta.ID(r.PathValue("id")),
	})
	if errr != nil {
		somethingWentWrong(w, errr)
		return
	}

	type extendedOrder struct {
		Order  *order.Order
		Orders []*order.Order
	}

	os, errr := getOrders()
	if errr != nil {
		somethingWentWrong(w, errr)
		return
	}

	eo := &extendedOrder{
		Order:  respGetOrderByID.GetOrder(),
		Orders: os,
	}

	w.Header().Set("Content-Type", "text/html")
	err := templOrder.Execute(w, eo)
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

	var wg sync.WaitGroup
	var u *user.User
	var subordinates []*user.User
	var orders []*order.Order
	cherr := make(chan errwrap.Error, 2)
	defer close(cherr)

	wg.Add(3)
	go func() {
		defer wg.Done()
		respGetUserByID, errr := mgmtuser.GetServiceMgmtUser().GetUserByID(context.Background(), &request.GetUserByIDRequest{
			ID: meta.ID(r.PathValue("id")),
		})
		if errr != nil {
			cherr <- errr
		} else {
			u = respGetUserByID.GetUser()
		}
	}()
	go func() {
		defer wg.Done()
		respGetSubOrdinates, errr := mgmtuser.GetServiceMgmtUser().GetUserSubOrdinates(context.Background(), &request.GetUserSubOrdinatesRequest{
			ID: meta.ID(r.PathValue("id")),
		})
		if errr != nil {
			cherr <- errr
		} else {
			subordinates = respGetSubOrdinates.GetSubOrdinates()
		}
	}()
	go func() {
		defer wg.Done()
		respGetOrders, errr := mgmtorder.GetServiceMgmtOrder().GetOrders(context.Background(), &request.GetOrdersRequest{
			Accountable: user.Email(r.PathValue("accountable")),
		})
		if errr != nil {
			cherr <- errr
		} else {
			orders = respGetOrders.GetOrders()
		}
	}()
	wg.Wait()

	if len(cherr) > 0 {
		err := templSomethingWrong.Execute(w, nil)
		if err != nil {
			log.Printf("[ERROR]: executing SomethingWrong html tmpl failed: %s\n", err)
		}
		return
	}

	type userExtended struct {
		*user.User
		SupervisorID   string
		SubOrdinates   []*user.User
		AccountableFor []*order.Order
	}

	ue := &userExtended{
		User: u,
		// TODO: get supervisor ID by getting user by email
		SubOrdinates:   subordinates,
		AccountableFor: orders,
	}

	w.Header().Set("Content-Type", "text/html")
	err := templUser.Execute(w, ue)
	if err != nil {
		log.Printf("[ERROR]: executing html tmpl failed: %s\n", err)
	}
}

func serveNewUser(w http.ResponseWriter, _ *http.Request) {
	err := templNewUser.Execute(w, nil)
	if err != nil {
		log.Printf("[ERROR]: executing new_user html tmpl failed: %s\n", err)
	}
}

func serveNewTask(w http.ResponseWriter, r *http.Request) {
	os, errr := getOrders()
	if errr != nil {
		somethingWentWrong(w, errr)
		return
	}

	type extendedOrder struct {
		Orders    []*order.Order
		Delegated bool
	}
	eo := &extendedOrder{
		Orders:    os,
		Delegated: true,
	}
	err := templNewTask.Execute(w, eo)
	if err != nil {
		log.Printf("[ERROR]: executing new_task html tmpl failed: %s\n", err)
	}
}

func serveNewOrder(w http.ResponseWriter, _ *http.Request) {
	os, errr := getOrders()
	if errr != nil {
		somethingWentWrong(w, errr)
		return
	}

	type extendedOrder struct {
		Orders    []*order.Order
		Delegated bool
	}
	eo := &extendedOrder{
		Orders: os,
	}

	err := templNewOrder.Execute(w, eo)
	if err != nil {
		log.Printf("[ERROR]: executing new_user html tmpl failed: %s\n", err)
	}
}

func main() {

	router.Route(&router.Service{
		MgmtOrder: mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()),
		MgmtUser:  mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()),
	})

	// http.HandleFunc("GET /", serveLogin) // NOTE: login

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("GET /", serveHome)

	http.HandleFunc("GET /orders", serveOrders)
	http.HandleFunc("GET /order/{id}", serveOrder)
	http.HandleFunc("GET /order/new", serveNewOrder)
	http.HandleFunc("GET /order/new/task", serveNewTask)

	http.HandleFunc("GET /users", serveUsers)
	http.HandleFunc("GET /user/{id}", serveUser)
	http.HandleFunc("GET /user/new", serveNewUser)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
