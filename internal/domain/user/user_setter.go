package user

import "github.com/moledoc/orderly/internal/domain/meta"

func (u *User) SetID(id meta.ID) {
	if u == nil || u.ID == nil {
		return
	}
	*u.ID = id
}

func (u *User) SetName(name string) {
	if u == nil || u.Name == nil {
		return
	}
	*u.Name = name
}

func (u *User) SetEmail(email Email) {
	if u == nil || u.Email == nil {
		return
	}
	*u.Email = email
}

func (u *User) SetSupervisor(supervisor Email) {
	if u == nil || u.Supervisor == nil {
		return
	}
	*u.Supervisor = supervisor
}

func (u *User) SetMeta(meta *meta.Meta) {
	if u == nil {
		return
	}
	u.Meta = meta
}
