package package1

import (
	"fmt"

	"amiame/test-cockroachdb-errors/server/package2"
	"github.com/cockroachdb/errors"
)

func Func1() error {
	fmt.Println("I'm in func1")

	err := package2.Func1()
	if err != nil {
		err = errors.Wrap(errors.WithTelemetry(err, "amir"), "wrapping error")
	}
	return err
}
