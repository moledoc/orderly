package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/moledoc/orderly/internal/domain/order"
)

var (
	SortOrderByID = func(a *order.Order, b *order.Order) bool {
		return a.GetTask().GetID() < b.GetTask().GetID()
	}

	SorterOrder = func(sorters ...func(a *order.Order, b *order.Order) bool) cmp.Option {
		return cmpopts.SortSlices(func(a *order.Order, b *order.Order) bool {
			for _, comparer := range sorters {
				if !comparer(a, b) {
					return false
				}
			}
			return true
		})
	}

	ComparerOrder = func(comparers ...func(a *order.Order, b *order.Order) bool) cmp.Option {
		return cmp.Comparer(func(a *order.Order, b *order.Order) bool {
			for _, comparer := range comparers {
				if !comparer(a, b) {
					return false
				}
			}
			return true
		})
	}
)
