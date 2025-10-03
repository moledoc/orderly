package user

import "github.com/moledoc/orderly/internal/domain/meta"

func (u *User) GetID() meta.ID {
	if u == nil || u.ID == nil {
		return meta.EmptyID()
	}
	return *u.ID
}

func (u *User) GetName() string {
	if u == nil || u.Name == nil {
		return ""
	}
	return *u.Name
}

func (u *User) GetEmail() Email {
	if u == nil || u.Email == nil {
		return ""
	}
	return *u.Email
}

func (u *User) GetSupervisor() Email {
	if u == nil || u.Supervisor == nil {
		return ""
	}
	return *u.Supervisor
}

func (u *User) GetMeta() *meta.Meta {
	if u == nil {
		return nil
	}
	return u.Meta
}
