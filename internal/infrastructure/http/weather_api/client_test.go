package weatherapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
	weatherapi "github.com/biraneves/fc-labs-weather/internal/infrastructure/http/weather_api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPClient(t *testing.T) {
	client := weatherapi.NewHTTPClient(nil, "", "token", 0)
	require.NotNil(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.FetchCurrent(ctx, dto.WeatherAPIRequestDto{Q: "São Paulo"})
	assert.Error(t, err)
}

func TestHTTPClient_FetchCurrent(t *testing.T) {
	type fields struct {
		status int
		body   string
		delay  time.Duration
	}

	tests := []struct {
		name          string
		fields        fields
		request       dto.WeatherAPIRequestDto
		want          dto.WeatherAPIResponseDto
		expectedError string
	}{
		{
			name: "success",
			fields: fields{
				status: http.StatusOK,
				body: `{
				  "location": {
				    "name": "Sao Paulo"
				  },
				  "current": {
				    "temp_C": 22.5
				  }
				}`,
			},
			request: dto.WeatherAPIRequestDto{Q: "São Paulo"},
			want: func() dto.WeatherAPIResponseDto {
				var resp dto.WeatherAPIResponseDto
				resp.Location.Name = "Sao Paulo"
				resp.Current.TempC = 22.5
				return resp
			}(),
		},
		{
			name: "unexpected status",
			fields: fields{
				status: http.StatusBadRequest,
				body:   `{"error": {"message": "invalid"}}`,
			},
			request:       dto.WeatherAPIRequestDto{Q: "São Paulo"},
			expectedError: "weatherapi: unexpected status: 400",
		},
		{
			name: "invalid json",
			fields: fields{
				status: http.StatusOK,
				body:   `{"location":`,
			},
			request:       dto.WeatherAPIRequestDto{Q: "São Paulo"},
			expectedError: "weatherapi: decode response",
		},
		{
			name: "context timeout",
			fields: fields{
				status: http.StatusOK,
				body:   `{}`,
				delay:  100 * time.Millisecond,
			},
			request:       dto.WeatherAPIRequestDto{Q: "São Paulo"},
			expectedError: "context deadline exceeded",
		},
		{
			name:          "empty query",
			request:       dto.WeatherAPIRequestDto{Q: ""},
			expectedError: weatherapi.ErrEmptyQuery.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.fields.delay > 0 {
					time.Sleep(tt.fields.delay)
				}

				if tt.fields.status != 0 {
					w.WriteHeader(tt.fields.status)
				}

				if tt.fields.body != "" {
					_, _ = w.Write([]byte(tt.fields.body))
				}
			}))
			defer server.Close()

			client := weatherapi.NewHTTPClient(nil, server.URL, "token", 50*time.Millisecond)

			ctx := context.Background()
			if tt.fields.delay > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 50*time.Millisecond)
				defer cancel()
			}

			got, err := client.FetchCurrent(ctx, tt.request)
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestHTTPClient_FetchCurrent_MissingAPIKey(t *testing.T) {
	client := weatherapi.NewHTTPClient(nil, "https://example.com", "", time.Second)

	_, err := client.FetchCurrent(context.Background(), dto.WeatherAPIRequestDto{Q: "São Paulo"})
	require.Error(t, err)
	assert.ErrorIs(t, err, weatherapi.ErrMissingAPIKey)
}
