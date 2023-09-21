package package2

import (
	"fmt"
	//"os"

	myErr "amiame/test-cockroachdb-errors/error"
	"github.com/cockroachdb/errors"
)

func Func1() error {
	fmt.Println("I'm in func1")
	/*
		dir, err := os.Getwd()
		if err != nil {
			return errors.Newf("os getwd: %w", err)
		}
	*/

	//errors.Formatter
	//return errors.Wrap(myErr.Error3, "something wrong here")
	err := errors.WithSafeDetails(myErr.Error1, "%+v\n", struct {
		name string
		id   int
	}{
		name: "amir",
		id:   56,
	})
	return err
	//return errors.Wrapf(myErr.Error3, "hey I'm here %s", errors.Safe(dir))
	//return errors.NewWithDepth(1, "I don't know")
	//return errors.New("I'm a new error")
}
