package table

import "errors"

var errTableNotFound = errors.New("table-not-found")
var errTableSameNameAlreadyExists = errors.New("table-same-name-already-exists")
