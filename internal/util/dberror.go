package util

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicate = errors.New("duplicate entry")
