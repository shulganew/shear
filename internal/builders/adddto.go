package builders

import "github.com/shulganew/shear.git/internal/config"

// DTO for batch request.
type AddRequestDTO struct {
	Origin    string
	CtxConfig config.CtxConfig
	// Responce server address.
	Resp string
}

// DTO for batch responce.
type AddResponsehDTO struct {
	AnwerURL string
	Status   RespStatus
	Err      error
}
