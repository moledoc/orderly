package compare_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/order"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
)

func TestCompare_Order(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		compare.RequireEqual(t, &order.Order{}, &order.Order{})
	})
	t.Run("same", func(t *testing.T) {
		u := setup.OrderObjWithIDs()
		compare.RequireEqual(t, u, u)
	})

	t.Run("different", func(t *testing.T) {
		compare.RequireNotEqual(t, setup.OrderObjWithIDs("u1"), setup.OrderObjWithIDs("u2"))
	})

	t.Run("differentParentID", func(t *testing.T) {
		u1 := setup.OrderObjWithIDs()
		u2 := utils.RePtr(u1)
		u2.SetParentOrderID(meta.NewID())

		compare.RequireNotEqual(t, u1, u2)

		t.Run("ignoreChange", func(t *testing.T) {
			opts := []cmp.Option{
				compare.IgnorePaths("ParentOrderID"),
			}
			compare.RequireEqual(t, u1, u2, opts...)
		})
	})

	t.Run("differentTask", func(t *testing.T) {
		t1 := setup.TaskObj()
		t2 := utils.RePtr(t1)

		t2.SetObjective("this is new objective")

		compare.RequireNotEqual(t, t1, t2)

		t.Run("ignoreChange", func(t *testing.T) {
			opts := []cmp.Option{
				compare.IgnorePaths("Objective"),
			}
			compare.RequireEqual(t, t1, t2, opts...)
		})
	})

	t.Run("differentSitRep", func(t *testing.T) {
		sr1 := setup.SitrepObj()
		sr2 := utils.RePtr(sr1)

		sr2.SetActions("this is new action")

		compare.RequireNotEqual(t, sr1, sr2)

		t.Run("ignoreChange", func(t *testing.T) {
			opts := []cmp.Option{
				compare.IgnorePaths("Actions"),
			}
			compare.RequireEqual(t, sr1, sr2, opts...)
		})
	})

}
