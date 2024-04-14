package ghandlers

import (
	"context"

	"github.com/shulganew/shear.git/internal/builders"
	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"google.golang.org/grpc/status"
)

// @Summary      Add origin URL
// @Description  set URL
// @Tags         gRPC
// @Success      nil
// @Failure      6 "Conflict. URL existed."
// @Failure      13 "Handling error"
func (u *UsersServer) AddURL(ctx context.Context, in *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	ctxConfig := ctx.Value(config.CtxConfig{}).(config.CtxConfig)

	// Convert to string array of origins.
	resDTO := u.serviceURL.AddURL(ctx, builders.AddRequestDTO{Origin: in.Origin, CtxConfig: ctxConfig, Resp: u.conf.GetResponse()})

	// Check errors in service.
	if resDTO.Err != nil {
		//return nil, status.Errorf(resDTO.Status.GetStatusGRPC(), "Error in service answer: %s", resDTO.Err.Error())
		return nil, status.Errorf(resDTO.Status.GetStatusGRPC(), resDTO.Err.Error())
	}

	return &pb.AddURLResponse{Brief: resDTO.AnwerURL}, nil
}
