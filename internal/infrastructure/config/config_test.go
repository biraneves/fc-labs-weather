package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/biraneves/fc-labs-weather/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("HTTP_TIMEOUT", "")
	t.Setenv("VIACEP_URL", "")
	t.Setenv("VIACEP_TIMEOUT", "")
	t.Setenv("WEATHER_URL", "")
	t.Setenv("WEATHER_API_KEY", "")
	t.Setenv("WEATHER_TIMEOUT", "")

	tests := []struct {
		name          string
		envContent    string
		expectedError string
		assertSuccess func(t *testing.T, cfg config.AppConfig)
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
			assertSuccess: func(t *testing.T, cfg config.AppConfig) {
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
VIACEP_TIMEOUT=5s
WEATHER_URL=https://api.weatherapi.com/v1
WEATHER_TIMEOUT=5s
`,
			expectedError: "missing WEATHER_API_KEY",
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
			expectedError: "invalid HTTP_TIMEOUT",
		},
		{
			name:          "missing file",
			expectedError: "load config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.envContent != "" {
				err := os.WriteFile(filepath.Join(dir, ".env"), []byte(tt.envContent), 0o644)
				require.NoError(t, err)
			}

			cfg, err := config.Load(dir)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, config.AppConfig{}, cfg)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, tt.assertSuccess)
			tt.assertSuccess(t, cfg)
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
