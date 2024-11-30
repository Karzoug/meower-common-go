package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/rs/xid"

	"github.com/Karzoug/meower-common-go/auth"
)

const userKey string = "x-user-id"

func Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		key, found := md[userKey]
		if found && len(key) != 0 {
			id, err := xid.FromString(key[0])
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, "invalid user id")
			}
			ctx = auth.WithUserID(ctx, id)
		}

		return handler(ctx, req)
	}
}
