package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
	"github.com/biraneves/fc-labs-weather/internal/application/ports/outbound"
)

type HTTPClient struct {
	httpClient *http.Client
	baseURL    string
	timeout    time.Duration
}

func NewHTTPClient(httpClient *http.Client, baseURL string, timeout time.Duration) *HTTPClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if baseURL == "" {
		baseURL = "https://viacep.com.br/ws"
	}

	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	return &HTTPClient{
		httpClient: httpClient,
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		timeout:    timeout,
	}
}

func (h *HTTPClient) Find(ctx context.Context, request dto.ViaCEPRequestDto) (dto.ViaCEPResponseDto, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	url := fmt.Sprintf("%s/%s/json", h.baseURL, request.CEP.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return dto.ViaCEPResponseDto{}, fmt.Errorf("viacep: create request: %w", err)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return dto.ViaCEPResponseDto{}, fmt.Errorf("viacep: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return dto.ViaCEPResponseDto{}, outbound.ErrZipcodeNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return dto.ViaCEPResponseDto{}, fmt.Errorf("viacep: unexpected status: %d", resp.StatusCode)
	}

	var payload struct {
		dto.ViaCEPResponseDto
		Erro bool `json:"erro"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return dto.ViaCEPResponseDto{}, fmt.Errorf("viacep: decode response: %w", err)
	}

	if payload.Erro {
		return dto.ViaCEPResponseDto{}, outbound.ErrZipcodeNotFound
	}

	return payload.ViaCEPResponseDto, nil
}

var _ outbound.ZipcodeLookupPort = (*HTTPClient)(nil)
