package state

import (
	"context"
	"log"
	"sync"

	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
)

type Label string

type APIs struct {
	User  mgmtuser.ServiceMgmtUserAPI
	Order mgmtorder.ServiceMgmtOrderAPI
}

type orderHierarchy struct {
	Order           *order.Order
	DelegatedOrders []*order.Order
}

type State struct {
	mu             sync.Mutex
	orderHierarchy *orderHierarchy
	orders         map[Label]*order.Order
	users          map[Label]*user.User
	apis           APIs
}

func New(apis APIs) *State {
	return &State{
		apis: apis,
		orders: map[Label]*order.Order{
			"root": apis.Order.GetRootOrder(context.Background()),
		},
		users: map[Label]*user.User{
			"root": apis.User.GetRootUser(context.Background()),
		},
	}
}

func (s *State) createUser() {
	for label, u := range s.users {
		if len(u.GetID()) > 0 { // NOTE: existing user
			continue
		}
		respPostUser, err := s.apis.User.PostUser(context.Background(), &request.PostUserRequest{
			User: u,
		})
		if err != nil {
			log.Printf("[ERROR]: failed to create user with label %q: %s\n", label, err)
		} else {
			s.users[label] = respPostUser.GetUser()
		}
	}
}

func (s *State) createOrders() {
	for label, o := range s.orders {
		if len(o.GetID()) > 0 { // NOTE: existing order
			continue
		}
		o.SetDelegatedOrders(nil)
		respPostOrder, err := s.apis.Order.PostOrder(context.Background(), &request.PostOrderRequest{
			Order: o,
		})
		if err != nil {
			log.Printf("[ERROR]: failed to create order with label %q: %s\n", label, err)
		} else {
			s.orders[label] = respPostOrder.GetOrder()
		}
	}
}
