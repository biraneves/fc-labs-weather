package inbound

import (
	"context"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
)

type GetWeatherByCEPUseCase interface {
	Execute(ctx context.Context, request dto.RequestInDto) (dto.RequestOutDto, error)
}
