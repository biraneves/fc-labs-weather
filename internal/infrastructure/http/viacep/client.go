package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
	"github.com/biraneves/fc-labs-weather/internal/application/ports/outbound"
	"github.com/biraneves/fc-labs-weather/internal/infrastructure/http/server"
)

type HTTPClient struct {
	httpClient *http.Client
	baseURL    string
	timeout    time.Duration
	logger     *slog.Logger
}

func NewHTTPClient(httpClient *http.Client, baseURL string, timeout time.Duration, logger *slog.Logger) *HTTPClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if baseURL == "" {
		baseURL = "https://viacep.com.br/ws"
	}

	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &HTTPClient{
		httpClient: httpClient,
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		timeout:    timeout,
		logger:     logger,
	}
}

func (h *HTTPClient) Find(ctx context.Context, request dto.ViaCEPRequestDto) (dto.ViaCEPResponseDto, error) {
	logger := server.LoggerFromContext(ctx, h.logger)

	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	url := fmt.Sprintf("%s/%s/json", h.baseURL, request.CEP.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.Error("viacep: create request failed",
			slog.String("type", "outbound_error"),
			slog.String("cep", request.CEP.String()),
			slog.String("error", err.Error()),
		)
		return dto.ViaCEPResponseDto{}, fmt.Errorf("viacep: create request: %w", err)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		logger.Error("viacep: http call failed",
			slog.String("type", "outbound_error"),
			slog.String("cep", request.CEP.String()),
			slog.String("error", err.Error()),
		)
		return dto.ViaCEPResponseDto{}, fmt.Errorf("viacep: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		logger.Info("viacep: zipcode not found",
			slog.String("type", "outbound_error"),
			slog.String("cep", request.CEP.String()),
		)
		return dto.ViaCEPResponseDto{}, outbound.ErrZipcodeNotFound
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("viacep: unexpected status",
			slog.String("type", "outbound_error"),
			slog.String("cep", request.CEP.String()),
			slog.Int("status", resp.StatusCode),
		)
		return dto.ViaCEPResponseDto{}, fmt.Errorf("viacep: unexpected status: %d", resp.StatusCode)
	}

	var payload struct {
		dto.ViaCEPResponseDto
		Erro any `json:"erro"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		logger.Error("viacep: decode response failed",
			slog.String("type", "outbound_error"),
			slog.String("cep", request.CEP.String()),
			slog.String("error", err.Error()),
		)
		return dto.ViaCEPResponseDto{}, fmt.Errorf("viacep: decode response: %w", err)
	}

	if normalizeErrorFlag(payload.Erro) {
		logger.Info("viacep: response flagged erro=true",
			slog.String("type", "outbound_error"),
			slog.String("cep", request.CEP.String()),
		)
		return dto.ViaCEPResponseDto{}, outbound.ErrZipcodeNotFound
	}

	logger.Info("viacep: lookup succeeded",
		slog.String("type", "outbound_success"),
		slog.String("cep", request.CEP.String()),
		slog.String("localidade", payload.Localidade),
		slog.String("uf", payload.UF),
	)

	return payload.ViaCEPResponseDto, nil
}

func normalizeErrorFlag(v any) bool {
	switch val := v.(type) {
	case bool:
		return val

	case string:
		return strings.EqualFold(strings.TrimSpace(val), "true")

	case json.Number:
		return val == "1"

	case nil:
		return false

	default:
		return false
	}
}

var _ outbound.ZipcodeLookupPort = (*HTTPClient)(nil)
