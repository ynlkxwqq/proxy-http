package proxy

import "simple-http-proxy/internal/domain/proxy"

type Service struct {
	cache proxy.Cache
}

type Configuration func(s *Service) error

func New(configs ...Configuration) (*Service, error) {
	s := &Service{}

	for _, config := range configs {
		if err := config(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func WithCache(c proxy.Cache) Configuration {
	return func(s *Service) error {
		s.cache = c
		return nil
	}
}
