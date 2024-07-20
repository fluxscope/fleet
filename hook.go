package fleet

import "context"

type Hook func(ctx context.Context) error
