package command

import (
	"context"
	"time"

	"github.com/ivan-prykhodko/go-cqrs"
)

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

func NewCreatePostCommandHandlerRegistration() cqrs.HandlerRegistration {
	return cqrs.NewHandlerRegistration(
		(CreatePostCommand{}).ObjectType(),
		newCreatePostCommandHandler(),
	)
}

type createPostCommandHandler struct {
}

func newCreatePostCommandHandler() cqrs.Handler[CreatePostCommand, CreatePostCommandResult] {
	return &createPostCommandHandler{}
}

func (h *createPostCommandHandler) Handle(ctx context.Context, input CreatePostCommand) (*CreatePostCommandResult, error) {
	// do whatever you need here
	time.Sleep(5 * time.Second)

	return &CreatePostCommandResult{EventId: "123"}, nil
}
