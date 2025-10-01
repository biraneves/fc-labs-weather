package outbound

import (
	"context"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
)

type WeatherProviderPort interface {
	FetchCurrent(ctx context.Context, req dto.WeatherAPIRequestDto) (dto.WeatherAPIResponseDto, error)
}
