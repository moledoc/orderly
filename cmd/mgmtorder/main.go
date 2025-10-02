package main

import (
	"fmt"
	"net/http"

	"github.com/moledoc/orderly/services/mgmtorder"
)

func main() {
	mgmtorder.New()

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
