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

var (
	SortTaskByID = func(a *order.Task, b *order.Task) bool {
		return a.GetID() < b.GetID()
	}

	SortTaskByAccountable = func(a *order.Task, b *order.Task) bool {
		return a.GetAccountable() < b.GetAccountable()
	}

	SorterTask = func(sorters ...func(a *order.Task, b *order.Task) bool) cmp.Option {
		return cmpopts.SortSlices(func(a *order.Task, b *order.Task) bool {
			for _, comparer := range sorters {
				if !comparer(a, b) {
					return false
				}
			}
			return true
		})
	}

	ComparerTask = func(comparers ...func(a *order.Task, b *order.Task) bool) cmp.Option {
		return cmp.Comparer(func(a *order.Task, b *order.Task) bool {
			for _, comparer := range comparers {
				if !comparer(a, b) {
					return false
				}
			}
			return true
		})
	}
)

var (
	SortSitRepByID = func(a *order.SitRep, b *order.SitRep) bool {
		return a.GetID() < b.GetID()
	}

	SorterSitRep = func(sorters ...func(a *order.SitRep, b *order.SitRep) bool) cmp.Option {
		return cmpopts.SortSlices(func(a *order.SitRep, b *order.SitRep) bool {
			for _, comparer := range sorters {
				if !comparer(a, b) {
					return false
				}
			}
			return true
		})
	}

	ComparerSitRep = func(comparers ...func(a *order.SitRep, b *order.SitRep) bool) cmp.Option {
		return cmp.Comparer(func(a *order.SitRep, b *order.SitRep) bool {
			for _, comparer := range comparers {
				if !comparer(a, b) {
					return false
				}
			}
			return true
		})
	}
)
