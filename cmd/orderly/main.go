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
)

// MAYBE: cache result
func getUsers() []*user.User {
	resp, _ := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{}) // TODO: handle error
	return resp.GetUsers()
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
	templNewDelegatedOrder = template.Must(template.New("new_task").Funcs(templFuncMap).ParseFiles(
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
	resp, errr := mgmtorder.GetServiceMgmtOrder().GetOrders(context.Background(), &request.GetOrdersRequest{})
	if errr != nil {
		somethingWentWrong(w, errr)
		return
	}
	type extendedOrder struct {
		Order           *order.Order
		AccountableUser *user.User
	}
	var emails []user.Email
	emailsMap := make(map[user.Email]struct{})
	for _, o := range resp.GetOrders() {
		emailsMap[o.GetAccountable()] = struct{}{}
	}
	for email, _ := range emailsMap {
		emails = append(emails, email)
	}

	respOrdersAccountables, errr := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{
		Emails: emails,
	})
	if errr != nil {
		somethingWentWrong(w, errr)
		return
	}
	accountablesMap := make(map[user.Email]*user.User)
	for _, u := range respOrdersAccountables.GetUsers() {
		_, ok := accountablesMap[u.GetEmail()]
		if ok {
			log.Printf("[WARNING]: multiple users with email %q, using first seen", u.GetEmail())
			continue
		}
		accountablesMap[u.GetEmail()] = u
	}

	eos := make([]*extendedOrder, len(resp.GetOrders()))
	for i, o := range resp.GetOrders() {
		accountable := &user.User{}
		a := accountablesMap[o.GetAccountable()]
		if a != nil {
			accountable = a
		}
		eos[i] = &extendedOrder{
			Order:           o,
			AccountableUser: accountable,
		}
	}

	w.Header().Set("Content-Type", "text/html")
	err := templOrders.Execute(w, eos)
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

	var wg sync.WaitGroup
	var accountable *user.User
	var orders []*order.Order
	var parentOrder *order.Order
	var emails []user.Email
	cherr := make(chan errwrap.Error, 5)
	defer close(cherr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		respOrderAccountable, errr := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{
			Emails: []user.Email{respGetOrderByID.GetOrder().GetAccountable()},
		})
		if errr != nil {
			cherr <- errr

		} else {
			if len(respOrderAccountable.GetUsers()) > 1 {
				log.Printf("[WARNING]: multiple users with same email: %s", respGetOrderByID.GetOrder().GetAccountable())
			}
			if len(respOrderAccountable.GetUsers()) > 0 {
				accountable = respOrderAccountable.GetUsers()[0]
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetOrders, errr := mgmtorder.GetServiceMgmtOrder().GetOrders(context.Background(), &request.GetOrdersRequest{})
		if errr != nil {
			cherr <- errr
		} else {
			orders = respGetOrders.GetOrders()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetParentOrder, errr := mgmtorder.GetServiceMgmtOrder().GetOrderByID(context.Background(), &request.GetOrderByIDRequest{
			ID: respGetOrderByID.GetOrder().GetParentOrderID(),
		})
		if errr != nil {
			cherr <- errr
		} else {
			parentOrder = respGetParentOrder.GetOrder()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetEmails, errr := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{
			Supervisor: mgmtuser.GetServiceMgmtUser().GetRootUser(context.Background()).GetEmail(),
		})
		if errr != nil {
			cherr <- errr
		} else {
			for _, u := range respGetEmails.GetUsers() {
				emails = append(emails, u.GetEmail())
			}
		}
	}()
	wg.Wait()

	type extendedOrder struct {
		Order           *order.Order
		AccountableUser *user.User
		Orders          []*order.Order
		ParentOrder     *order.Order
		Emails          []user.Email
	}

	eo := &extendedOrder{
		Order:           respGetOrderByID.GetOrder(),
		Orders:          orders,
		AccountableUser: accountable,
		ParentOrder:     parentOrder,
		Emails:          emails,
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

	respGetUserByID, errr := mgmtuser.GetServiceMgmtUser().GetUserByID(context.Background(), &request.GetUserByIDRequest{
		ID: meta.ID(r.PathValue("id")),
	})
	if errr != nil {
		err := templSomethingWrong.Execute(w, nil)
		if err != nil {
			log.Printf("[ERROR]: executing SomethingWrong html tmpl failed: %s\n", err)
		}
		return
	}

	var wg sync.WaitGroup
	var supervisor *user.User
	var subordinates []*user.User
	var emails []user.Email
	var orders []*order.Order
	cherr := make(chan errwrap.Error, 5)
	defer close(cherr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetSubOrdinates, errr := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{
			Supervisor: user.Email(respGetUserByID.GetUser().GetEmail()),
		})
		if errr != nil {
			cherr <- errr
		} else {
			subordinates = respGetSubOrdinates.GetUsers()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetSupervisor, errr := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{
			Emails: []user.Email{respGetUserByID.GetUser().GetSupervisor()},
		})
		if errr != nil {
			cherr <- errr

		} else if len(respGetSupervisor.GetUsers()) != 1 {
			cherr <- errwrap.NewError(http.StatusConflict, "incorrect nr of users with email '%s': expected 1, got %d", respGetUserByID.GetUser().GetSupervisor(), len(respGetSupervisor.GetUsers()))
		} else {
			supervisor = respGetSupervisor.GetUsers()[0]
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetEmails, errr := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{})
		if errr != nil {
			cherr <- errr

		} else {
			for _, u := range respGetEmails.GetUsers() {
				emails = append(emails, u.GetEmail())
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetOrders, errr := mgmtorder.GetServiceMgmtOrder().GetOrders(context.Background(), &request.GetOrdersRequest{
			Accountable: user.Email(respGetUserByID.GetUser().GetEmail()),
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
		User           *user.User
		SupervisorID   meta.ID
		SubOrdinates   []*user.User
		AccountableFor []*order.Order
		Emails         []user.Email
	}

	ue := &userExtended{
		User:           respGetUserByID.GetUser(),
		SupervisorID:   supervisor.GetID(),
		SubOrdinates:   subordinates,
		AccountableFor: orders,
		Emails:         emails,
	}

	w.Header().Set("Content-Type", "text/html")
	err := templUser.Execute(w, ue)
	if err != nil {
		log.Printf("[ERROR]: executing html tmpl failed: %s\n", err)
	}
}

func serveNewUser(w http.ResponseWriter, _ *http.Request) {
	var wg sync.WaitGroup
	var emails []user.Email
	cherr := make(chan errwrap.Error, 5)
	defer close(cherr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetEmails, errr := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{})
		if errr != nil {
			cherr <- errr

		} else {
			for _, u := range respGetEmails.GetUsers() {
				emails = append(emails, u.GetEmail())
			}
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

	type extendedUser struct {
		Emails []user.Email
	}

	eu := &extendedUser{
		Emails: emails,
	}

	err := templNewUser.Execute(w, eu)
	if err != nil {
		log.Printf("[ERROR]: executing new_user html tmpl failed: %s\n", err)
	}
}

func serveNewDelegatedOrder(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var emails []user.Email
	var orders []*order.Order
	cherr := make(chan errwrap.Error, 5)
	defer close(cherr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetOrders, errr := mgmtorder.GetServiceMgmtOrder().GetOrders(context.Background(), &request.GetOrdersRequest{})
		if errr != nil {
			cherr <- errr

		} else {
			orders = respGetOrders.GetOrders()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetEmails, errr := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{
			Supervisor: mgmtuser.GetServiceMgmtUser().GetRootUser(context.Background()).GetEmail(),
		})
		if errr != nil {
			cherr <- errr

		} else {
			for _, u := range respGetEmails.GetUsers() {
				emails = append(emails, u.GetEmail())
			}
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

	type extendedOrder struct {
		Orders    []*order.Order
		Emails    []user.Email
		Delegated bool
	}
	eo := &extendedOrder{
		Orders:    orders,
		Emails:    emails,
		Delegated: true,
	}
	err := templNewDelegatedOrder.Execute(w, eo)
	if err != nil {
		log.Printf("[ERROR]: executing new_task html tmpl failed: %s\n", err)
	}
}

func serveNewOrder(w http.ResponseWriter, _ *http.Request) {
	var wg sync.WaitGroup
	var users []*user.User
	var emails []user.Email
	var orders []*order.Order
	cherr := make(chan errwrap.Error, 5)
	defer close(cherr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetOrders, errr := mgmtorder.GetServiceMgmtOrder().GetOrders(context.Background(), &request.GetOrdersRequest{})
		if errr != nil {
			cherr <- errr
		} else {
			orders = respGetOrders.GetOrders()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		respGetUsers, errr := mgmtuser.GetServiceMgmtUser().GetUsers(context.Background(), &request.GetUsersRequest{
			Supervisor: mgmtuser.GetServiceMgmtUser().GetRootUser(context.Background()).GetEmail(),
		})
		if errr != nil {
			cherr <- errr
		} else {
			users = respGetUsers.GetUsers()
			for _, u := range users {
				emails = append(emails, u.GetEmail())
			}
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
	type extendedOrder struct {
		Orders    []*order.Order
		Users     []*user.User
		Emails    []user.Email
		Delegated bool
	}
	eo := &extendedOrder{
		Orders:    orders,
		Users:     users,
		Emails:    emails,
		Delegated: false,
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
	http.HandleFunc("GET /order/new/task", serveNewDelegatedOrder)

	http.HandleFunc("GET /users", serveUsers)
	http.HandleFunc("GET /user/{id}", serveUser)
	http.HandleFunc("GET /user/new", serveNewUser)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
