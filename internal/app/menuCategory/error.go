package menuCategory

import "errors"

var errPrinterNotFound = errors.New("printer-not-found")
var errMenuCategoryNotFound = errors.New("menu-category-not-found")
var errMenuCategorySameTitleAlreadyExists = errors.New("menu-category-same-title-already-exists")
