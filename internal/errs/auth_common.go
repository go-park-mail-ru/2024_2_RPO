package errs

import "fmt"

var ErrWrongCredentials = fmt.Errorf("wrong credentials")
var ErrBusyEmail = fmt.Errorf("this email is used in another account")
var ErrBusyNickname = fmt.Errorf("this nickname is used in another account")

const SessionCookieName string = "session_id"
