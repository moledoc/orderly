package models

import (
	"time"
)

type Meta struct {
	Version uint      `json:"version,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
	Deleted bool      `json:"deleted,omitempty"`
}

func (m *Meta) Clone() *Meta {
	clone := Meta{}
	clone.Version = m.Version
	clone.Created = m.Created
	clone.Updated = m.Updated
	clone.Deleted = m.Deleted
	return &clone
}
