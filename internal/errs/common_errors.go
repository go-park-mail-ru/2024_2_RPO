package errs

import "fmt"

var ErrNotFound = fmt.Errorf("Not found")
var ErrNotPermitted = fmt.Errorf("Not permitted")
var ErrAlreadyExists = fmt.Errorf("Already exists")
