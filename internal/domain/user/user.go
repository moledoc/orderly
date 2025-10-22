package user

import (
	"github.com/moledoc/orderly/internal/domain/meta"
)

type Email string

type User struct {
	ID         meta.ID    `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"` // TODO: maybe remove
	Email      Email      `json:"email,omitempty"`
	Supervisor Email      `json:"supervisor,omitempty"`
	Meta       *meta.Meta `json:"meta,omitempty"`
}

func (u *User) Clone() *User {
	if u == nil {
		return nil
	}

	var clone User = User{
		ID:         u.GetID(),
		Name:       u.GetName(),
		Email:      u.GetEmail(),
		Supervisor: u.GetSupervisor(),
		Meta:       u.GetMeta().Clone(),
	}

	return &clone
}
