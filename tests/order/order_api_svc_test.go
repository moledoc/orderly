package tests

import (
	"testing"

	apiordersvc "github.com/moledoc/orderly/tests/api/order/svc"
	apiusersvc "github.com/moledoc/orderly/tests/api/user/svc"
	"github.com/stretchr/testify/suite"
)

func TestOrderSvcSuite(t *testing.T) {

	t.Run("OrderAPISvc", func(t *testing.T) {
		suite.Run(t, &OrderSuite{
			API:     apiordersvc.NewOrderAPISvc(),
			UserAPI: apiusersvc.NewUserAPISvc(),
		})
	})
}

func TestOrderSvcPerformanceSuite(t *testing.T) {

	t.Run("OrderAPISvcPerformance", func(t *testing.T) {
		suite.Run(t, &OrderPerformanceSuite{
			API:     apiordersvc.NewOrderAPISvc(),
			UserAPI: apiusersvc.NewUserAPISvc(),
		})
	})
}
