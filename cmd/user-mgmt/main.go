package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type meta struct {
	Version int       `json:"version"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type user struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Supervisor string `json:"supervisor"`
	Meta       *meta  `json:"meta"`
}

type storageAPI interface {
	Close()
	Read(ctx context.Context, ID int) (user *user, err error)
	Write(ctx context.Context, user *user) (err error)
}

func handlePostUser(w http.ResponseWriter, r *http.Request) {}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) {}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {}

func handleGetUserSubOrdinates(w http.ResponseWriter, r *http.Request) {}

func handlePatchUser(w http.ResponseWriter, r *http.Request) {}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {}

func main() {
	http.HandleFunc("POST /user", handlePostUser)
	http.HandleFunc("GET /user/{id}", handleGetUserByID)
	http.HandleFunc("GET /users", handleGetUsers)
	http.HandleFunc("GET /user/{id}/subordinates", handleGetUserSubOrdinates)
	http.HandleFunc("PATCH /user/{id}", handlePatchUser)
	http.HandleFunc("DELETE /user/{id}", handleDeleteUser)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
