package menuOption

import "errors"

var errMenuItemNotFound = errors.New("menu-item-not-found")
var errMenuOptionNotFound = errors.New("menu-option-not-found")
var errMenuOptionSameTitleAlreadyExists = errors.New("menu-option-same-title-already-exists")
