package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
)

func Ptr[T any](t T) *T {
	return &t
}

func Deref[T any](t *T) T {
	if t == nil {
		var tt T
		return tt
	}
	return *t
}

func RePtr[T any](t *T) *T {
	return Ptr(Deref(t))
}

const (
	RandAlphanumLen = 32
)

func RandAlphanum() string {
	v := ""
	for len(v) < RandAlphanumLen {
		v = fmt.Sprintf("%v%v", v, strconv.FormatInt(rand.Int63(), 16))
	}
	v = v[:32]
	return v
}

func IsZeroValue(x any) bool {
	if x == nil {
		return true
	}
	v := reflect.ValueOf(x)
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}
