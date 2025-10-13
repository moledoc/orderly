package tests

import (
	"github.com/moledoc/orderly/tests/api"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	API api.User
}

type UserPerformanceSuite struct {
	suite.Suite
	API api.User
}
