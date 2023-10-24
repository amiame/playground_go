package main

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"zerolog-error-stack/pkgerrors"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	out := &bytes.Buffer{}
	log := zerolog.New(out)

	err := errors.Wrap(errors.New("error message"), "from error")
	log.Log().Stack().Err(err).Msg("")

	fmt.Println(out.String())
}
