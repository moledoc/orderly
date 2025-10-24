package tests

import (
	"net/http"
	"sync"
	"testing"

	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	apiorderreq "github.com/moledoc/orderly/tests/api/order/req"
	apiuserreq "github.com/moledoc/orderly/tests/api/user/req"
	"github.com/stretchr/testify/suite"
)

func TestOrderReqSuite(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		router.RouteOrder(mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()))
		router.RouteUser(mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()))
		go http.ListenAndServe(":8080", nil)
		wg.Done()
	}()
	wg.Wait()

	t.Run("OrderAPIReq", func(t *testing.T) {
		suite.Run(t, &OrderSuite{
			API:     apiorderreq.NewOrderAPIReq(),
			UserAPI: apiuserreq.NewUserAPIReq(),
		})
	})
}

func TestOrderReqPerformanceSuite(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		router.RouteOrder(mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()))
		router.RouteUser(mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()))
		go http.ListenAndServe(":8080", nil)
		wg.Done()
	}()
	wg.Wait()

	t.Run("OrderAPIReqPerformance", func(t *testing.T) {
		suite.Run(t, &OrderPerformanceSuite{
			API:     apiorderreq.NewOrderAPIReq(),
			UserAPI: apiuserreq.NewUserAPIReq(),
		})
	})
}
