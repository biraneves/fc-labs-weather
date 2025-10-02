package weatherapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
	"github.com/biraneves/fc-labs-weather/internal/application/ports/outbound"
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
}

func NewHTTPClient(httpClient *http.Client, baseURL, apiKey string, timeout time.Duration) *HTTPClient {
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
	}
}

func (h *HTTPClient) FetchCurrent(ctx context.Context, request dto.WeatherAPIRequestDto) (dto.WeatherAPIResponseDto, error) {
	if h.apiKey == "" {
		return dto.WeatherAPIResponseDto{}, ErrMissingAPIKey
	}

	query := strings.TrimSpace(request.Q)
	if query == "" {
		return dto.WeatherAPIResponseDto{}, ErrEmptyQuery
	}

	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	endpoint := fmt.Sprintf("%s/current.json", h.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return dto.WeatherAPIResponseDto{}, fmt.Errorf("weatherapi: create request: %w", err)
	}

	q := url.Values{}
	q.Set("key", h.apiKey)
	q.Set("q", query)
	req.URL.RawQuery = q.Encode()

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return dto.WeatherAPIResponseDto{}, fmt.Errorf("weatherapi: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.WeatherAPIResponseDto{}, fmt.Errorf("weatherapi: unexpected status: %d", resp.StatusCode)
	}

	var payload dto.WeatherAPIResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return dto.WeatherAPIResponseDto{}, fmt.Errorf("weatherapi: decode response: %w", err)
	}

	return payload, nil
}

var _ outbound.WeatherProviderPort = (*HTTPClient)(nil)
