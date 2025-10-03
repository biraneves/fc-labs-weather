package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type RecoveryMiddleware struct {
	logger *slog.Logger
}

func NewRecoveryMiddleware(logger *slog.Logger) *RecoveryMiddleware {
	if logger == nil {
		logger = slog.Default()
	}

	return &RecoveryMiddleware{logger: logger}
}

func (m *RecoveryMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger := LoggerFromContext(r.Context(), m.logger)

				logger.Error("panic recovered",
					slog.String("type", "panic"),
					slog.String("panic", fmt.Sprintf("%v", rec)),
					slog.String("stacktrace", string(debug.Stack())),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error":             "internal_server_error",
					"error_description": "An unexpected error occurred while processing the request.",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
