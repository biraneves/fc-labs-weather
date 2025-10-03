package server_test

import (
	"bufio"
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

func TestLoggerMiddleware_Wrap(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	mw := server.NewLoggerMiddleware(logger)

	var (
		nextCalled   bool
		capturedID   string
		innerMsgSeen bool
	)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true

		ctxLogger := server.LoggerFromContext(r.Context(), logger)
		ctxLogger.Info("inner", slog.String("type", "inner"))
		capturedID = server.RequestIDFromContext(r.Context())
		require.NotEmpty(t, capturedID)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("result"))
	})

	handler := mw.Wrap(next)

	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01001000", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.True(t, nextCalled)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "result", rec.Body.String())

	logs := decodeJSONLines(t, buf.String())
	require.Len(t, logs, 3)

	t.Run("request_in log", func(t *testing.T) {
		record := logs[0]
		assert.Equal(t, "request_in", record["msg"])
		assert.Equal(t, "request_in", record["type"])
		assert.Equal(t, "GET", record["method"])
		assert.Equal(t, "/weather", record["path"])
		assert.Equal(t, "cep=01001000", record["query"])
	})

	t.Run("inner log propagates request_id", func(t *testing.T) {
		record := logs[1]
		assert.Equal(t, "inner", record["msg"])
		assert.Equal(t, "inner", record["type"])

		reqID, ok := record["request_id"].(string)
		require.True(t, ok)
		assert.Equal(t, capturedID, reqID)

		innerMsgSeen = true
	})
	require.True(t, innerMsgSeen)

	t.Run("request_out log", func(t *testing.T) {
		record := logs[2]
		assert.Equal(t, "request_out", record["msg"])
		assert.Equal(t, "request_out", record["type"])

		status, ok := record["status"].(float64)
		require.True(t, ok)
		assert.Equal(t, float64(http.StatusCreated), status)

		payloadSize, ok := record["bytes"].(float64)
		require.True(t, ok)
		assert.Equal(t, float64(len("result")), payloadSize)

		_, hasDuration := record["duration"]
		assert.True(t, hasDuration)
	})
}

func TestLoggerFromContext_Fallback(t *testing.T) {
	base := slog.New(slog.NewJSONHandler(io.Discard, nil))
	got := server.LoggerFromContext(context.Background(), base)
	assert.Same(t, base, got)
}

func TestRequestIDFromContext_WithoutValue(t *testing.T) {
	assert.Empty(t, server.RequestIDFromContext(context.Background()))
}

func decodeJSONLines(t *testing.T, raw string) []map[string]any {
	t.Helper()

	var logs []map[string]any
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(raw)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var payload map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &payload))
		logs = append(logs, payload)
	}

	require.NoError(t, scanner.Err())

	return logs
}
