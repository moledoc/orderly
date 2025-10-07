package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/moledoc/orderly/internal/domain/user"
)

var (
	SortUserByID = func(a *user.User, b *user.User) bool {
		return a.GetID() < b.GetID()
	}

	SorterUser = func(sorters ...func(a *user.User, b *user.User) bool) cmp.Option {
		return cmpopts.SortSlices(func(a *user.User, b *user.User) bool {
			for _, comparer := range sorters {
				if !comparer(a, b) {
					return false
				}
			}
			return true
		})
	}

	ComparerUser = func(comparers ...func(a *user.User, b *user.User) bool) cmp.Option {
		return cmp.Comparer(func(a *user.User, b *user.User) bool {
			for _, comparer := range comparers {
				if !comparer(a, b) {
					return false
				}
			}
			return true
		})
	}
)
