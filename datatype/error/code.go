package error

import "net/http"

var (
	ServerError  = NewError(http.StatusInternalServerError, 1000, "Server Error", nil)
	BadRequest   = NewError(http.StatusBadRequest, 1001, "Bad Request", nil)
	Forbidden    = NewError(http.StatusForbidden, 1002, "Permission Forbidden", nil)
	Unauthorized = NewError(http.StatusUnauthorized, 1003, "Request is unauthorized", nil)
)
