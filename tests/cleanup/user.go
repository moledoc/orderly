package cleanup

import (
	"context"
	"testing"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/tests/api"
)

func User(t *testing.T, api api.User, u *user.User) {
	t.Helper()
	t.Cleanup(func() {
		id := meta.ID(u.GetID())
		if len(id) == 0 {
			return
		}
		_, err := api.DeleteUser(t, context.Background(), &request.DeleteUserRequest{
			ID: id,
		})
		if err != nil {
			t.Logf("[WARNING]: user '%v' wasn't cleaned up: %s\n", id, err)
		}
	})
}
