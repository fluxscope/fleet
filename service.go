package fleet

import "context"

type Service interface {
	ID() string
	Run(context.Context) error
	Shutdown(context.Context) error
}
