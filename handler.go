package cqrs

import (
	"context"
	"fmt"
)

type Input interface {
	ObjectType() string
}

type Handler[I any, R any] interface {
	Handle(context.Context, I) (R, error)
}

type HandlerFunc func(context.Context, any) (any, error)

type HandlerResolver interface {
	Resolve(input any) (HandlerRegistration, error)
}

type handlerResolver struct {
	handlers map[string]HandlerRegistration
}

func NewHandlerResolver(registrations []HandlerRegistration) HandlerResolver {
	handlers := make(map[string]HandlerRegistration)
	for _, registration := range registrations {
		handlers[registration.ObjectType()] = registration
	}

	return &handlerResolver{
		handlers: handlers,
	}
}

func (r *handlerResolver) Resolve(input any) (HandlerRegistration, error) {
	i, ok := (input).(Input)
	if !ok {
		return nil, fmt.Errorf("%w: input must implement Input interface, got %T", ErrInvalidInput, input)
	}

	key := i.ObjectType()
	handler, ok := r.handlers[key]
	if !ok {
		return nil, fmt.Errorf("%w: input %T", ErrHandlerNotFound, input)
	}

	return handler, nil
}

type HandlerRegistration interface {
	ObjectType() string
	Invoke(ctx context.Context, i any) (any, error)
}

type handlerRegistration struct {
	objectType string
	handleFunc HandlerFunc
}

func NewHandlerRegistration[I any, R any](objectType string, handler Handler[I, R]) HandlerRegistration {
	return &handlerRegistration{
		objectType: objectType,
		handleFunc: func(ctx context.Context, i any) (any, error) {
			input, ok := i.(I)
			if !ok {
				return nil, fmt.Errorf("%w: expected %T, got %T", ErrInvalidInputType, *new(I), i)
			}

			return handler.Handle(ctx, input)
		},
	}
}

func (r *handlerRegistration) ObjectType() string {
	return r.objectType
}

func (r *handlerRegistration) Invoke(ctx context.Context, i any) (any, error) {
	return r.handleFunc(ctx, i)
}
