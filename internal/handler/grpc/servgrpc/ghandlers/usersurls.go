package ghandlers

import (
	"context"
	"encoding/binary"

	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Return all users briefs URLs (shorts).
// @Summary      Get all user's origin URLs
// @Description  Return Short object with origin and brief URL.
// @Tags         gRPC
// @Success      0 {object} []ResponseAuth "OK"
// @Failure      7 "User unauthorized"
// @Failure      2 "No contenet for user"
func (u *UsersServer) GetUserURLs(ctx context.Context, in *pb.GetURLs) (*pb.GetUserURLsResponse, error) {
	// Get UserID from cxt values.
	ctxConfig := ctx.Value(config.CtxConfig{}).(config.CtxConfig)
	if ctxConfig.IsNewUser() {
		return nil, status.Errorf(codes.PermissionDenied, "User not athorized")
	}

	userID := ctxConfig.GetUserID()

	// Get Short URLs for userID.
	shorts := u.serviceURL.GetUserAll(ctx, userID)
	zap.S().Infof("Found: %d saved URL for User with ID: %s", len(shorts), userID)

	// If no data - codes.Code = 2.
	if len(shorts) == 0 {
		return nil, status.Errorf(codes.Unknown, "No contenet for user.")
	}

	var responce []*pb.Short
	for _, short := range shorts {
		_, answerURL, err := u.serviceURL.GetAnsURLFast("http", u.conf.GetResponse(), short.Brief)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Can't get ansfer URL")
		}
		responce = append(responce, &pb.Short{Brief: answerURL.String(), Origin: short.Origin})
	}

	zap.S().Infoln("Server answer with user's short URLs in gRPC: ", binary.Size(responce))
	return &pb.GetUserURLsResponse{Short: responce}, nil
}
