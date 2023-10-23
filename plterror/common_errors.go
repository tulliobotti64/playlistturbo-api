package plterror

import "net/http"

var ErrBadSyntax = &PLTError{message: "ERR_BAD_SYNTAX", status: http.StatusBadRequest}
var ErrServerError = &PLTError{message: "ERR_INTERNAL_SERVER_ERROR", status: http.StatusInternalServerError}
var ErrNoDataAdd = &PLTError{message: "ERR_NO_DATA_ADDED", status: http.StatusInternalServerError}
var InvalidExtension = &PLTError{message: "ERR_INVALID_EXTENSION", status: http.StatusBadRequest}
