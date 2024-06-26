package ghandlers

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/golang/mock/gomock"
	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"github.com/stretchr/testify/assert"

	"github.com/shulganew/shear.git/internal/handler/grpc/servgrpc/interceptors"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/service/mocks"
)

func TestGRPCDel(t *testing.T) {
	// Buffer for gRPC connection emulation.
	bufSize := 1024 * 1024
	var lis *bufconn.Listener
	// init configApp
	configApp := config.DefaultConfig(false)

	initCtx := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = context.WithValue(ctx, config.CtxIP{}, configApp.GetIP())
		ctx = context.WithValue(ctx, config.CtxPassKey{}, configApp.GetPass())
		return handler(ctx, req)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// crete mock storege
	storeMock := mocks.NewMockStorageURL(ctrl)
	serviceURL := service.NewService(storeMock)

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(initCtx, interceptors.AuthInterceptor))

	// Create Del service.
	delCh := make(chan service.DelBatch)
	del := service.NewDelete(delCh, &configApp)
	service.DeleteShort(context.Background(), serviceURL, delCh)

	// Register gRPC server.
	us := NewUsersServer(serviceURL, &configApp, nil, del)
	pb.RegisterUsersServer(s, us)

	lis = bufconn.Listen(bufSize)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	tests := []struct {
		name        string
		userID      string
		briefs      []string
		statusError error
		userAuth    bool
		mockCalls   int
	}{
		{
			name:        "Del URL gRPC Ok",
			userID:      "018dea9b-7085-75f5-91c5-2ba674052348",
			briefs:      []string{"afdgad", "sdfasgf"},
			statusError: nil,
			userAuth:    true,
			mockCalls:   1,
		},
		{
			name:        "Del URL gRPC Ok",
			userID:      "018dea9b-7085-75f5-91c5-2ba674052348",
			briefs:      []string{"afdgad", "sdfasgf"},
			statusError: status.Errorf(codes.PermissionDenied, "User not athorized"),
			userAuth:    false,
			mockCalls:   0,
		},
	}
	for _, tt := range tests {
		t.Log("Test: ", tt.name)
		_ = storeMock.EXPECT().
			DeleteBatch(gomock.Any(), tt.userID, tt.briefs).
			Times(tt.mockCalls).
			Return(nil)

		// Add MD to request.
		ctx := context.Background()
		if tt.userAuth {
			cUserID, _ := service.EncodeCookie(tt.userID, "mypass")
			md := metadata.New(map[string]string{"user_id": cUserID})
			ctx = metadata.NewOutgoingContext(ctx, md)
		}
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
		assert.NoError(t, err)
		defer func() {
			err := conn.Close()
			assert.NoError(t, err)
		}()
		client := pb.NewUsersClient(conn)

		// Del request.
		ok, err := client.DelUserURLs(ctx, &pb.DelRequest{Briefs: tt.briefs})
		if tt.userAuth {
			assert.True(t, ok.GetOk())
		}
		// Check error.
		assert.Equal(t, tt.statusError, err)
	}
}
