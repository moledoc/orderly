package models

import "github.com/moledoc/orderly/utils"

type Email string

type User struct {
	ID         *uint   `json:"id,omitempty"`
	Name       *string `json:"name,omitempty"`
	Email      *Email  `json:"email,omitempty"`
	Supervisor *Email  `json:"supervisor,omitempty"`
	Meta       *Meta   `json:"meta,omitempty"`
}

func (u *User) Clone() *User {
	var clone User
	clone.ID = utils.RePtr(u.ID)
	clone.Name = utils.RePtr(u.Name)
	clone.Email = utils.RePtr(u.Email)
	clone.Supervisor = utils.RePtr(u.Supervisor)
	clone.Meta = u.Meta.Clone()

	return &clone
}
