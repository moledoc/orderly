package tests

import (
	"testing"

	apiusersvc "github.com/moledoc/orderly/tests/api/user/svc"
	"github.com/stretchr/testify/suite"
)

func TestUserSvcSuite(t *testing.T) {
	t.Run("UserAPISvc", func(t *testing.T) {
		suite.Run(t, &UserSuite{
			API: apiusersvc.NewUserAPISvc(),
		})
	})
}

func TestUserSvcPerformanceSuite(t *testing.T) {
	t.Run("UserAPISvcPerformance", func(t *testing.T) {
		suite.Run(t, &UserPerformanceSuite{
			API: apiusersvc.NewUserAPISvc(),
		})
	})
}
