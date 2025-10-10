package compare

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	IgnoreID = cmp.FilterPath(func(pathsteps cmp.Path) bool {
		last := pathsteps.Last().String()
		return last == ".ID"
	}, cmp.Ignore())
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
)

func AssertEqual(t *testing.T, expected any, actual any, opts ...cmp.Option) {
	t.Helper()
	var r DiffReporter
	opts = append(opts, cmp.Reporter(&r))
	assert.Empty(t, cmp.Diff(expected, actual, opts...), r.String())
}

func RequireEqual(t *testing.T, expected any, actual any, opts ...cmp.Option) {
	t.Helper()
	var r DiffReporter
	opts = append(opts, cmp.Reporter(&r))
	require.Empty(t, cmp.Diff(expected, actual, opts...), r.String())
}

func AssertNotEqual(t *testing.T, expected any, actual any, opts ...cmp.Option) {
	t.Helper()
	var r DiffReporter
	opts = append(opts, cmp.Reporter(&r))
	assert.NotEmpty(t, cmp.Diff(expected, actual, opts...), r.String())
}

func RequireNotEqual(t *testing.T, expected any, actual any, opts ...cmp.Option) {
	t.Helper()
	var r DiffReporter
	opts = append(opts, cmp.Reporter(&r))
	require.NotEmpty(t, cmp.Diff(expected, actual, opts...), r.String())
}
