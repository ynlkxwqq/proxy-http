package handler

import (
	"simple-http-proxy/docs"
	"simple-http-proxy/internal/config"
	"simple-http-proxy/internal/handler/http"
	"simple-http-proxy/internal/service/proxy"
	"simple-http-proxy/pkg/server/router"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Dependencies struct {
	Configs config.Configs

	ProxyService *proxy.Service
}

type Configuration func(h *Handler) error

type Handler struct {
	deps Dependencies

	HTTP *chi.Mux
}

func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	h = &Handler{
		deps: d,
	}

	for _, config := range configs {
		if err := config(h); err != nil {
			return nil, err
		}
	}

	return
}

func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		proxyHandler := http.NewProxyHandler(h.deps.ProxyService)

		h.HTTP = router.New()

		docs.SwaggerInfo.BasePath = h.deps.Configs.APP.Path
		h.HTTP.Get("/swagger/*", httpSwagger.WrapHandler)

		h.HTTP.Route("/", func(r chi.Router) {
			r.Mount("/proxy", proxyHandler.Routes())
		})

		return nil
	}
}
