package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/biraneves/fc-labs-weather/internal/application/usecase"
	"github.com/biraneves/fc-labs-weather/internal/infrastructure/config"
	"github.com/biraneves/fc-labs-weather/internal/infrastructure/http/server"
	viacep "github.com/biraneves/fc-labs-weather/internal/infrastructure/http/viacep"
	weatherapi "github.com/biraneves/fc-labs-weather/internal/infrastructure/http/weather_api"
)

const shutdownTimeout = 10 * time.Second

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	httpClient := &http.Client{Timeout: cfg.HTTP.Timeout}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	zipcodeClient := viacep.NewHTTPClient(httpClient, cfg.ViaCEP.BaseURL, cfg.ViaCEP.Timeout, logger)
	weatherClient := weatherapi.NewHTTPClient(httpClient, cfg.Weather.BaseURL, cfg.Weather.APIKey, cfg.Weather.Timeout, logger)

	getWeatherUC := usecase.NewGetWeatherByCEPUseCase(zipcodeClient, weatherClient)

	handler := server.NewHandler(getWeatherUC, logger)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	logging := server.NewLoggerMiddleware(logger)
	recovery := server.NewRecoveryMiddleware(logger)
	rootHandler := logging.Wrap(recovery.Wrap(mux))

	srv := &http.Server{
		Addr:    cfg.HTTP.Addr,
		Handler: rootHandler,
	}

	go func() {
		slog.Info("server listening:", "port", cfg.HTTP.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	slog.Info("shutdown requested")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed:", "error", err.Error())
		return
	}

	slog.Info("server stopped")
}
