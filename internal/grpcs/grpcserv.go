package grpcs

import (
	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/grpcs/proto"
	"github.com/shulganew/shear.git/internal/service"
)

var _ pb.UsersServer = (*UsersServer)(nil)

// Base gRPC struct.
type UsersServer struct {
	pb.UnimplementedUsersServer
	serviceURL *service.Shorten
	conf       *config.Config
}

// UserServer gRPC constructor.
func NewUsersServer(serviceURL *service.Shorten, conf *config.Config) *UsersServer {
	return &UsersServer{serviceURL: serviceURL, conf: conf}
}
