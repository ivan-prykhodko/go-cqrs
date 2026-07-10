# Simple CQRS in Go

A lightweight and type-safe CQRS (Command Query Responsibility Segregation) library for Go.

## Overview

This repository provides a minimalistic implementation of the CQRS pattern. It allows you to decouple your application's command and query logic using a central bus and strongly-typed handlers.

### Key Components:
- **CommandBus**: Dispatches commands that change the system state.
- **QueryBus**: Asks queries that retrieve data without modifying state.
- **HandlerResolver**: Matches inputs to their corresponding handlers.
- **Generics-based API**: Ensures type safety for both inputs and results.

## Usage

Below is a basic example of how to define and use a Command.

### 1. Define Command and Result

Your command must implement the `cqrs.Input` interface.

```go
type CreatePostCommand struct {
    Id    int
    Title string
}

func (c CreatePostCommand) ObjectType() string {
    return "blog.CreatePostCommand"
}

type CreatePostCommandResult struct {
    EventId string
}
```

### 2. Implement the Handler

```go
type createPostCommandHandler struct{}

func (h *createPostCommandHandler) Handle(ctx context.Context, input CreatePostCommand) (CreatePostCommandResult, error) {
    // Your business logic here
    return CreatePostCommandResult{EventId: "123"}, nil
}
```

### 3. Setup and Dispatch

```go
func main() {
    // Register handlers
    postHandlerRegistration := cqrs.NewHandlerRegistration(
        (CreatePostCommand{}).ObjectType(),
        &createPostCommandHandler{},
    )
    
    resolver := cqrs.NewHandlerResolver([]cqrs.HandlerRegistration{postHandlerRegistration})
    commandBus := cqrs.NewCommandBus(resolver)
    
    // Dispatch command
    result, err := cqrs.Dispatch[CreatePostCommand, CreatePostCommandResult](
        context.Background(),
        commandBus,
        CreatePostCommand{Id: 100, Title: "My New Post"},
    )
    
    if err == nil {
        fmt.Printf("Post created with EventId: %s\n", result.EventId)
    }
}
```

For a more detailed example, check the [example](example) folder.

## TODOs

- [ ] Implement middleware support (see `middleware.go`).
- [ ] Add built-in logging and instrumentation.

## License

Distributed under the MIT License. See `LICENSE` for more information.
