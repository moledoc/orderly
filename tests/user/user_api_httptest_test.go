package tests

import (
	"testing"

	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	apiuserhttptest "github.com/moledoc/orderly/tests/api/user/httptest"
	"github.com/stretchr/testify/suite"
)

func TestUserHTTPTestSuite(t *testing.T) {
	mux := router.RouteUser(mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()))

	t.Run("UserAPIHTTPTest", func(t *testing.T) {
		suite.Run(t, &UserSuite{
			API: apiuserhttptest.NewUserAPIHTTPTest(mux),
		})
	})
}

func TestUserHTTPTestPerformanceSuite(t *testing.T) {
	mux := router.RouteUser(mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()))
	t.Run("UserAPIHTTPTestPerformance", func(t *testing.T) {
		suite.Run(t, &UserPerformanceSuite{
			API: apiuserhttptest.NewUserAPIHTTPTest(mux),
		})
	})
}
