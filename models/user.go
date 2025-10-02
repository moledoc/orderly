package models

import "time"

type Meta struct {
	Version uint      `json:"version,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
	Deleted bool      `json:"deleted,omitempty"`
}

type Email string

type User struct {
	ID         *uint   `json:"id,omitempty"`
	Name       *string `json:"name,omitempty"`
	Email      *Email  `json:"email,omitempty"`
	Supervisor *Email  `json:"supervisor,omitempty"`
	Meta       *Meta   `json:"meta,omitempty"`
}
