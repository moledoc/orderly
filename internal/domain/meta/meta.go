package meta

import (
	"time"

	"github.com/moledoc/orderly/pkg/utils"
)

type ID string

func EmptyID() ID {
	return ID("")
}

func NewID() ID {
	return ID(utils.RandAlphanum())
}

type Meta struct {
	Version uint      `json:"version,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
	Deleted bool      `json:"deleted,omitempty"`
}

func (m *Meta) VersionIncr() {
	if m == nil {
		return
	}
	m.Version += 1
}

func (m *Meta) Clone() *Meta {
	if m == nil {
		return nil
	}
	clone := Meta{}
	clone.Version = m.Version
	clone.Created = m.Created
	clone.Updated = m.Updated
	clone.Deleted = m.Deleted
	return &clone
}
