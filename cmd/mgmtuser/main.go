package main

import (
	"fmt"
	"net/http"

	"github.com/moledoc/orderly/services/mgmtuser"
)

func main() {

	mgmtuser.New()

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
