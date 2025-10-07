package user

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/pkg/utils"
)

type Email string

type userDerefFields struct {
	ID         meta.ID
	Name       string
	Email      Email
	Supervisor Email
	Meta       meta.Meta
}

type User struct {
	ID         *meta.ID   `json:"id,omitempty"`
	Name       *string    `json:"name,omitempty"`
	Email      *Email     `json:"email,omitempty"`
	Supervisor *Email     `json:"supervisor,omitempty"`
	Meta       *meta.Meta `json:"meta,omitempty"`
}

func (u *User) Clone() *User {
	if u == nil {
		return nil
	}

	var clone User = User{
		ID:         utils.Ptr(u.GetID()),
		Name:       utils.Ptr(u.GetName()),
		Email:      utils.Ptr(u.GetEmail()),
		Supervisor: utils.Ptr(u.GetSupervisor()),
		Meta:       u.GetMeta().Clone(),
	}

	return &clone
}

func (u *User) Deref() *userDerefFields {
	if u == nil {
		return nil
	}

	var deref userDerefFields = userDerefFields{
		ID:         u.GetID(),
		Name:       u.GetName(),
		Email:      u.GetEmail(),
		Supervisor: u.GetSupervisor(),
		Meta:       *u.GetMeta().Clone(),
	}

	return &deref
}
