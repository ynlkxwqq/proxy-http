package proxy

import "context"

type Cache interface {
	Get(ctx context.Context, id string) (requestResponsePair Entity, err error)
	Set(ctx context.Context, id string, requestResponsePair Entity)
}
