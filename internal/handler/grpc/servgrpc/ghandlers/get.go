package ghandlers

import (
	"context"

	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GET and redirect by brief.
// @Summary      Get origin URL by brief (short) URL
// @Description  get short by id
// @Tags         gRPC
// @Param        id   path  string  true  "brief URL"
// @Success      nil
// @Failure      2
// @Failure      5
func (u *UsersServer) GetURL(ctx context.Context, in *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	var resp pb.GetURLResponse
	origin, exist, isDeleted := u.serviceURL.GetOrigin(ctx, in.GetBrief())
	if exist {
		if isDeleted {
			return nil, status.Errorf(codes.Unknown, "Deleted: StatusGone")
		}
		resp.Origin = origin
		return &resp, nil
	}
	return nil, status.Errorf(codes.NotFound, "NotFound")
}
