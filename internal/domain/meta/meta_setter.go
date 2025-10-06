package meta

import (
	"time"
)

// func (m *Meta) SetVersion(version uint) {
// 	if m == nil {
// 		return
// 	}
// 	m.Version = version
// }

func (m *Meta) SetCreated(created time.Time) {
	if m == nil {
		return
	}
	m.Created = created
}

func (m *Meta) SetUpdated(updated time.Time) {
	if m == nil {
		return
	}
	m.Updated = updated
}
