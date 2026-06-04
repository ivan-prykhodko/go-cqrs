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

func Dispatch[C any, R any](ctx context.Context, bus CommandBus, command C) (*R, error) {
	handler, err := bus.resolver().Resolve(command)
	if err != nil {
		return nil, err
	}

	result, err := handler.Invoke(ctx, command)
	if err != nil {
		return nil, err
	}

	typedResult, ok := result.(*R)
	if !ok {
		return nil, fmt.Errorf("%w: got %T", ErrUnexpectedResultType, result)
	}

	return typedResult, nil
}
