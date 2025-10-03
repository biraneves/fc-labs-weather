package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/biraneves/fc-labs-weather/internal/infrastructure/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecoveryMiddleware_Wrap(t *testing.T) {
	type fields struct {
		panicValue any
	}

	type expectations struct {
		status      int
		body        map[string]string
		logContains map[string]string
	}

	tests := []struct {
		name string
		f    fields
		exp  expectations
	}{
		{
			name: "panic recovered",
			f: fields{
				panicValue: "whatever",
			},
			exp: expectations{
				status: http.StatusInternalServerError,
				body: map[string]string{
					"error":             "internal_server_error",
					"error_description": "An unexpected error occurred while processing the request.",
				},
				logContains: map[string]string{
					"msg":   "panic recovered",
					"type":  "panic",
					"panic": "whatever",
				},
			},
		},
		{
			name: "no panic",
			exp: expectations{
				status: http.StatusOK,
				body: map[string]string{
					"status": "ok",
				},
				logContains: map[string]string{
					"msg": "no log expected",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))
			mw := server.NewRecoveryMiddleware(logger)

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.f.panicValue != nil {
					panic(tt.f.panicValue)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			})

			handler := mw.Wrap(next)

			req := httptest.NewRequest(http.MethodGet, "/weather", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			require.Equal(t, tt.exp.status, rec.Code)

			var gotBody map[string]string
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &gotBody))
			assert.Equal(t, tt.exp.body, gotBody)

			logContent := buf.String()
			if tt.f.panicValue == nil {
				assert.Empty(t, strings.TrimSpace(logContent))
				return
			}

			var logEntry map[string]any
			require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))

			for k, v := range tt.exp.logContains {
				value, ok := logEntry[k].(string)
				require.True(t, ok, "expected field %s in log", k)
				assert.Equal(t, v, value)
			}

			stacktrace, ok := logEntry["stacktrace"].(string)
			assert.True(t, ok)
			assert.NotEmpty(t, stacktrace)
		})
	}
}

func TestRecoveryMiddleware_FallbackLogger(t *testing.T) {
	mw := server.NewRecoveryMiddleware(nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := mw.Wrap(next)

	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLoggerFromContextIntegration(t *testing.T) {
	base := slog.New(slog.NewJSONHandler(io.Discard, nil))
	mw := server.NewRecoveryMiddleware(base)

	ctx := context.WithValue(context.Background(), "logger", base)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("fail")
	})

	handler := mw.Wrap(next)

	req := httptest.NewRequest(http.MethodGet, "/weather", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
