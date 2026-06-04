package main

import (
	"context"
	"fmt"

	"github.com/ivan-prykhodko/go-cqrs"
	"github.com/ivan-prykhodko/go-cqrs/example/application/command"
)

func main() {
	handlerRegistrations := []cqrs.HandlerRegistration{
		command.NewCreatePostCommandHandlerRegistration(),
		// add more handlers here
	}
	handlerResolver := cqrs.NewHandlerResolver(handlerRegistrations)
	commandBus := cqrs.NewCommandBus(handlerResolver)

	fmt.Println("Dispatching command...")

	result, err := cqrs.Dispatch[command.CreatePostCommand, command.CreatePostCommandResult](
		context.TODO(),
		commandBus,
		command.CreatePostCommand{
			Id:    100500,
			Title: "Sample Post",
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	fmt.Println("Done")
}
