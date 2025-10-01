package outbound

import (
	"context"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
)

type ZipcodeLookupPort interface {
	Find(ctx context.Context, request dto.ViaCEPRequestDto) (dto.ViaCEPResponseDto, error)
}
