package user

import (
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/pkg/utils"
)

type Email string

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
