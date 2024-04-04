package grpcs

import (
	"context"

	pb "github.com/shulganew/shear.git/internal/grpcs/proto"
)

func (us *UsersServer) Get(ctx context.Context, in *pb.GetURL) (*pb.GetURLResponse, error) {
	var resp pb.GetURLResponse

	origin, _, _ := us.serviceURL.GetOrigin(ctx, in.GetBrief())

	resp.Origin = origin
	return &resp, nil
}
