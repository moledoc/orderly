package compare

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DiffReporter struct {
	path  cmp.Path
	diffs []string
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\t-: %+v\n\t+: %+v\n", r.path, vx, vy))
	}
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}

var (
	IgnoreID = cmp.FilterPath(func(pathsteps cmp.Path) bool {
		last := pathsteps.Last().String()
		return last == ".ID"
	}, cmp.Ignore())
	IgnoreMeta = cmp.FilterPath(func(pathsteps cmp.Path) bool {
		last := pathsteps.Last().String()
		return last == ".Meta"
	}, cmp.Ignore())
)

var (
	IgnorePath = func(paths ...string) cmp.Option {
		return cmp.FilterPath(func(ppaths cmp.Path) bool {
			for _, path := range paths {
				if ppaths[1:].String() == path {
					return true
				}
			}
			return false
		}, cmp.Ignore())
	}
)

var (
	SorterString = cmpopts.SortSlices(func(a string, b string) bool {
		return a < b
	})

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

	CompareUser = func(a *user.User, b *user.User) bool {
		return cmp.Equal(a.Deref(), b.Deref())
	}

	ComparerUser = func(comparers ...func(a *user.User, b *user.User) bool) cmp.Option {
		if len(comparers) == 0 {
			return cmp.Comparer(CompareUser)
		}
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

func AssertEqual(t *testing.T, expected any, actual any, opts ...cmp.Option) {
	opts = append(opts, IgnoreMeta)
	var r DiffReporter
	opts = append(opts, cmp.Reporter(&r))
	assert.True(t, cmp.Equal(expected, actual, opts...), r.String())
}

func RequireEqual(t *testing.T, expected any, actual any, opts ...cmp.Option) {
	opts = append(opts, IgnoreMeta)
	var r DiffReporter
	opts = append(opts, cmp.Reporter(&r))
	require.True(t, cmp.Equal(expected, actual, opts...), r.String())
}
