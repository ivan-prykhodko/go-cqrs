package cqrs

import (
	"context"
	"errors"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

// Mock objects for testing
type testInput struct {
	ID string
}

func (i testInput) ObjectType() string {
	return "testInput"
}

type testOutput struct {
	Message string
}

type testHandler struct {
	shouldFail bool
}

func (h *testHandler) Handle(ctx context.Context, input testInput) (testOutput, error) {
	if h.shouldFail {
		var zero testOutput
		return zero, errors.New("handler error")
	}
	return testOutput{Message: "Handled " + input.ID}, nil
}

func TestHandlerResolver(t *testing.T) {
	t.Run("Successfully resolve handler", func(t *testing.T) {
		handler := &testHandler{}
		reg := NewHandlerRegistration("testInput", handler)
		resolver := NewHandlerResolver([]HandlerRegistration{reg})

		input := testInput{ID: "123"}
		resolvedReg, err := resolver.Resolve(input)

		a := assert.New(t)
		a.NoError(err)
		a.NotNil(resolvedReg)
		a.Eq("testInput", resolvedReg.ObjectType())
		a.Eq(reg, resolvedReg)
	})

	t.Run("Fail to resolve handler - not found", func(t *testing.T) {
		resolver := NewHandlerResolver([]HandlerRegistration{})

		input := testInput{ID: "123"}
		resolvedReg, err := resolver.Resolve(input)

		a := assert.New(t)
		a.Error(err)
		a.True(errors.Is(err, ErrHandlerNotFound))
		a.Nil(resolvedReg)
	})

	t.Run("Fail to resolve handler - invalid input", func(t *testing.T) {
		resolver := NewHandlerResolver([]HandlerRegistration{})

		// something that doesn't implement Input
		errInput := "not an input"
		resolvedReg, err := resolver.Resolve(errInput)

		a := assert.New(t)
		a.Error(err)
		a.True(errors.Is(err, ErrInvalidInput))
		a.Nil(resolvedReg)
	})
}

func TestHandlerRegistration(t *testing.T) {
	t.Run("Successfully invoke handler", func(t *testing.T) {
		handler := &testHandler{}
		reg := NewHandlerRegistration("testInput", handler)

		input := testInput{ID: "123"}
		result, err := reg.Invoke(t.Context(), input)

		a := assert.New(t)
		a.NoError(err)
		a.NotNil(result)

		output, ok := result.(testOutput)
		a.True(ok)
		a.Eq("Handled 123", output.Message)
	})

	t.Run("Fail to invoke handler - wrong input type", func(t *testing.T) {
		handler := &testHandler{}
		reg := NewHandlerRegistration("testInput", handler)

		// Different input type than expected by reg.handleFunc
		type otherInput struct{ testInput }
		result, err := reg.Invoke(t.Context(), otherInput{})

		a := assert.New(t)
		a.Error(err)
		a.True(errors.Is(err, ErrInvalidInputType))
		a.Nil(result)
	})

	t.Run("Fail to invoke handler - handler returns error", func(t *testing.T) {
		handler := &testHandler{shouldFail: true}
		reg := NewHandlerRegistration("testInput", handler)

		input := testInput{ID: "123"}
		result, err := reg.Invoke(t.Context(), input)

		a := assert.New(t)
		a.Error(err)
		a.Eq("handler error", err.Error())
		a.Eq(testOutput{}, result)
	})
}
