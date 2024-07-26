package runx

import "context"

func RunForever(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
