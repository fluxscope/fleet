package fleet

import "context"

type Hook func(ctx context.Context) error

//type Command func(ctx context.Context, args ...string) error

type Command interface {
	Run(context.Context, ...string) error
}
