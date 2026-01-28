package order

import "errors"

var errTableNotFound = errors.New("table-not-found")
var errOrderNotFound = errors.New("order-not-found")
var errCourseMismatch = errors.New("course-mismatch")
var errInvalidPrintRequest = errors.New("invalid-print-request")
