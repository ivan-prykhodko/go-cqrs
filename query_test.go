package cqrs

import (
	"errors"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func TestAsk(t *testing.T) {
	handler := &testHandler{}
	reg := NewHandlerRegistration("testInput", handler)
	resolver := NewHandlerResolver([]HandlerRegistration{reg})
	bus := NewQueryBus(resolver)

	t.Run("Successfully ask query", func(t *testing.T) {
		input := testInput{ID: "query1"}
		result, err := Ask[testInput, testOutput](t.Context(), bus, input)

		a := assert.New(t)
		a.NoError(err)
		a.NotNil(result)
		a.Eq("Handled query1", result.Message)
	})

	t.Run("Fail to ask query - resolver error", func(t *testing.T) {
		input := "invalid input"
		result, err := Ask[string, testOutput](t.Context(), bus, input)

		a := assert.New(t)
		a.Error(err)
		a.True(errors.Is(err, ErrInvalidInput))
		a.Eq(testOutput{}, result)
	})

	t.Run("Fail to ask query - handler error", func(t *testing.T) {
		handler.shouldFail = true
		defer func() { handler.shouldFail = false }()

		input := testInput{ID: "query2"}
		result, err := Ask[testInput, testOutput](t.Context(), bus, input)

		a := assert.New(t)
		a.Error(err)
		a.Eq("handler error", err.Error())
		a.Eq(testOutput{}, result)
	})

	t.Run("Fail to ask query - unexpected result type", func(t *testing.T) {
		type unexpectedResult struct{}

		input := testInput{ID: "query3"}
		result, err := Ask[testInput, unexpectedResult](t.Context(), bus, input)

		a := assert.New(t)
		a.Error(err)
		a.True(errors.Is(err, ErrUnexpectedResultType))
		a.Eq(unexpectedResult{}, result)
	})
}
