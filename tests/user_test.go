package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	API UserAPI
}

type UserSvc struct{} // NOTE: tests service layer methods
type UserReq struct{} // NOTE: tests service through HTTP requests

func TestUserSuite(t *testing.T) {
	t.Run("UserSvc", func(t *testing.T) {
		suite.Run(t, &UserSuite{
			API: new(UserSvc),
		})
	})
	t.Run("UserReq", func(t *testing.T) {
		suite.Run(t, &UserSuite{
			API: new(UserReq),
		})
	})
}
