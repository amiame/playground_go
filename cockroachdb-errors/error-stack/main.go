package main

import (
	"fmt"
	"github.com/cockroachdb/errors"
)

type User struct {
	name string
}

func func1() error {
	user := User{
		name: "amir",
	}
	return errors.Newf("this is a new error: %+v", user)
}

func func2() error {
	if err := func1(); err != nil {
		return errors.Wrapf(err, "this is a wrapping error")
	}

	return nil
}

func func3() error {
	if err := func2(); err != nil {
		return errors.Wrapf(err, "this is a wrapping error")
	}

	return nil
}

func main() {
	if err := func3(); err != nil {
		fmt.Printf("error: %+v\n", err)
		fmt.Println("\n=================================")
		fmt.Printf("Unwrap(error): %+v\n", errors.Unwrap(err))
		fmt.Println("\n=================================")
		fmt.Printf("UnwrapAll(error): %+v\n", errors.UnwrapAll(err))
	}
}
