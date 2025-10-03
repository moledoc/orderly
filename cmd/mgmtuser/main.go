package main

import (
	"fmt"
	"net/http"

	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtuser"
)

func main() {

	router.RouteUser(mgmtuser.NewServiceMgmtUser(local.NewLocalRepositoryUser()))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
