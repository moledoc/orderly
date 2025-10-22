package tests

import (
	"net/http"
	"sync"
	"testing"

	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	apiuserreq "github.com/moledoc/orderly/tests/api/user/req"
	"github.com/stretchr/testify/suite"
)

func TestUserReqSuite(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		router.RouteUser(mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()))
		wg.Done()
		http.ListenAndServe(":8080", nil)
	}()
	wg.Wait()

	t.Run("UserAPIReq", func(t *testing.T) {
		suite.Run(t, &UserSuite{
			API: apiuserreq.NewUserAPIReq(),
		})
	})
}

func TestUserReqPerformanceSuite(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		router.RouteUser(mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()))
		wg.Done()
		http.ListenAndServe(":8080", nil)
	}()
	wg.Wait()

	t.Run("UserAPIReqPerformance", func(t *testing.T) {
		suite.Run(t, &UserPerformanceSuite{
			API: apiuserreq.NewUserAPIReq(),
		})
	})
}
