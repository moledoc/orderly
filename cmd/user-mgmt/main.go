package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ierror interface {
	String() string
}

type erro struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

func (e *erro) String() string {
	return fmt.Sprintf(`{"code": %d, "message": %q}`, e.Code, e.Message)
}

func NewError(code uint, format string, a ...any) ierror {
	return &erro{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}

type Meta struct {
	Version *uint      `json:"version"`
	Created *time.Time `json:"created"`
	Updated *time.Time `json:"updated"`
	Deleted *bool      `json:"deleted"`
}

type User struct {
	ID         *uint   `json:"id"`
	Name       *string `json:"name"`
	Email      *string `json:"email"`
	Supervisor *string `json:"supervisor"`
	Meta       *Meta   `json:"meta"`
}

type Action int

const (
	CREATE Action = iota
	UPDATE
	SOFTDELETE
	HARDDELETE
)

type storageAPI interface {
	Close()
	Read(ctx context.Context, id int) (*User, ierror)
	Write(ctx context.Context, action Action, user *User) ierror
}

type LocalStorage map[uint][]*User

func (s LocalStorage) Close() {
	s = nil
}

func (s LocalStorage) Read(ctx context.Context, id uint) (*User, ierror) {
	if s == nil {
		return nil, NewError(http.StatusInternalServerError, "localstorage not initialized for read")
	}
	us, ok := s[id]
	if !ok || len(us) == 0 {
		return nil, NewError(http.StatusNotFound, "not found during read")
	}
	return us[len(us)-1], nil
}

func (s LocalStorage) Write(ctx context.Context, action Action, user *User) ierror {
	if s == nil {
		return NewError(http.StatusInternalServerError, "localstorage not initialized for write")
	}
	if user == nil || user.ID == nil {
		return NewError(http.StatusInternalServerError, "invalid user object in write")
	}

	us, ok := s[*user.ID]

	switch action {

	case CREATE:
		if ok || len(us) > 0 {
			return NewError(http.StatusConflict, "already exists during write")
		}
		s[*user.ID] = append(s[*user.ID], user)
		return nil

	case UPDATE:
		if !ok || len(us) == 0 {
			return NewError(http.StatusNotFound, "not found during write")
		}
		var updUser User = *(us[len(us)-1])
		updated := false
		if user.Name != nil {
			updUser.Name = user.Name
			updated = true
		}
		if user.Email != nil {
			updUser.Email = user.Email
			updated = true
		}
		if user.Supervisor != nil {
			updUser.Supervisor = user.Supervisor
			updated = true
		}
		if updated {
			now := time.Now().UTC()
			updUser.Meta.Updated = &now
			*updUser.Meta.Version += 1
			s[*user.ID] = append(s[*user.ID], &updUser)
		}
		return nil

	case SOFTDELETE:
		if ok {
			for _, u := range us {
				b := true
				u.Meta.Deleted = &b
			}
		}
		return nil

	case HARDDELETE:
		if ok {
			delete(s, *user.ID)
		}
		return nil

	default:
		return NewError(http.StatusInternalServerError, "undefined write action")
	}
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
	http.HandleFunc("PATCH /user", handlePatchUser)
	http.HandleFunc("DELETE /user/{id}", handleDeleteUser)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
