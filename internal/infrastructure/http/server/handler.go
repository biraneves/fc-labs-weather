package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
	"github.com/biraneves/fc-labs-weather/internal/application/ports/inbound"
	"github.com/biraneves/fc-labs-weather/internal/application/usecase"
	"github.com/biraneves/fc-labs-weather/internal/domain/entity"
)

type Handler struct {
	useCase inbound.GetWeatherByCEPUseCase
	logger  *slog.Logger
}

func NewHandler(uc inbound.GetWeatherByCEPUseCase, logger *slog.Logger) *Handler {
	return &Handler{useCase: uc, logger: logger}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /weather", h.handleWeather)
}

type errorBody struct {
	Error string `json:"error"`
}

func (h *Handler) handleWeather(w http.ResponseWriter, r *http.Request) {
	logger := LoggerFromContext(r.Context(), h.logger)

	if r.Method != http.MethodGet {
		logger.Warn("unsupported method",
			slog.String("type", "handler_error"),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		logger.Warn("missing cep query parameter",
			slog.String("type", "handler_error"),
			slog.String("query", r.URL.RawQuery),
		)
		writeError(w, http.StatusBadRequest, "missing query parameter: cep")
		return
	}

	cepToSearch, err := entity.NewCep(cepParam)
	if err != nil {
		logger.Warn("invalid cep received",
			slog.String("type", "handler_error"),
			slog.String("cep", cepParam),
			slog.String("error", err.Error()),
		)
		writeError(w, http.StatusUnprocessableEntity, usecase.ErrInvalidZipCode.Error())
		return
	}

	out, err := h.useCase.Execute(r.Context(), dto.RequestInDto{CEP: cepToSearch})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidZipCode):
			logger.Warn("use case rejected cep as invalid",
				slog.String("type", "handler_error"),
				slog.String("cep", cepParam),
			)
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return

		case errors.Is(err, usecase.ErrZipcodeNotFound):
			logger.Info("zipcode not found",
				slog.String("type", "handler_error"),
				slog.String("cep", cepParam),
			)
			writeError(w, http.StatusNotFound, err.Error())
			return

		default:
			logger.Error("unexpected failure executing use case",
				slog.String("type", "handler_error"),
				slog.String("cep", cepParam),
				slog.String("error", err.Error()),
			)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorBody{Error: message})
}
