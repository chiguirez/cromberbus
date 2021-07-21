package multierror

import (
	"errors"
	"fmt"
)

type MultiErr interface {
	error
	Append(errors ...error)
	NilOrError() error
	Errors() []error
}

type multiErr struct {
	err []error
}

func New(errorList ...error) MultiErr {
	m := &multiErr{}

	for _, err := range errorList {
		if errors.As(err, new(MultiErr)) {
			m.Append(err.(MultiErr).Errors()...) //nolint:errorlint

			continue
		}

		if err != nil {
			m.err = append(m.err, err)
		}
	}

	return m
}

func (m *multiErr) Append(errors ...error) {
	for _, err := range errors {
		if err != nil {
			m.err = append(m.err, err)
		}
	}
}

func (m *multiErr) Error() string {
	errString := "\nmultiple errors happened : \n"
	for _, err := range m.err {
		errString += fmt.Sprintf("	%v \n", err)
	}

	return errString
}

func (m *multiErr) NilOrError() error {
	switch len(m.err) {
	case 0:
		return nil
	case 1:
		return m.err[0]
	default:
		return m
	}
}

func (m *multiErr) Errors() []error {
	return m.err
}

func (m *multiErr) Is(target error) bool {
	for _, err := range m.err {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}
