package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/biraneves/fc-labs-weather/internal/infrastructure/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	resetEnv := func(t *testing.T) {
		t.Setenv("PORT", "")
		t.Setenv("HTTP_TIMEOUT", "")
		t.Setenv("VIACEP_URL", "")
		t.Setenv("VIACEP_TIMEOUT", "")
		t.Setenv("VIACEP_RETURN_TYPE", "")
		t.Setenv("WEATHER_URL", "")
		t.Setenv("WEATHER_API_KEY", "")
		t.Setenv("WEATHER_TIMEOUT", "")
	}

	tests := []struct {
		name       string
		envContent string
		assertions func(t *testing.T, cfg config.AppConfig)
	}{
		{
			name: "success",
			envContent: `PORT=9090
HTTP_TIMEOUT=3s
VIACEP_URL=https://viacep.com.br/ws/
VIACEP_RETURN_TYPE=json
VIACEP_TIMEOUT=4s
WEATHER_URL=https://api.weatherapi.com/v1/
WEATHER_API_KEY=abc123
WEATHER_TIMEOUT=6s
`,
			assertions: func(t *testing.T, cfg config.AppConfig) {
				assert.Equal(t, ":9090", cfg.HTTP.Addr)
				assert.Equal(t, 3*time.Second, cfg.HTTP.Timeout)

				assert.Equal(t, "https://viacep.com.br/ws", cfg.ViaCEP.BaseURL)
				assert.Equal(t, "json", cfg.ViaCEP.ReturnType)
				assert.Equal(t, 4*time.Second, cfg.ViaCEP.Timeout)

				assert.Equal(t, "https://api.weatherapi.com/v1", cfg.Weather.BaseURL)
				assert.Equal(t, "abc123", cfg.Weather.APIKey)
				assert.Equal(t, 6*time.Second, cfg.Weather.Timeout)
			},
		},
		{
			name: "missing api key",
			envContent: `PORT=8080
HTTP_TIMEOUT=5s
VIACEP_URL=https://viacep.com.br/ws/
VIACEP_RETURN_TYPE=xml
VIACEP_TIMEOUT=5s
WEATHER_URL=https://api.weatherapi.com/v1
WEATHER_TIMEOUT=5s
`,
			assertions: func(t *testing.T, cfg config.AppConfig) {
				assert.Equal(t, ":8080", cfg.HTTP.Addr)
				assert.Equal(t, 5*time.Second, cfg.HTTP.Timeout)

				assert.Equal(t, "https://viacep.com.br/ws", cfg.ViaCEP.BaseURL)
				assert.Equal(t, "xml", cfg.ViaCEP.ReturnType)
				assert.Equal(t, 5*time.Second, cfg.ViaCEP.Timeout)

				assert.Equal(t, "https://api.weatherapi.com/v1", cfg.Weather.BaseURL)
				assert.Equal(t, "default_key", cfg.Weather.APIKey)
				assert.Equal(t, 5*time.Second, cfg.Weather.Timeout)
			},
		},
		{
			name: "invalid http timeout",
			envContent: `PORT=8080
HTTP_TIMEOUT=invalid
VIACEP_URL=https://viacep.com.br/ws/
VIACEP_TIMEOUT=5s
WEATHER_URL=https://api.weatherapi.com/v1
WEATHER_API_KEY=foo
WEATHER_TIMEOUT=5s
`,
			assertions: func(t *testing.T, cfg config.AppConfig) {
				assert.Equal(t, ":8080", cfg.HTTP.Addr)
				assert.Equal(t, 5*time.Second, cfg.HTTP.Timeout)
				assert.Equal(t, 5*time.Second, cfg.ViaCEP.Timeout)
				assert.Equal(t, 5*time.Second, cfg.Weather.Timeout)
			},
		},
		{
			name: "missing file",
			assertions: func(t *testing.T, cfg config.AppConfig) {
				assert.Equal(t, ":8080", cfg.HTTP.Addr)
				assert.Equal(t, 5*time.Second, cfg.HTTP.Timeout)

				assert.Equal(t, "", cfg.ViaCEP.BaseURL)
				assert.Equal(t, "json", cfg.ViaCEP.ReturnType)
				assert.Equal(t, 5*time.Second, cfg.ViaCEP.Timeout)

				assert.Equal(t, "", cfg.Weather.BaseURL)
				assert.Equal(t, "default_key", cfg.Weather.APIKey)
				assert.Equal(t, 5*time.Second, cfg.Weather.Timeout)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetEnv(t)

			if tt.envContent != "" {
				dir := t.TempDir()
				envPath := filepath.Join(dir, ".env")

				require.NoError(t, os.WriteFile(envPath, []byte(tt.envContent), 0o644))
				require.NoError(t, godotenv.Overload(envPath))
			}

			cfg, err := config.Load(".")
			require.NoError(t, err)
			tt.assertions(t, cfg)
		})
	}
}

func TestParseDuration(t *testing.T) {
	t.Run("valid duration", func(t *testing.T) {
		got, err := config.ParseDuration("2s")
		require.NoError(t, err)
		assert.Equal(t, 2*time.Second, got)
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := config.ParseDuration("")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "empty duration string")
	})

	t.Run("invalid format", func(t *testing.T) {
		_, err := config.ParseDuration("abc")
		require.Error(t, err)
	})
}
