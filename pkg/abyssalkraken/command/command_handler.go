package command

import "context"

type CommandHandler[C Command[R], R any] interface {
	Handle(ctx context.Context, command C) (R, error)
}
