package ghandlers

import (
	"context"
	"errors"
	"net/url"

	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary      Add origin URL
// @Description  set URL
// @Tags         gRPC
// @Success      nil
// @Failure      6 "Conflict. URL existed."
// @Failure      13 "Handling error"
func (u *UsersServer) AddURL(ctx context.Context, in *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	// Get userID from context.
	ctxConfig := ctx.Value(config.CtxConfig{}).(config.CtxConfig)
	redirectURL, err := url.Parse(in.Origin)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parse URL error")
	}
	zap.S().Infoln("redirectURL: ", redirectURL)
	brief := service.GenerateShortLinkByte()
	mainURL, answerURL, err := u.serviceURL.GetAnsURLFast(redirectURL.Scheme, u.conf.GetResponse(), brief)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parse new URL error")
	}

	// Save map to storage.
	err = u.serviceURL.AddURL(ctx, ctxConfig.GetUserID(), brief, (*redirectURL).String())
	if err != nil {
		var tagErr *service.ErrDuplicatedURL
		if errors.As(err, &tagErr) {
			// Send existed string from error.
			var answer string
			answer, err = url.JoinPath(mainURL, tagErr.Brief)
			if err != nil {
				zap.S().Errorln("Error during JoinPath", err)
			}
			return &pb.AddURLResponse{Brief: answer}, status.Errorf(codes.AlreadyExists, "StatusConflict AlreadyExists")
		}

		zap.S().Errorln(err)
		return nil, status.Errorf(codes.Internal, "Error saving in Storage")
	}

	return &pb.AddURLResponse{Brief: answerURL.String()}, nil
}
