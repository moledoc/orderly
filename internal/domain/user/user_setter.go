package user

import "github.com/moledoc/orderly/internal/domain/meta"

func (u *User) SetID(id meta.ID) {
	if u == nil {
		return
	}
	u.ID = id
}

func (u *User) SetEmail(email Email) {
	if u == nil {
		return
	}
	u.Email = email
}

func (u *User) SetSupervisor(supervisorID meta.ID) {
	if u == nil {
		return
	}
	u.SupervisorID = supervisorID
}

func (u *User) SetMeta(meta *meta.Meta) {
	if u == nil {
		return
	}
	u.Meta = meta
}
