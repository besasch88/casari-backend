package printer

import "errors"

var errPrinterNotFound = errors.New("printer-not-found")
var errPrinterSameTitleAlreadyExists = errors.New("printer-same-title-already-exists")
