package user

import (
	"github.com/moledoc/orderly/internal/domain/meta"
)

type Email string

type User struct {
	ID           meta.ID    `json:"id,omitempty"`
	Email        Email      `json:"email,omitempty"`
	SupervisorID meta.ID    `json:"supervisor_id,omitempty"`
	Meta         *meta.Meta `json:"meta,omitempty"`
}

func (u *User) Clone() *User {
	if u == nil {
		return nil
	}

	var clone User = User{
		ID:           u.GetID(),
		Email:        u.GetEmail(),
		SupervisorID: u.GetSupervisorID(),
		Meta:         u.GetMeta().Clone(),
	}

	return &clone
}
