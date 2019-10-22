package util

import (
	"errors"
	"fmt"
)

type ErrorCollector []error

func (c *ErrorCollector) Collect(err error) {
	if err == nil {
		return
	}

	*c = append(*c, err)
}

func (c *ErrorCollector) GetError() error {
	if len(*c) == 0 {
		return nil
	}

	msg := "collected errors: \n"
	for i, err := range *c {
		msg += fmt.Sprintf("\tError %d: %s\n", i, err.Error())
	}

	return errors.New(msg)
}