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
	clone := Meta{
		Version: m.Version,
		Created: m.Created,
		Updated: m.Updated,
	}
	return &clone
}
