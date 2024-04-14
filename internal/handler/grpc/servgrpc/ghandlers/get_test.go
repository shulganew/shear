package ghandlers

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
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

func TestGRPCGet(t *testing.T) {
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
	us := NewUsersServer(serviceURL, &configApp, nil, nil)
	pb.RegisterUsersServer(s, us)

	lis = bufconn.Listen(bufSize)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	tests := []struct {
		name              string
		origin            string
		brief             string
		statusError       error
		responseExist     bool
		responseIsDeleted bool
	}{
		{
			name:              "Get URL gRPC Ok",
			origin:            "http://yandex.ru/",
			statusError:       nil,
			brief:             "qwerqewr",
			responseExist:     true,
			responseIsDeleted: false,
		},
		{
			name:              "Get URL gRPC Deleted",
			origin:            "http://yandex.ru/",
			statusError:       status.Errorf(codes.Unknown, "Deleted: StatusGone"),
			brief:             "qwerqewr",
			responseExist:     true,
			responseIsDeleted: true,
		},
		{
			name:              "Get URL gRPC Deleted",
			origin:            "http://yandex.ru/",
			statusError:       status.Errorf(codes.NotFound, "NotFound"),
			brief:             "qwerqewr",
			responseExist:     false,
			responseIsDeleted: false,
		},
	}
	for _, tt := range tests {
		t.Log("Test: ", tt.name)
		_ = storeMock.EXPECT().
			GetOrigin(gomock.Any(), tt.brief).
			Times(1).
			Return(tt.origin, tt.responseExist, tt.responseIsDeleted)

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
		assert.NoError(t, err)
		defer func() {
			err := conn.Close()
			assert.NoError(t, err)
		}()
		client := pb.NewUsersClient(conn)
		// Get request.
		_, err = client.GetURL(ctx, &pb.GetURLRequest{Brief: tt.brief})
		// Check error.
		assert.Equal(t, tt.statusError, err)
	}
}
