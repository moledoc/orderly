package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/pkg/utils"
)

var (
	SortOrderByID = func(a *order.Order, b *order.Order) bool {
		return a.GetOrder().GetID() < b.GetOrder().GetID()
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
	SortState = func(a *order.State, b *order.State) bool {
		return utils.Deref(a) < utils.Deref(b)
	}

	SorterState = func(sorters ...func(a *order.State, b *order.State) bool) cmp.Option {
		return cmpopts.SortSlices(func(a *order.State, b *order.State) bool {
			for _, comparer := range sorters {
				if !comparer(a, b) {
					return false
				}
			}
			return true
		})
	}

	ComparerState = func(comparers ...func(a *order.State, b *order.State) bool) cmp.Option {
		return cmp.Comparer(func(a *order.State, b *order.State) bool {
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
	SortOrderByID = func(a *order.Order, b *order.Order) bool {
		return a.GetID() < b.GetID()
	}

	SortOrderByAccountable = func(a *order.Order, b *order.Order) bool {
		return a.GetAccountable() < b.GetAccountable()
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
	SortSitRepByID = func(a *order.SitRep, b *order.SitRep) bool {
		return a.GetID() < b.GetID()
	}
	SortSitRepByDateTime = func(a *order.SitRep, b *order.SitRep) bool {
		return a.GetDateTime().Sub(b.GetDateTime()) < 0
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
