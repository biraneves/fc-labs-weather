package server

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestResponseRecorder(t *testing.T) {
    tests := []struct {
        name           string
        action         func(t *testing.T, rec *responseRecorder, base *httptest.ResponseRecorder)
        expectedStatus int
        expectedBytes  int
        expectedBody   string
    }{
        {
            name:           "default status",
            expectedStatus: http.StatusOK,
        },
        {
            name: "write header",
            action: func(t *testing.T, rec *responseRecorder, base *httptest.ResponseRecorder) {
                rec.WriteHeader(http.StatusAccepted)
            },
            expectedStatus: http.StatusAccepted,
        },
        {
            name: "write body",
            action: func(t *testing.T, rec *responseRecorder, base *httptest.ResponseRecorder) {
                _, err := rec.Write([]byte("hello"))
                require.NoError(t, err)
            },
            expectedStatus: http.StatusOK,
            expectedBytes:  5,
            expectedBody:   "hello",
        },
        {
            name: "multiple writes accumulate",
            action: func(t *testing.T, rec *responseRecorder, base *httptest.ResponseRecorder) {
                _, err := rec.Write([]byte("foo"))
                require.NoError(t, err)
                _, err = rec.Write([]byte("bar"))
                require.NoError(t, err)
            },
            expectedStatus: http.StatusOK,
            expectedBytes:  6,
            expectedBody:   "foobar",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            base := httptest.NewRecorder()
            rec := newResponseRecorder(base)

            if tt.action != nil {
                tt.action(t, rec, base)
            }

            assert.Equal(t, tt.expectedStatus, rec.status)
            assert.Equal(t, tt.expectedStatus, base.Result().StatusCode)
            assert.Equal(t, tt.expectedBytes, rec.bytes)
            if tt.expectedBody != "" {
                assert.Equal(t, tt.expectedBody, base.Body.String())
            } else {
                assert.Equal(t, "", base.Body.String())
            }
        })
    }
}
