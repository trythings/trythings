package api

import (
	"errors"

	"golang.org/x/net/context"
)

type key int

const superuserKey key = 0

func AsSuperuser(ctx context.Context) context.Context {
	return context.WithValue(ctx, superuserKey, true)
}

func IsSuperuser(ctx context.Context) (bool, error) {
	v := ctx.Value(superuserKey)
	if v == nil {
		return false, nil
	}

	su, ok := v.(bool)
	if !ok {
		return false, errors.New("unexpected superuser type")
	}

	return su, nil
}
