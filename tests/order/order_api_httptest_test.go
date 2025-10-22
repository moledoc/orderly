package tests

import (
	"testing"

	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	apiorderhttptest "github.com/moledoc/orderly/tests/api/order/httptest"
	apiuserhttptest "github.com/moledoc/orderly/tests/api/user/httptest"
	"github.com/stretchr/testify/suite"
)

func TestOrderHTTPTestSuite(t *testing.T) {
	orderMux := router.RouteOrder(mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()))
	userMux := router.RouteUser(mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()))
	t.Run("OrderAPIHTTPTest", func(t *testing.T) {
		suite.Run(t, &OrderSuite{
			API:     apiorderhttptest.NewOrderAPIHTTPTest(orderMux),
			UserAPI: apiuserhttptest.NewUserAPIHTTPTest(userMux),
		})
	})
}

func TestOrderHTTPTestPerformanceSuite(t *testing.T) {
	orderMux := router.RouteOrder(mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()))
	userMux := router.RouteUser(mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()))

	t.Run("OrderAPIHTTPTestPerformance", func(t *testing.T) {
		suite.Run(t, &OrderPerformanceSuite{
			API:     apiorderhttptest.NewOrderAPIHTTPTest(orderMux),
			UserAPI: apiuserhttptest.NewUserAPIHTTPTest(userMux),
		})
	})
}
