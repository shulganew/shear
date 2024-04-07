package ghandlers

import (
	"context"
	"net"
	"strings"

	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// @Summary      API get shortener statistic
// @Description  Get num of URLs and Users in the Shortener.
// @Tags         gRPC
// @Success      0 "OK"
// @Router       /api/internal/stats [get]
func (u *UsersServer) GetStat(ctx context.Context, in *pb.GetStatRequest) (*pb.GetStatResponse, error) {

	p, _ := peer.FromContext(ctx)
	ipport := p.Addr.String()
	ipStr, _, err := net.SplitHostPort(strings.TrimSpace(ipport))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Can't parse split host and port")
	}
	trust := ctx.Value(config.CtxIP{}).(string)

	_, nip, err := net.ParseCIDR(trust)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Can't parse CIDR IP form config.")
	}

	// Get source IP.
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, status.Errorf(codes.Internal, "Error get user ip.: ")
	}

	// Check if ip is allow for CIDR IP/mask.
	if !nip.Contains(ip) {
		return nil, status.Errorf(codes.PermissionDenied, "Source ip denay")
	}
	zap.S().Infoln("Source ip: ", ip)
	shorts, err := u.serviceURL.GetNumShorts(ctx)
	if err != nil {
		et := "Error during getting num of shorts (URLs): " + err.Error()
		zap.S().Errorln(et)
		return nil, status.Errorf(codes.Internal, et)
	}
	users, err := u.serviceURL.GetNumUsers(ctx)
	if err != nil {
		et := "Error during getting num of users: " + err.Error()
		zap.S().Errorln(et)
		return nil, status.Errorf(codes.Internal, et)
	}
	return &pb.GetStatResponse{Shorts: int64(shorts), Users: int64(users)}, nil
}
