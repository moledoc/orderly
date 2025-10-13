package compare_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/compare"
	"github.com/moledoc/orderly/tests/setup"
)

func TestCompare_User(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		compare.RequireEqual(t, &user.User{}, &user.User{})
	})
	t.Run("same", func(t *testing.T) {
		u := setup.UserObjWithID()
		compare.RequireEqual(t, u, u)
	})

	t.Run("sameWithDifferentID", func(t *testing.T) {
		u1 := setup.UserObjWithID()
		u2 := utils.RePtr(u1)
		u2.SetID(meta.NewID())

		t.Run("assertOnID", func(t *testing.T) {
			compare.RequireNotEqual(t, u1, u2)
		})

		t.Run("ignoreID", func(t *testing.T) {
			opts := []cmp.Option{
				compare.IgnorePaths("ID"),
				// compare.IgnoreID, // or use this
			}
			compare.RequireEqual(t, u1, u2, opts...)
		})

	})
	t.Run("differentWithDiffIDs", func(t *testing.T) {
		compare.RequireNotEqual(t, setup.UserObjWithID("u1"), setup.UserObjWithID("u2"))
	})

	t.Run("differentWithSameIDs", func(t *testing.T) {
		u1 := setup.UserObjWithID()
		u2 := utils.RePtr(u1)
		u2.SetName("new name")

		compare.RequireNotEqual(t, u1, u2)

		t.Run("ignoreChange", func(t *testing.T) {
			opts := []cmp.Option{
				compare.IgnorePaths("Name"),
			}
			compare.RequireEqual(t, u1, u2, opts...)
		})
	})
}
