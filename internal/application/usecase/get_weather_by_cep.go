package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
	"github.com/biraneves/fc-labs-weather/internal/application/ports/inbound"
	"github.com/biraneves/fc-labs-weather/internal/application/ports/outbound"
	"github.com/biraneves/fc-labs-weather/internal/domain/entity"
)

var (
	ErrInvalidZipCode  = errors.New("invalid zipcode")     // -> 422
	ErrZipcodeNotFound = errors.New("cannot find zipcode") // -> 404
)

type GetWeatherByCEPUseCase struct {
	Zipcode outbound.ZipcodeLookupPort
	Weather outbound.WeatherProviderPort
}

func NewGetWeatherByCEPUseCase(zipcode outbound.ZipcodeLookupPort, weather outbound.WeatherProviderPort) inbound.GetWeatherByCEPUseCase {
	return &GetWeatherByCEPUseCase{zipcode, weather}
}

func (g GetWeatherByCEPUseCase) Execute(ctx context.Context, request dto.RequestInDto) (dto.RequestOutDto, error) {
	cep, err := entity.NewCep(request.CEP.String())
	if err != nil {
		return dto.RequestOutDto{}, ErrInvalidZipCode
	}

	viaResp, err := g.Zipcode.Find(ctx, dto.ViaCEPRequestDto{CEP: cep})
	if err != nil {
		if errors.Is(err, ErrZipcodeNotFound) {
			return dto.RequestOutDto{}, ErrZipcodeNotFound
		}
		return dto.RequestOutDto{}, fmt.Errorf("zipcode lookup failed: %w", err)
	}

	city := strings.TrimSpace(viaResp.Localidade)
	if city == "" {
		return dto.RequestOutDto{}, ErrZipcodeNotFound
	}

	weatherResp, err := g.Weather.FetchCurrent(ctx, dto.WeatherAPIRequestDto{Q: city})
	if err != nil {
		return dto.RequestOutDto{}, fmt.Errorf("weather provider failed: %w", err)
	}

	tempC, err := entity.NewTemperatureCelsius(weatherResp.Current.TempC)
	if err != nil {
		return dto.RequestOutDto{}, fmt.Errorf("weather provider returned invalid celsius temperature: %w", err)
	}

	tempF, err := entity.NewTemperatureFahrenheit(tempC.ToFahrenheit())
	if err != nil {
		return dto.RequestOutDto{}, fmt.Errorf("invalid fahrenheit conversion: %w", err)
	}

	tempK, err := entity.NewTemperatureKelvin(tempC.ToKelvin())
	if err != nil {
		return dto.RequestOutDto{}, fmt.Errorf("invalid kelvin conversion: %w", err)
	}

	return dto.RequestOutDto{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}
