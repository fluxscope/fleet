package handler

import "context"

type ContextExtractor func(ctx context.Context) map[string]interface{}

func NoOpExtractor() ContextExtractor {
	return func(ctx context.Context) map[string]interface{} {
		return map[string]interface{}{}
	}
}
