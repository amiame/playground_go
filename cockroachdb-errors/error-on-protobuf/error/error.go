package my_error

import (
	"github.com/cockroachdb/errors"
	//"errors"
)

var (
	Error1 = errors.New("error1")
	Error2 = errors.New("error2")
	Error3 = errors.New("error3")
)
