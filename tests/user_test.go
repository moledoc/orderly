package tests

import (
	"net/http"
	"testing"

	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	API UserAPI
}

type UserAPISvc struct { // NOTE: tests service layer methods
	Svc mgmtuser.ServiceMgmtUserAPI
}

func NewUserAPISvc() *UserAPISvc {
	// TODO: local vs db
	return &UserAPISvc{
		Svc: mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()),
	}
}

type UserAPIReq struct { // NOTE: tests service through HTTP requests
	// TODO: local vs db
	HttpClient *http.Client
	BaseURL    string
}

func NewUserAPIReq() *UserAPIReq {
	// TODO: local vs db
	return &UserAPIReq{
		HttpClient: &http.Client{},
		BaseURL:    "http://127.0.0.1:8080",
	}
}

var (
	_ UserAPI = (*UserAPISvc)(nil)
	_ UserAPI = (*UserAPIReq)(nil)
)

func TestUserSvcSuite(t *testing.T) {
	t.Run("UserAPISvc", func(t *testing.T) {
		suite.Run(t, &UserSuite{
			API: NewUserAPISvc(),
		})
	})
}

func TestUserReqSuite(t *testing.T) {
	t.Run("UserAPIReq", func(t *testing.T) {
		suite.Run(t, &UserSuite{
			API: NewUserAPIReq(),
		})
	})
}
