package meta

import (
	"time"
)

func (m *Meta) GetVersion() uint {
	if m == nil {
		return 0
	}
	return m.Version
}

func (m *Meta) GetCreated() time.Time {
	if m == nil {
		return time.Time{}
	}
	return m.Created
}

func (m *Meta) GetUpdated() time.Time {
	if m == nil {
		return time.Time{}
	}
	return m.Updated
}
