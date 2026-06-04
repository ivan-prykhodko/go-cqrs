package cqrs

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func TestErrors(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrUnexpectedResultType", ErrUnexpectedResultType, "handler returned unexpected result type"},
		{"ErrInvalidInput", ErrInvalidInput, "invalid input"},
		{"ErrInvalidInputType", ErrInvalidInputType, "invalid input type"},
		{"ErrHandlerNotFound", ErrHandlerNotFound, "handler not found"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)
			a.NotNil(tc.err)
			a.Eq(tc.msg, tc.err.Error())
		})
	}
}
