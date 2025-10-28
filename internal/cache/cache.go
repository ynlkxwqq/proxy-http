package cache

import (
	"simple-http-proxy/internal/cache/memory"
	"simple-http-proxy/internal/domain/proxy"
)

type Dependencies struct {
}

type Cache struct {
	deps Dependencies

	Proxy proxy.Cache
}

type Configuration func(*Cache) error

func New(d Dependencies, configs ...Configuration) (c *Cache, err error) {
	c = &Cache{
		deps: d,
	}

	for _, config := range configs {
		if err = config(c); err != nil {
			return
		}
	}

	return
}

func WithMemoryCache() Configuration {
	return func(c *Cache) error {
		c.Proxy = memory.NewProxyCache()
		return nil
	}
}
