package errs

import "fmt"

var ErrNotFound = fmt.Errorf("not found")
var ErrNotPermitted = fmt.Errorf("not permitted")
var ErrAlreadyExists = fmt.Errorf("already exists")
var ErrValidation = fmt.Errorf("validation error")
