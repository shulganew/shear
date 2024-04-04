package handlers

import (
	"context"

	pb "github.com/shulganew/shear.git/internal/grpcs/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	//brief := service.GenerateShortLinkByte()
	//mainURL, answerURL, err := us.serviceURL.GetAnsURLFast(service.SchemaHTTP, us.conf.GetResponse(), brief)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "Error parse URL")
	// }

	// Get userID from context.
	var userID string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("user_id")
		if len(values) > 0 {
			userID = values[0]
		}
	} else {
		return nil, status.Errorf(codes.PermissionDenied, "Can't find user in metadata")
	}
	zap.S().Infoln("IserID: ", userID)
	return nil, nil

	/*
		// Save map to storage
		err = us.serviceURL.SetURL(req.Context(), userID.Value, brief, (*redirectURL).String())

		if err != nil {
			var tagErr *storage.ErrDuplicatedURL
			if errors.As(err, &tagErr) {
				// set status code 409 Conflict
				res.WriteHeader(http.StatusConflict)

				//send existed string from error
				var answer string
				answer, err = url.JoinPath(mainURL, tagErr.Brief)
				if err != nil {
					zap.S().Errorln("Error during JoinPath", err)
				}
				res.Write([]byte(answer))
				return
			}

			zap.S().Errorln(err)
			http.Error(res, "Error saving in Storage.", http.StatusInternalServerError)
		}
		// set status code 201
		res.WriteHeader(http.StatusCreated)
		// send generate and saved string
		res.Write([]byte(answerURL.String()))
	*/
}
