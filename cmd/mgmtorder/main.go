package main

import (
	"fmt"
	"net/http"

	"github.com/moledoc/orderly/internal/repository/local"
	"github.com/moledoc/orderly/internal/router"
	"github.com/moledoc/orderly/internal/service/mgmtorder"
)

func main() {

	router.RouteOrder(mgmtorder.NewServiceMgmtOrder(local.NewLocalRepositoryOrder()))

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
