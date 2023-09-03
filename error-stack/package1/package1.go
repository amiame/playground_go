package package1

import (
	std_errors "errors"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
)

func Package1Function1(arg1 int) (error, error, error) {
	if err, err2, err3 := package1Function2(arg1); err != nil || err2 != nil || err3 != nil {
		return fmt.Errorf("this is an error in function 1: %w", err), errors.Wrap(err2, "this is an error in function 1"), err3
	}
	return nil, nil, nil
}

func package1Function2(arg1 int) (error, error, error) {
	if arg1 == 1 {
		return nil, nil, nil
	}
	return std_errors.New("this is an error in function 2"), errors.New("this is an error in function 2"), eris.Errorf("this is an error in function 2")
}
