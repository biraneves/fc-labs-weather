package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultPort          = "8080"
	defaultTimeout       = 5 * time.Second
	defaultAPIReturnType = "json"
	defaultWeatherAPIKey = "default_key"
)

type HTTPConfig struct {
	Addr    string
	Timeout time.Duration
}

type ViaCEPConfig struct {
	BaseURL    string
	ReturnType string
	Timeout    time.Duration
}

type WeatherAPIConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

type AppConfig struct {
	HTTP    HTTPConfig
	ViaCEP  ViaCEPConfig
	Weather WeatherAPIConfig
}

func Load(dir string) (AppConfig, error) {
	dotEnvFile := "./.env"
	if err := godotenv.Load(dotEnvFile); err != nil {
		slog.Warn("unable to read .env and relying on environment variables")
	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		slog.Warn("invalid PORT:", "default_value", defaultPort)
		appPort = defaultPort
	}

	httpTimeout, err := parseDuration(os.Getenv("HTTP_TIMEOUT"))
	if err != nil {
		slog.Warn("invalid HTTP_TIMEOUT:", "default_value", defaultTimeout)
		httpTimeout = defaultTimeout
	}

	viaCEPTimeout, err := parseDuration(os.Getenv("VIACEP_TIMEOUT"))
	if err != nil {
		slog.Warn("invalid VIACEP_TIMEOUT:", "default_value", defaultTimeout)
		viaCEPTimeout = defaultTimeout
	}

	weatherTimeout, err := parseDuration(os.Getenv("WEATHER_TIMEOUT"))
	if err != nil {
		slog.Warn("invalid WEATHER_TIMEOUT:", "default_value", defaultTimeout)
		weatherTimeout = defaultTimeout
	}

	returnType := os.Getenv("VIACEP_RETURN_TYPE")
	if returnType == "" {
		slog.Warn("invalid VIACEP_RETURN_TYPE:", "default_value", defaultAPIReturnType)
		returnType = defaultAPIReturnType
	}

	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		slog.Warn("invalid WEATHER_API_KEY - using default value")
		weatherAPIKey = defaultWeatherAPIKey
	}

	cfg := AppConfig{
		HTTP: HTTPConfig{
			Addr:    fmt.Sprintf(":%s", appPort),
			Timeout: httpTimeout,
		},
		ViaCEP: ViaCEPConfig{
			BaseURL:    strings.TrimSuffix(os.Getenv("VIACEP_URL"), "/"),
			ReturnType: returnType,
			Timeout:    viaCEPTimeout,
		},
		Weather: WeatherAPIConfig{
			BaseURL: strings.TrimSuffix(os.Getenv("WEATHER_URL"), "/"),
			APIKey:  weatherAPIKey,
			Timeout: weatherTimeout,
		},
	}

	return cfg, nil
}

func parseDuration(raw string) (time.Duration, error) {
	if raw == "" {
		return 0, errors.New("empty duration string")
	}

	return time.ParseDuration(raw)
}
