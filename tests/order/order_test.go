package tests

import (
	"github.com/moledoc/orderly/tests/api"
	"github.com/stretchr/testify/suite"
)

type OrderSuite struct {
	suite.Suite
	API api.Order
}

type OrderPerformanceSuite struct {
	suite.Suite
	API api.Order
}
