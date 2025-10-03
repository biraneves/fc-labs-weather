package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
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
	v := viper.New()

	v.SetConfigFile(filepath.Join(dir, ".env"))
	v.SetConfigType("env")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("HTTP_ADDR", ":8080")
	v.SetDefault("HTTP_TIMEOUT", "5s")
	v.SetDefault("VIACEP_URL", "https://viacep.com.br/ws")
	v.SetDefault("VIACEP_RETURN_TYPE", "json")
	v.SetDefault("VIACEP_TIMEOUT", "5s")
	v.SetDefault("WEATHER_URL", "https://api.weatherapi.com/v1")
	v.SetDefault("WEATHER_TIMEOU", "5s")

	if err := v.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return AppConfig{}, fmt.Errorf("load config file: %w", err)
	}

	httpTimeout, err := parseDuration(v.GetString("HTTP_TIMEOUT"))
	if err != nil {
		return AppConfig{}, fmt.Errorf("invalid HTTP_TIMEOUT: %w", err)
	}

	viaCEPTimeout, err := parseDuration(v.GetString("VIACEP_TIMEOUT"))
	if err != nil {
		return AppConfig{}, fmt.Errorf("invalid VIACEP_TIMEOUT: %w", err)
	}

	weatherTimeout, err := parseDuration(v.GetString("WEATHER_TIMEOUT"))
	if err != nil {
		return AppConfig{}, fmt.Errorf("invalid WEATHER_TIMEOUT: %w", err)
	}

	apiKey := strings.TrimSpace(v.GetString("WEATHER_API_KEY"))
	if apiKey == "" {
		return AppConfig{}, errors.New("missing WEATHER_API_KEY")
	}

	cfg := AppConfig{
		HTTP: HTTPConfig{
			Addr:    v.GetString("HTTP_ADDR"),
			Timeout: httpTimeout,
		},
		ViaCEP: ViaCEPConfig{
			BaseURL:    strings.TrimSuffix(v.GetString("VIACEP_URL"), "/"),
			ReturnType: v.GetString("VIACEP_RETURN_TYPE"),
			Timeout:    viaCEPTimeout,
		},
		Weather: WeatherAPIConfig{
			BaseURL: strings.TrimSuffix(v.GetString("WEATHER_URL"), "/"),
			APIKey:  apiKey,
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
