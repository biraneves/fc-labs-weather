package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
	"github.com/biraneves/fc-labs-weather/internal/application/ports/outbound"
	"github.com/biraneves/fc-labs-weather/internal/application/usecase"
	"github.com/biraneves/fc-labs-weather/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGetWeatherByCEPUseCase(t *testing.T) {
	tests := []struct {
		name        string
		zipMock     outbound.ZipcodeLookupPort
		weatherMock outbound.WeatherProviderPort
	}{
		{
			name:        "success",
			zipMock:     &fakeZipcodePort{},
			weatherMock: &fakeWeatherPort{},
		},
		{
			name:        "nil zipcode port",
			zipMock:     nil,
			weatherMock: &fakeWeatherPort{},
		},
		{
			name:        "nil weather port",
			zipMock:     &fakeZipcodePort{},
			weatherMock: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := usecase.NewGetWeatherByCEPUseCase(tt.zipMock, tt.weatherMock)

			require.NotNil(t, uc)

			concrete := uc.(*usecase.GetWeatherByCEPUseCase)
			require.Equal(t, tt.zipMock, concrete.Zipcode)
			require.Equal(t, tt.weatherMock, concrete.Weather)
		})
	}
}

func TestGetWeatherByCEPUseCase_Execute(t *testing.T) {
	validCEP, _ := entity.NewCep("01001000")
	invalidCEP := entity.Cep("123")

	baseWeatherResp := dto.WeatherAPIResponseDto{}
	baseWeatherResp.Current.TempC = 25.0

	tests := []struct {
		name          string
		input         dto.RequestInDto
		zipcodeStub   outbound.ZipcodeLookupPort
		weatherStub   outbound.WeatherProviderPort
		expectedError string
		assertSuccess func(t *testing.T, out dto.RequestOutDto)
	}{
		{
			name:          "invalid zipcode",
			input:         dto.RequestInDto{CEP: invalidCEP},
			zipcodeStub:   fakeZipcodePort{},
			weatherStub:   fakeWeatherPort{},
			expectedError: usecase.ErrInvalidZipCode.Error(),
		},
		{
			name:          "zipcode not found",
			input:         dto.RequestInDto{CEP: validCEP},
			zipcodeStub:   fakeZipcodePort{err: usecase.ErrZipcodeNotFound},
			weatherStub:   fakeWeatherPort{},
			expectedError: usecase.ErrZipcodeNotFound.Error(),
		},
		{
			name:          "zipcode service failed",
			input:         dto.RequestInDto{CEP: validCEP},
			zipcodeStub:   fakeZipcodePort{err: errors.New("service unavailable")},
			weatherStub:   fakeWeatherPort{},
			expectedError: "zipcode lookup failed",
		},
		{
			name:          "empty locality",
			input:         dto.RequestInDto{CEP: validCEP},
			zipcodeStub:   fakeZipcodePort{resp: dto.ViaCEPResponseDto{Localidade: ""}},
			weatherStub:   fakeWeatherPort{},
			expectedError: usecase.ErrZipcodeNotFound.Error(),
		},
		{
			name:          "weather service failed",
			input:         dto.RequestInDto{CEP: validCEP},
			zipcodeStub:   fakeZipcodePort{resp: dto.ViaCEPResponseDto{Localidade: "São Paulo"}},
			weatherStub:   fakeWeatherPort{err: errors.New("weather api unavailable")},
			expectedError: "weather provider failed",
		},
		{
			name:        "success",
			input:       dto.RequestInDto{CEP: validCEP},
			zipcodeStub: fakeZipcodePort{resp: dto.ViaCEPResponseDto{Localidade: "São Paulo"}},
			weatherStub: fakeWeatherPort{resp: baseWeatherResp},
			assertSuccess: func(t *testing.T, out dto.RequestOutDto) {
				require.True(t, out.TempC.IsValid())
				require.True(t, out.TempF.IsValid())
				require.True(t, out.TempK.IsValid())
				assert.InDelta(t, 25.0, out.TempC.Value(), 1e-3)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := usecase.GetWeatherByCEPUseCase{
				Zipcode: tt.zipcodeStub,
				Weather: tt.weatherStub,
			}

			out, err := uc.Execute(context.Background(), tt.input)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tt.assertSuccess)
				tt.assertSuccess(t, out)
			}
		})
	}
}

type fakeZipcodePort struct {
	resp dto.ViaCEPResponseDto
	err  error
}

func (z fakeZipcodePort) Find(ctx context.Context, request dto.ViaCEPRequestDto) (dto.ViaCEPResponseDto, error) {
	return z.resp, z.err
}

type fakeWeatherPort struct {
	resp dto.WeatherAPIResponseDto
	err  error
}

func (w fakeWeatherPort) FetchCurrent(ctx context.Context, req dto.WeatherAPIRequestDto) (dto.WeatherAPIResponseDto, error) {
	return w.resp, w.err
}
