package handlers

import (
	"context"
	"errors"
	"net/url"

	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/grpcs/proto"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
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
func (us *UsersServer) GetURL(ctx context.Context, in *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	var resp pb.GetURLResponse
	origin, exist, isDeleted := us.serviceURL.GetOrigin(ctx, in.GetBrief())
	if exist {
		if isDeleted {
			return nil, status.Errorf(codes.Unknown, "Deleted: StatusGone")
		}
		resp.Origin = origin
		return &resp, nil
	}
	return nil, status.Errorf(codes.NotFound, "NotFound")
}

// @Summary      Set origin URL
// @Description  set URL
// @Tags         gRPC
// @Success      nil
// @Failure      7 "User unauthorized. PermissionDenied."
// @Failure      404 "Conflict. URL existed."
// @Failure      500 "Handling error"
func (us *UsersServer) SetURL(ctx context.Context, in *pb.SetURLRequest) (*pb.SetURLResponse, error) {
	// Get userID from context.
	ctxConfig := ctx.Value(config.CtxConfig{}).(config.CtxConfig)

	zap.S().Infoln("IsNewUser: ", ctxConfig.IsNewUser())
	zap.S().Infoln("IserID: ", ctxConfig.GetUserID())

	redirectURL, err := url.Parse(in.Origin)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parse URL error")
	}
	zap.S().Infoln("redirectURL: ", redirectURL)
	brief := service.GenerateShortLinkByte()
	mainURL, answerURL, err := us.serviceURL.GetAnsURLFast(redirectURL.Scheme, us.conf.GetResponse(), brief)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "parse new URL error")
	}

	// Save map to storage.
	err = us.serviceURL.SetURL(ctx, ctxConfig.GetUserID(), brief, (*redirectURL).String())
	if err != nil {
		var tagErr *storage.ErrDuplicatedURL
		if errors.As(err, &tagErr) {
			// Send existed string from error.
			var answer string
			answer, err = url.JoinPath(mainURL, tagErr.Brief)
			if err != nil {
				zap.S().Errorln("Error during JoinPath", err)
			}
			return &pb.SetURLResponse{Brief: answer}, status.Errorf(codes.AlreadyExists, "StatusConflict AlreadyExists")
		}

		zap.S().Errorln(err)
		return nil, status.Errorf(codes.Internal, "Error saving in Storage")
	}

	return &pb.SetURLResponse{Brief: answerURL.String()}, nil
}
