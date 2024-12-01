package interceptor

import (
	"context"
	"errors"
	"syscall"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Karzoug/meower-common-go/ucerr"
)

type grpcStatus interface {
	GRPCStatus() *status.Status
}

func Error(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if nil == err {
				return
			}

			switch e := err.(type) { //nolint:errorlint
			case grpcStatus:
				logOnlyServiceError(ctx, err, info.FullMethod, logger)
				err = e.GRPCStatus().Err()
			default:
				switch {
				//  got network error -> it's ok, so ignore it
				case isNetworkError(e):

				// it's unknown (= untrusted) error,
				// log it and return internal server error
				default:
					logger.Error().
						Err(e).
						Ctx(ctx). // for trace_id
						Str("method", info.FullMethod).
						Msg("error handler")
					err = status.Error(codes.Internal, codes.Internal.String())
				}
			}
		}()

		return handler(ctx, req)
	}
}

func logOnlyServiceError(ctx context.Context, err error, method string, logger zerolog.Logger) {
	var e ucerr.Error
	if !errors.As(err, &e) {
		return
	}

	var ev *zerolog.Event
	switch e.Code() {
	case codes.OK,
		codes.Canceled,
		codes.InvalidArgument,
		codes.DeadlineExceeded,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.FailedPrecondition,
		codes.OutOfRange,
		codes.Unimplemented, codes.Unauthenticated:
		return

	case codes.ResourceExhausted, codes.Aborted:
		ev = logger.Warn()

	case codes.Internal, codes.Unavailable, codes.Unknown, codes.DataLoss:
		ev = logger.Error()
	}

	ev.Err(e.Unwrap()).
		Ctx(ctx). // for trace_id
		Str("method", method).
		Msg("error handler")
}

func isNetworkError(err error) bool {
	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return true

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return true
	}

	return false
}
