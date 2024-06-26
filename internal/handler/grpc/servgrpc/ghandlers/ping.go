package ghandlers

import (
	"context"

	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Test DB connection.
// @Summary      Test database
// @Description  Ping service for database connection check
// @Tags         gRPC
// @Success      0 "Available"
// @Failure      13 "Database connection failed"
func (u *UsersServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	if err := u.db.PingContext(ctx); err == nil {
		return &pb.PingResponse{Ok: true}, nil
	} else {
		return nil, status.Errorf(codes.Internal, "Database connection failed.")
	}
}
