package user

import "github.com/moledoc/orderly/internal/domain/meta"

func (u *User) GetID() meta.ID {
	if u == nil {
		return meta.EmptyID()
	}
	return u.ID
}

func (u *User) GetEmail() Email {
	if u == nil {
		return ""
	}
	return u.Email
}

func (u *User) GetSupervisorID() meta.ID {
	if u == nil {
		return ""
	}
	return u.SupervisorID
}

func (u *User) GetMeta() *meta.Meta {
	if u == nil {
		return nil
	}
	return u.Meta
}
