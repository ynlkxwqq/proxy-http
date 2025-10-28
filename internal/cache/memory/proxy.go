package memory

import (
	"context"
	"sync"

	"simple-http-proxy/internal/domain/proxy"
	"simple-http-proxy/pkg/log"
)

type ProxyCache struct {
	c sync.Map
}

func NewProxyCache() *ProxyCache {
	return &ProxyCache{}
}

func (p *ProxyCache) Get(ctx context.Context, id string) (pair proxy.Entity, err error) {
	logger := log.LoggerFromContext(ctx).With().Str("ProxyCache", id).Logger()

	if val, ok := p.c.Load(id); ok {
		logger.Info().Msg("request-response pair retrieved from cache")
		return val.(proxy.Entity), nil
	}

	logger.Warn().Msg("request-response pair not found in cache")

	return
}

func (p *ProxyCache) Set(ctx context.Context, id string, requestResponsePair proxy.Entity) {
	p.c.Store(id, requestResponsePair)
	logger := log.LoggerFromContext(ctx).With().Str("ProxyCache", id).Logger()

	logger.Info().Msg("request-response pair added to cache")
}
