package ghandlers

import (
	"context"
	"errors"
	"net/url"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/entities"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary      Add several original user's URLs and return generated briefs from response
// @Description  Add json URLs
// @Tags         gRPC
// @Success      0   "Created"
// @Failure      3 "Wrong URL format"
// @Failure      6 "Conflict. URL existed."
// @Failure      13 "Handling error"
func (u *UsersServer) Batch(ctx context.Context, in *pb.BatchRequest) (*pb.BatchResponse, error) {
	// Get UserID from cxt values.
	ctxConfig := ctx.Value(config.CtxConfig{}).(config.CtxConfig)
	if ctxConfig.IsNewUser() {
		return nil, status.Errorf(codes.PermissionDenied, "User not athorized")
	}
	var briefs []string
	shorts := []entities.Short{}
	for i, r := range in.Origins {
		var origin *url.URL
		origin, err := url.Parse(r)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Wrong URL format")
		}
		// get short brief and full answer URL
		brief := service.GenerateShortLinkByte()
		var answerURL *url.URL
		_, answerURL, err = u.serviceURL.GetAnsURLFast(origin.Scheme, u.conf.GetResponse(), brief)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Error parse answer")
		}

		// add batches
		briefs = append(briefs, answerURL.String())

		shortSession := entities.NewShort(i, ctxConfig.GetUserID(), brief, (*origin).String(), "")
		shorts = append(shorts, *shortSession)
	}
	zap.S().Infof("Original URLS: %+v \n", shorts)
	// save to storage
	err := u.serviceURL.SetAll(ctx, shorts)

	// check duplicated strings
	var tagErr *storage.ErrDuplicatedShort
	if err != nil {
		if errors.As(err, &tagErr) {
			// conflictR
			return &pb.BatchResponse{Briefs: []string{tagErr.Short.Brief}}, status.Errorf(codes.AlreadyExists, "Has existed original URL")

		}
		// create Ok answer
	}
	return &pb.BatchResponse{Briefs: briefs}, nil
}
