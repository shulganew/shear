package ghandlers

import (
	"context"

	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Delete User's URLs from array (mark as deleted with saving in DB)
// @Summary      Delete user's URLs
// @Description  Delete array  of user's URLs in database, async
// @Tags         gRPC
// @Success      0 "Accepted"
// @Failure      7 "User unauthorized"
func (u *UsersServer) DelUserURLs(ctx context.Context, in *pb.DelRequest) (*pb.DelResponse, error) {
	// Get UserID from cxt values.
	ctxConfig := ctx.Value(config.CtxConfig{}).(config.CtxConfig)
	if ctxConfig.IsNewUser() {
		return nil, status.Errorf(codes.PermissionDenied, "User not athorized")
	}

	// Async delete Shorts from body
	u.servDelete.AsyncDelete(ctxConfig.GetUserID(), in.Briefs)
	return &pb.DelResponse{Ok: true}, nil
}
