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

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
func TestGRPCGet(t *testing.T) {

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
		request           string
		body              string
		origin            string
		method            string
		contentType       string
		brief             string
		statusError       error
		responseExist     bool
		responseIsDeleted bool
	}{
		{
			name:              "Add URL gRPC Ok",
			request:           "http://localhost:8080",
			body:              "http://yandex.ru/",
			origin:            "http://yandex.ru/",
			contentType:       "text/plain",
			statusError:       nil,
			brief:             "qwerqewr",
			responseExist:     true,
			responseIsDeleted: false,
		},
		{
			name:              "Add URL gRPC Deleted",
			request:           "http://localhost:8080",
			body:              "http://yandex.ru/",
			origin:            "http://yandex.ru/",
			contentType:       "text/plain",
			statusError:       status.Errorf(codes.Unknown, "Deleted: StatusGone"),
			brief:             "qwerqewr",
			responseExist:     true,
			responseIsDeleted: true,
		},
		{
			name:              "Add URL gRPC Deleted",
			request:           "http://localhost:8080",
			body:              "http://yandex.ru/",
			origin:            "http://yandex.ru/",
			contentType:       "text/plain",
			statusError:       status.Errorf(codes.NotFound, "NotFound"),
			brief:             "qwerqewr",
			responseExist:     false,
			responseIsDeleted: false,
		},
	}
	for _, tt := range tests {

		_ = storeMock.EXPECT().
			GetOrigin(gomock.Any(), tt.brief).
			Times(1).
			Return(tt.origin, tt.responseExist, tt.responseIsDeleted)

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
		assert.NoError(t, err)
		defer func() {
			err := conn.Close()
			assert.NoError(t, err)
		}()

		client := pb.NewUsersClient(conn)
		_, err = client.GetURL(ctx, &pb.GetURLRequest{Brief: tt.brief})

		assert.Equal(t, tt.statusError, err)
	}
}
