package server

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
)

type Server struct {
	http *http.Server
}

type Configuration func(s *Server) error

func New(configs ...Configuration) (s *Server, err error) {
	s = &Server{}

	for _, config := range configs {
		if err := config(s); err != nil {
			return nil, err
		}
	}

	return
}

func (s *Server) Run(logger *zerolog.Logger) (err error) {
	if s.http != nil {
		go func() {
			err := s.http.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				logger.Error().Err(err).Msg("ERR_SERVE_HTTP")
				return
			}
		}()
	}

	return
}

func (s *Server) Stop(ctx context.Context) (err error) {
	if s.http != nil {
		return s.http.Shutdown(ctx)
	}

	return
}

func WithHTTPServer(h http.Handler, port string) Configuration {
	return func(s *Server) (err error) {
		s.http = &http.Server{
			Addr:    ":" + port,
			Handler: h,
		}
		return nil
	}
}
