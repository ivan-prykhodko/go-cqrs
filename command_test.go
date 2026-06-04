package cqrs

import (
	"errors"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func TestDispatch(t *testing.T) {
	handler := &testHandler{}
	reg := NewHandlerRegistration("testInput", handler)
	resolver := NewHandlerResolver([]HandlerRegistration{reg})
	bus := NewCommandBus(resolver)

	t.Run("Successfully dispatch command", func(t *testing.T) {
		input := testInput{ID: "cmd1"}
		result, err := Dispatch[testInput, testOutput](t.Context(), bus, input)

		a := assert.New(t)
		a.NoError(err)
		a.NotNil(result)
		a.Eq("Handled cmd1", result.Message)
	})

	t.Run("Fail to dispatch command - resolver error", func(t *testing.T) {
		input := "invalid input" // Not implementing Input interface
		// Use a type that won't be resolved
		result, err := Dispatch[string, testOutput](t.Context(), bus, input)

		a := assert.New(t)
		a.Error(err)
		a.True(errors.Is(err, ErrInvalidInput))
		a.Nil(result)
	})

	t.Run("Fail to dispatch command - handler error", func(t *testing.T) {
		handler.shouldFail = true
		defer func() {
			handler.shouldFail = false
		}()

		input := testInput{ID: "cmd2"}
		result, err := Dispatch[testInput, testOutput](t.Context(), bus, input)

		a := assert.New(t)
		a.Error(err)
		a.Eq("handler error", err.Error())
		a.Nil(result)
	})

	t.Run("Fail to dispatch command - unexpected result type", func(t *testing.T) {
		// Expecting a different result type than what handler returns
		type unexpectedResult struct{}

		input := testInput{ID: "cmd3"}
		result, err := Dispatch[testInput, unexpectedResult](t.Context(), bus, input)

		a := assert.New(t)
		a.Error(err)
		a.True(errors.Is(err, ErrUnexpectedResultType))
		a.Nil(result)
	})
}
