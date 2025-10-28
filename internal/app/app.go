package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"simple-http-proxy/internal/cache"
	"simple-http-proxy/internal/config"
	"simple-http-proxy/internal/handler"
	"simple-http-proxy/internal/service/proxy"
	"simple-http-proxy/pkg/log"
	"simple-http-proxy/pkg/server"
)

func Run() {
	logger := log.LoggerFromContext(context.Background())

	config, err := config.New()
	if err != nil {
		logger.Error().Err(err).Msg("ERR_LOAD_CONFIG")
	}

	caches, err := cache.New(
		cache.Dependencies{},
		cache.WithMemoryCache(),
	)
	if err != nil {
		logger.Error().Err(err).Msg("ERR_CREATE_CACHES")
	}

	proxyService, err := proxy.New(
		proxy.WithCache(caches.Proxy),
	)
	if err != nil {
		logger.Error().Err(err).Msg("ERR_CREATE_PROXY_SERVICE")
	}

	handlers, err := handler.New(handler.Dependencies{
		ProxyService: proxyService,
	}, handler.WithHTTPHandler())
	if err != nil {
		logger.Error().Err(err).Msg("ERR_CREATE_HANDLERS")
	}

	server, err := server.New(server.WithHTTPServer(handlers.HTTP, config.APP.Port))
	if err != nil {
		logger.Error().Err(err).Msg("ERR_CREATE_SERVER")
	}

	if err = server.Run(logger); err != nil {
		logger.Error().Err(err).Msg("ERR_RUN_SERVER")
		return
	}
	logger.Info().Msg("http server started on http://localhost:" + config.APP.Port + ", swagger is at /swagger/index.html")

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err = server.Stop(ctx); err != nil {
		panic(err)
	}

	fmt.Println("running cleanup tasks...")

	fmt.Println("server was successful shutdown.")
}
