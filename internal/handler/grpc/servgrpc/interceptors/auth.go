package interceptors

import (
	"context"

	"github.com/google/uuid"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Check user auth from MD value (cookie analog of authmiddleware).
func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Pass from config.
	pass := ctx.Value(config.CtxPassKey{}).(string)

	// Check userID in gRPC ctx.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata: user id")
	}
	isNewUser := false
	values := md.Get("user_id")
	var cUserID string
	if len(values) > 0 {
		cUserID = values[0]
	} else {
		isNewUser = true
	}
	// User UUID existed in ctx.
	if !isNewUser {
		userID, err := service.DecodeCookie(cUserID, pass)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "decode user id failed")
		}

		// check correct UUID
		_, err = uuid.Parse(userID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "parce user UUID failed")
		}

		rctx := context.WithValue(ctx, config.CtxConfig{}, config.NewCtxConfig(userID, isNewUser))

		m, err := handler(rctx, req)
		if err != nil {
			zap.S().Infoln("RPC failed with error: %v", err)
		}
		return m, err
	}

	// UUID not set.
	userID, err := uuid.NewV7()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error generate user uuid")
	}
	// Request contex.
	rctx := context.WithValue(ctx, config.CtxConfig{}, config.NewCtxConfig(userID.String(), isNewUser))

	cUserID, err = service.EncodeCookie(userID.String(), pass)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error encode cookie")
	}
	// Add userID to response.
	md = metadata.New(map[string]string{"user_id": cUserID})
	ctx = metadata.NewOutgoingContext(rctx, md)
	m, err := handler(ctx, req)
	if err != nil {
		zap.S().Infoln("RPC failed with error: %v", err)
	}
	return m, err
}
