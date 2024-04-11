package ghandlers

import (
	"context"
	"strconv"

	"github.com/shulganew/shear.git/internal/builders"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/entities"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"google.golang.org/grpc/status"
)

// @Summary      Add several original user's URLs and return generated briefs from response
// @Description  Add URLs. No user auth.
// @Tags         gRPC
// @Success      0   "Created"
// @Failure      3 "Wrong URL format"
// @Failure      6 "Conflict. URL existed."
// @Failure      13 "Handling error"
func (u *UsersServer) Batch(ctx context.Context, in *pb.BatchRequest) (*pb.BatchResponse, error) {
	// Get UserID from cxt values.
	ctxConfig := ctx.Value(config.CtxConfig{}).(config.CtxConfig)

	// Handle bach requests.
	var requests []entities.BatchRequest
	for i, r := range in.Origins {
		// Add batches.
		requests = append(requests, entities.BatchRequest{SessionID: strconv.Itoa(i), Origin: r})
	}

	// Convert to string array of origins.
	resDTO := u.serviceURL.AddBatch(ctx, builders.BatchRequestDTO{Origins: requests, CtxConfig: ctxConfig, Resp: u.conf.GetResponse()})

	// Check errors in service.
	if resDTO.Err != nil {
		//return nil, status.Errorf(resDTO.Status.GetStatusGRPC(), "Error in service answer: %s", resDTO.Err.Error())
		return nil, status.Errorf(resDTO.Status.GetStatusGRPC(), resDTO.Err.Error())
	}

	br := []string{}
	for _, b := range resDTO.AnwerURLs {
		br = append(br, b.Answer)
	}

	return &pb.BatchResponse{Briefs: br}, nil
}
