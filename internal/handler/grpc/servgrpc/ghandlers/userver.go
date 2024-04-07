package ghandlers

import (
	"database/sql"

	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"github.com/shulganew/shear.git/internal/service"
)

var _ pb.UsersServer = (*UsersServer)(nil)

// Base gRPC struct.
type UsersServer struct {
	pb.UnimplementedUsersServer
	serviceURL *service.Shorten
	conf       *config.Config
	db         *sql.DB
	servDelete *service.Delete
}

// UserServer gRPC constructor.
func NewUsersServer(serviceURL *service.Shorten, conf *config.Config, db *sql.DB, servDelete *service.Delete) *UsersServer {
	return &UsersServer{serviceURL: serviceURL, conf: conf, db: db, servDelete: servDelete}

}
