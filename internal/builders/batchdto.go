package builders

import (
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/entities"
)

// DTO for batch request.
type BatchRequestDTO struct {
	Origins   []entities.BatchRequest
	CtxConfig config.CtxConfig
	Resp      string
}

// DTO for batch responce.
type BatchResponsehDTO struct {
	AnwerURLs []entities.BatchResponse
	Status    RespStatus
	Err       error
}
