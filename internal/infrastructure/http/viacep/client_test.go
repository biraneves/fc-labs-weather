package viacep_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
	"github.com/biraneves/fc-labs-weather/internal/infrastructure/http/viacep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPClient(t *testing.T) {
	client := viacep.NewHTTPClient(nil, "", 0)
	require.NotNil(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.Find(ctx, dto.ViaCEPRequestDto{})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.DeadlineExceeded) || err != nil)
}

func TestHTTPClient_Find(t *testing.T) {
	type fields struct {
		status int
		body   string
		delay  time.Duration
	}

	tests := []struct {
		name          string
		fields        fields
		want          dto.ViaCEPResponseDto
		wantError     error
		expectedError string
	}{
		{
			name: "success",
			fields: fields{
				status: http.StatusOK,
				body: `{
				  "cep": "01001-000",
				  "logradouro": "Praça da Sé",
				  "localidade": "São Paulo",
				  "uf": "SP"
				}`,
			},
			want: dto.ViaCEPResponseDto{
				CEP:        "01001-000",
				Logradouro: "Praça da Sé",
				Localidade: "São Paulo",
				UF:         "SP",
			},
		},
		{
			name: "zipcode not found - status 404",
			fields: fields{
				status: http.StatusNotFound,
				body:   `{}`,
			},
			wantError:     viacep.ErrZipcodeNotFound,
			expectedError: viacep.ErrZipcodeNotFound.Error(),
		},
		{
			name: "zipcode not found - field erro equals true",
			fields: fields{
				status: http.StatusOK,
				body:   `{"erro": true}`,
			},
			wantError:     viacep.ErrZipcodeNotFound,
			expectedError: viacep.ErrZipcodeNotFound.Error(),
		},
		{
			name: "unexpected error",
			fields: fields{
				status: http.StatusInternalServerError,
				body:   `{}`,
			},
			expectedError: "viacep: unexpected status: 500",
		},
		{
			name: "invalid json",
			fields: fields{
				status: http.StatusOK,
				body:   `{"cep":`,
			},
			expectedError: "viacep: decode response",
		},
		{
			name: "context timeout",
			fields: fields{
				status: http.StatusOK,
				body: `{
				  "cep": "01001-000",
				  "localidade": "São Paulo",
				  "uf": "SP"
				}`,
				delay: 200 * time.Millisecond,
			},
			expectedError: "context deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.fields.delay > 0 {
					time.Sleep(tt.fields.delay)
				}

				w.WriteHeader(tt.fields.status)
				_, _ = w.Write([]byte(tt.fields.body))
			}))
			defer server.Close()

			client := viacep.NewHTTPClient(nil, server.URL, 50*time.Millisecond)

			ctx := context.Background()
			if tt.fields.delay > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 50*time.Millisecond)
				defer cancel()
			}

			got, err := client.Find(ctx, dto.ViaCEPRequestDto{})
			if tt.wantError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantError)
				return
			}

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
