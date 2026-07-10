package cqrs

import (
	"context"
	"fmt"
)

type QueryBus interface {
	resolver() HandlerResolver
}

type queryBus struct {
	r HandlerResolver
}

func NewQueryBus(r HandlerResolver) QueryBus {
	return &queryBus{
		r: r,
	}
}

func (b *queryBus) resolver() HandlerResolver {
	return b.r
}

func Ask[Q any, R any](ctx context.Context, bus QueryBus, query Q) (R, error) {
	var zero R

	handler, err := bus.resolver().Resolve(query)
	if err != nil {
		return zero, err
	}

	result, err := handler.Invoke(ctx, query)
	if err != nil {
		return zero, err
	}

	typedResult, ok := result.(R)
	if !ok {
		return zero, fmt.Errorf("%w: got %T", ErrUnexpectedResultType, result)
	}

	return typedResult, nil
}
