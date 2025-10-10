package common

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/moledoc/orderly/internal/domain/errwrap"
	"github.com/moledoc/orderly/internal/domain/meta"
	"github.com/moledoc/orderly/internal/domain/user"
	"github.com/moledoc/orderly/internal/service/common/validation"
	"github.com/stretchr/testify/assert"
)

func TestValidation_ID(t *testing.T) {

	type input struct {
		Name   string
		ID     meta.ID
		Expect errwrap.Error
	}

	inputs := []input{
		{
			Name:   "empty",
			ID:     meta.EmptyID(),
			Expect: errwrap.NewError(http.StatusBadRequest, "invalid id length"),
		},
		{
			Name:   "short",
			ID:     meta.NewID()[:10],
			Expect: errwrap.NewError(http.StatusBadRequest, "invalid id length"),
		},
		{
			Name:   "long",
			ID:     meta.NewID() + meta.NewID(),
			Expect: errwrap.NewError(http.StatusBadRequest, "invalid id length"),
		},
		{
			Name:   "correct",
			ID:     meta.NewID(),
			Expect: nil,
		},
	}
	for c := 0; c < 128; c++ {
		if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
			continue
		}
		inputs = append(inputs, input{
			Name:   fmt.Sprintf("disallowed.char.%c", c),
			ID:     func() meta.ID { id := []rune(meta.NewID()); id[0] = rune(c); return meta.ID(id) }(),
			Expect: errwrap.NewError(http.StatusBadRequest, "disallowed character ('%c') in id", c),
		})
	}
	for _, in := range inputs {
		t.Run(in.Name, func(t *testing.T) {
			err := validation.ValidateID(in.ID)
			assert.Equal(t, in.Expect, err)
		})
	}
}

func TestValidation_Email(t *testing.T) {

	type input struct {
		Name   string
		Email  user.Email
		Expect errwrap.Error
	}

	inputs := []input{
		{
			Name:   "empty",
			Email:  "",
			Expect: errwrap.NewError(http.StatusBadRequest, "invalid email length"),
		},
		{
			Name:   "invalid",
			Email:  "not an email",
			Expect: errwrap.NewError(http.StatusBadRequest, "invalid email"),
		},
		// NOTE: more emails are not tested, as we're using net/mail.ParseAddress func that is from std lib - we assume it's properly tested
	}
	for _, in := range inputs {
		t.Run(in.Name, func(t *testing.T) {
			err := validation.ValidateEmail(in.Email)
			assert.Equal(t, in.Expect, err)
		})
	}
}
