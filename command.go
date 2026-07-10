package cqrs

import (
	"context"
	"fmt"
)

type CommandBus interface {
	resolver() HandlerResolver
}

type commandBus struct {
	r HandlerResolver
}

func NewCommandBus(r HandlerResolver) CommandBus {
	return &commandBus{
		r: r,
	}
}

func (b *commandBus) resolver() HandlerResolver {
	return b.r
}

func Dispatch[C any, R any](ctx context.Context, bus CommandBus, command C) (R, error) {
	var zero R

	handler, err := bus.resolver().Resolve(command)
	if err != nil {
		return zero, err
	}

	result, err := handler.Invoke(ctx, command)
	if err != nil {
		return zero, err
	}

	typedResult, ok := result.(R)
	if !ok {
		return zero, fmt.Errorf("%w: got %T", ErrUnexpectedResultType, result)
	}

	return typedResult, nil
}
