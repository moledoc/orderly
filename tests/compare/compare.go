package compare

import (
	"fmt"
	"slices"
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
		r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\texpected: %+v\n\tactual: %+v\n", r.path, vx, vy))
	}
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}

var (
	IgnoreMeta = cmp.FilterPath(func(pathsteps cmp.Path) bool {
		last := pathsteps.Last().String()
		return last == ".Meta"
	}, cmp.Ignore())
	IgnorePaths = func(paths ...string) cmp.Option {
		return cmp.FilterPath(func(path cmp.Path) bool {
			return slices.Contains(paths, path.String())
		}, cmp.Ignore())
	}
	// IgnoreUserFields = func(fields ...string) cmp.Option {
	// 	return cmpopts.IgnoreFields(user.User{}, fields...)
	// }
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

func AssertEqual(t *testing.T, expected any, actual any, opts ...cmp.Option) {
	opts = append(opts, IgnoreMeta)
	// var r DiffReporter
	// opts = append(opts, cmp.Reporter(&r))
	assert.Empty(t, cmp.Diff(expected, actual, opts...) /*, r.String()*/)
}

func RequireEqual(t *testing.T, expected any, actual any, opts ...cmp.Option) {
	opts = append(opts, IgnoreMeta)
	// var r DiffReporter
	// opts = append(opts, cmp.Reporter(&r))
	require.Empty(t, cmp.Diff(expected, actual, opts...) /*, r.String()*/)
}
