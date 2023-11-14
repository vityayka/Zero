package context

import (
	"context"

	"github.com/vityayka/go-zero/models"
)

type key string

const (
	ctxKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, ctxKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(ctxKey)
	user, isOk := val.(*models.User)
	if !isOk {
		return nil
	}
	return user
}
