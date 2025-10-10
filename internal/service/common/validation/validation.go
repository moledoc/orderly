package validation

import (
	"net/http"
	"net/mail"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/pkg/utils"
)

type IgnoreField int

const (
	IgnoreNothing IgnoreField = 0
	IgnoreID      IgnoreField = 1 << iota
	IgnoreEmpty
)

func IsFieldIgnored(field IgnoreField, bitmask IgnoreField) bool {
	return field&bitmask == field
}

func IsIgnoreEmpty(field any, bitmask IgnoreField) bool {
	return IgnoreEmpty&bitmask == IgnoreEmpty && utils.IsZeroValue(field)
}

func ValidateID(id meta.ID) errwrap.Error {
	if len(id) != utils.RandAlphanumLen {
		return errwrap.NewError(http.StatusBadRequest, "invalid id length")
	}
	for _, c := range id {
		if !('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9') {
			return errwrap.NewError(http.StatusBadRequest, "disallowed character ('%c') in id", c)
		}
	}
	return nil
}

func ValidateEmail(email user.Email) errwrap.Error {
	if len(email) == 0 {
		return errwrap.NewError(http.StatusBadRequest, "invalid email length")
	}
	_, err := mail.ParseAddress(string(email))
	if err != nil {
		return errwrap.NewError(http.StatusBadRequest, "invalid email")
	}
	return nil
}
