package outbound

import (
	"context"
	"errors"

	"github.com/biraneves/fc-labs-weather/internal/application/dto"
)

var ErrZipcodeNotFound = errors.New("zipcode lookup: not found")

type ZipcodeLookupPort interface {
	Find(ctx context.Context, request dto.ViaCEPRequestDto) (dto.ViaCEPResponseDto, error)
}
