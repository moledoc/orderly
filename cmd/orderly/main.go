package main

import (
	"fmt"
	"net/http"

	"github.com/moledoc/orderly/services/mgmtorder"
	"github.com/moledoc/orderly/services/mgmtuser"
)

func main() {
	mgmtorder.New()
	mgmtuser.New()

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
