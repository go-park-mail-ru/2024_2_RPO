package errs

import "fmt"

var ErrWrongCredentials = fmt.Errorf("Wrong credentials")
var ErrBusyEmail = fmt.Errorf("This email is used in another account")
var ErrBusyNickname = fmt.Errorf("This nickname is used in another account")

const SessionCookieName string = "session_id"
