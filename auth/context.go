package auth

import (
	"context"

	"github.com/rs/xid"
)

type authKey struct{}

var authUserIDKey authKey

func UserIDFromContext(ctx context.Context) xid.ID {
	username, ok := ctx.Value(authUserIDKey).(xid.ID)
	if !ok {
		return xid.NilID()
	}
	return username
}

func WithUserID(ctx context.Context, id xid.ID) context.Context {
	return context.WithValue(ctx, authUserIDKey, id)
}
