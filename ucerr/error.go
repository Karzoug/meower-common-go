package ucerr

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const statusClientClosedRequest = 499

// Error is a service/usecase level error.
type Error struct {
	msg  string
	err  error
	code codes.Code
}

func NewError(err error, msg string, code codes.Code) Error {
	return Error{
		msg:  msg,
		err:  err,
		code: code,
	}
}

func NewInternalError(err error) Error {
	return Error{
		msg:  "Internal error",
		err:  err,
		code: codes.Internal,
	}
}

// Error returns error message which can be returned to the client.
func (e Error) Error() string {
	return e.msg
}

func (e Error) Code() codes.Code {
	return e.code
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) HTTPCode() int {
	switch e.code {
	case codes.Aborted, codes.AlreadyExists:
		return http.StatusConflict
	case codes.Canceled:
		return statusClientClosedRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.OK:
		return http.StatusOK
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.Unimplemented:
		return http.StatusNotImplemented
	default: // codes.Unknown, codes.Internal, codes.DataLoss
		return http.StatusInternalServerError
	}
}

func (e Error) GRPCStatus() *status.Status {
	return status.New(e.code, e.msg)
}
