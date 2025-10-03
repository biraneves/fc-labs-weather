package server_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
	"github.com/biraneves/fc-labs-weather/internal/application/ports/inbound"
	"github.com/biraneves/fc-labs-weather/internal/application/usecase"
	"github.com/biraneves/fc-labs-weather/internal/domain/entity"
	"github.com/biraneves/fc-labs-weather/internal/infrastructure/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeUseCase struct {
	resp     dto.RequestOutDto
	err      error
	called   bool
	received dto.RequestInDto
}

func (f *fakeUseCase) Execute(ctx context.Context, request dto.RequestInDto) (dto.RequestOutDto, error) {
	f.called = true
	f.received = request
	if f.err != nil {
		return dto.RequestOutDto{}, f.err
	}

	return f.resp, nil
}

var _ inbound.GetWeatherByCEPUseCase = (*fakeUseCase)(nil)

var noopLogger = slog.New(slog.NewJSONHandler(io.Discard, nil))

func TestHandler_HandleWeather(t *testing.T) {
	tempC, err := entity.NewTemperatureCelsius(28.5)
	require.NoError(t, err)

	tempF, err := entity.NewTemperatureFahrenheit(tempC.ToFahrenheit())
	require.NoError(t, err)

	tempK, err := entity.NewTemperatureKelvin(tempC.ToKelvin())
	require.NoError(t, err)

	type fields struct {
		useCaseResp dto.RequestOutDto
		useCaseErr  error
	}

	type expectations struct {
		status            int
		bodyEquals        string
		bodyContains      string
		expectUseCaseCall bool
		expectedCEP       entity.Cep
		method            string
		url               string
	}

	tests := []struct {
		name string
		f    fields
		exp  expectations
	}{
		{
			name: "success",
			f: fields{
				useCaseResp: dto.RequestOutDto{
					TempC: tempC,
					TempF: tempF,
					TempK: tempK,
				},
			},
			exp: expectations{
				status:            http.StatusOK,
				bodyEquals:        `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.7}`,
				expectUseCaseCall: true,
				expectedCEP:       entity.Cep("01001000"),
				method:            http.MethodGet,
				url:               "/weather?cep=01001000",
			},
		},
		{
			name: "method not allowed",
			exp: expectations{
				status:            http.StatusMethodNotAllowed,
				bodyContains:      http.StatusText(http.StatusMethodNotAllowed),
				expectUseCaseCall: false,
				method:            http.MethodPost,
				url:               "/weather?cep=01001000",
			},
		},
		{
			name: "missing cep parameter",
			exp: expectations{
				status:            http.StatusBadRequest,
				bodyEquals:        `{"error":"missing query parameter: cep"}`,
				expectUseCaseCall: false,
				method:            http.MethodGet,
				url:               "/weather",
			},
		},
		{
			name: "invalid cep format",
			exp: expectations{
				status:            http.StatusUnprocessableEntity,
				bodyEquals:        `{"error":"invalid zipcode"}`,
				expectUseCaseCall: false,
				method:            http.MethodGet,
				url:               "/weather?cep=123",
			},
		},
		{
			name: "use case invalid zipcode",
			f: fields{
				useCaseErr: usecase.ErrInvalidZipCode,
			},
			exp: expectations{
				status:            http.StatusUnprocessableEntity,
				bodyEquals:        `{"error":"invalid zipcode"}`,
				expectUseCaseCall: true,
				expectedCEP:       entity.Cep("01001000"),
				method:            http.MethodGet,
				url:               "/weather?cep=01001000",
			},
		},
		{
			name: "use case zipcode not found",
			f: fields{
				useCaseErr: usecase.ErrZipcodeNotFound,
			},
			exp: expectations{
				status:            http.StatusNotFound,
				bodyEquals:        `{"error":"cannot find zipcode"}`,
				expectUseCaseCall: true,
				expectedCEP:       entity.Cep("01001000"),
				method:            http.MethodGet,
				url:               "/weather?cep=01001000",
			},
		},
		{
			name: "use case generic error",
			f: fields{
				useCaseErr: errors.New("whatever"),
			},
			exp: expectations{
				status:            http.StatusInternalServerError,
				bodyEquals:        `{"error":"internal error"}`,
				expectUseCaseCall: true,
				expectedCEP:       entity.Cep("01001000"),
				method:            http.MethodGet,
				url:               "/weather?cep=01001000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &fakeUseCase{resp: tt.f.useCaseResp, err: tt.f.useCaseErr}
			handler := server.NewHandler(uc, noopLogger)
			mux := http.NewServeMux()
			handler.RegisterRoutes(mux)

			req := httptest.NewRequest(tt.exp.method, tt.exp.url, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			require.Equal(t, tt.exp.status, rec.Code)

			if tt.exp.bodyEquals != "" {
				assert.JSONEq(t, tt.exp.bodyEquals, rec.Body.String())
			} else if tt.exp.bodyContains != "" {
				assert.True(t, strings.Contains(rec.Body.String(), tt.exp.bodyContains))
			}

			assert.Equal(t, tt.exp.expectUseCaseCall, uc.called)
			if tt.exp.expectUseCaseCall {
				assert.Equal(t, tt.exp.expectedCEP, uc.received.CEP)
			}
		})
	}
}
