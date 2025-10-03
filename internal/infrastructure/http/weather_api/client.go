package weatherapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
	"github.com/biraneves/fc-labs-weather/internal/application/ports/outbound"
	"github.com/biraneves/fc-labs-weather/internal/infrastructure/http/server"
)

var (
	ErrMissingAPIKey = errors.New("weatherapi: missing api key")
	ErrEmptyQuery    = errors.New("weatherapi: empty query parameter")
)

type HTTPClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	timeout    time.Duration
	logger     *slog.Logger
}

func NewHTTPClient(httpClient *http.Client, baseURL, apiKey string, timeout time.Duration, logger *slog.Logger) *HTTPClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if baseURL == "" {
		baseURL = "https://api.weatherapi.com/v1"
	}

	if timeout == 0 {
		timeout = 5 * time.Second
	}

	return &HTTPClient{
		httpClient: httpClient,
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		apiKey:     apiKey,
		timeout:    timeout,
		logger:     logger,
	}
}

func (h *HTTPClient) FetchCurrent(ctx context.Context, request dto.WeatherAPIRequestDto) (dto.WeatherAPIResponseDto, error) {
	logger := server.LoggerFromContext(ctx, h.logger)

	if h.apiKey == "" {
		logger.Error("weatherapi: missing api key",
			slog.String("type", "outbound_error"),
			slog.String("query", request.Q),
		)
		return dto.WeatherAPIResponseDto{}, ErrMissingAPIKey
	}

	query := strings.TrimSpace(request.Q)
	if query == "" {
		logger.Warn("weatherapi: empty query parameter",
			slog.String("type", "outbound_error"),
		)
		return dto.WeatherAPIResponseDto{}, ErrEmptyQuery
	}

	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/current.json", h.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		logger.Error("weatherapi: create request failed",
			slog.String("type", "outbound_error"),
			slog.String("query", query),
			slog.String("error", err.Error()),
		)
		return dto.WeatherAPIResponseDto{}, fmt.Errorf("weatherapi: create request: %w", err)
	}

	q := url.Values{}
	q.Set("key", h.apiKey)
	q.Set("q", query)
	req.URL.RawQuery = q.Encode()

	resp, err := h.httpClient.Do(req)
	if err != nil {
		logger.Error("weatherapi: http call failed",
			slog.String("type", "outbound_error"),
			slog.String("query", query),
			slog.String("error", err.Error()),
		)
		return dto.WeatherAPIResponseDto{}, fmt.Errorf("weatherapi: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("weatherapi: unexpected status",
			slog.String("type", "outbound_error"),
			slog.String("query", query),
			slog.Int("status", resp.StatusCode),
		)
		return dto.WeatherAPIResponseDto{}, fmt.Errorf("weatherapi: unexpected status: %d", resp.StatusCode)
	}

	var payload dto.WeatherAPIResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		logger.Error("weatherapi: decode response failed",
			slog.String("type", "outbound_error"),
			slog.String("query", query),
			slog.String("error", err.Error()),
		)
		return dto.WeatherAPIResponseDto{}, fmt.Errorf("weatherapi: decode response: %w", err)
	}

	logger.Info("weatherapi: lookup succeeded",
		slog.String("type", "outbound_success"),
		slog.String("query", query),
		slog.Float64("temp_c", payload.Current.TempC),
	)

	return payload, nil
}

var _ outbound.WeatherProviderPort = (*HTTPClient)(nil)
