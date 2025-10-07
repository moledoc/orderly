package cleanup

import (
	"context"
	"testing"

	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/request"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
	"github.com/moledoc/orderly/tests/api"
)

func User(t *testing.T, api api.User, u *user.User) {
	t.Cleanup(func() {
		id := meta.ID(u.GetID())
		_, err := api.DeleteUser(t, context.Background(), &request.DeleteUserRequest{
			ID: utils.Ptr(id),
		})
		if err != nil {
			t.Logf("[WARNING]: user '%v' wasn't cleaned up: %s\n", id, err)
		}
	})
}
