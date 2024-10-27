package boards

import "fmt"

var ErrNotFound = fmt.Errorf("Not found")
var ErrNotPermitted = fmt.Errorf("Not permitted")

type BoardPutRequest struct {
	NewName        string `json:"name"`
	NewDescription string `json:"description"`
}
