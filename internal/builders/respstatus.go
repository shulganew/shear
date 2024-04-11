package builders

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// Stract for universal responce.
type RespStatus struct {
	respStatusREST int
}

// Constructor for RespStatus
func NewRespStatus(errStatusREST int) RespStatus {
	o := new(RespStatus)
	o.respStatusREST = errStatusREST
	return *o
}

// Get status error for REST.
func (e RespStatus) GetStatusREST() int {
	return e.respStatusREST
}

// Get status error for gRPC.
func (e RespStatus) GetStatusGRPC() codes.Code {
	switch e.respStatusREST {
	// 200
	case http.StatusOK:
		return codes.OK
	// 201
	case http.StatusCreated:
		return codes.OK
	// 400
	case http.StatusBadRequest:
		return codes.OutOfRange
	// 401
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	// 403
	case http.StatusForbidden:
		return codes.PermissionDenied
	//404
	case http.StatusNotFound:
		return codes.NotFound
	// 409
	case http.StatusConflict:
		return codes.AlreadyExists
	// 410
	case http.StatusGone:
		return codes.Unknown
	// 500
	case http.StatusInternalServerError:
		return codes.Internal
	// 503
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	default:
		return codes.Unimplemented
	}

}

// Setter method for the field respStatusREST of type int in the object RespStatus.
func (e *RespStatus) SetStatusREST(status int) {
	e.respStatusREST = status
}
