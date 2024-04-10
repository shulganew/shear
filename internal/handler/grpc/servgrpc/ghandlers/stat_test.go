package ghandlers

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
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

func TestGRPCStat(t *testing.T) {
	// Buffer for gRPC connection emulation.
	bufSize := 1024 * 1024
	var lis *bufconn.Listener
	// init configApp
	configApp := config.DefaultConfig(false)

	tests := []struct {
		name            string
		origin          string
		brief           string
		statusError     error
		dbErrorShorts   error
		dbErrorUsers    error
		numShorts       int
		numUsers        int
		ipportSource    string
		ipportConfig    string
		mockCallsShorts int
		mockCallsUsers  int
	}{
		
			{
				name:            "Stat gRPC Ok",
				origin:          "http://yandex.ru/",
				statusError:     nil,
				dbErrorShorts:   nil,
				dbErrorUsers:    nil,
				brief:           "qwerqewr",
				numShorts:       10,
				numUsers:        20,
				ipportSource:    "192.168.2.33:4321",
				ipportConfig:    "192.168.2.1/24",
				mockCallsShorts: 1,
				mockCallsUsers:  1,
			},
			{
				name:            "Stat denay ip",
				origin:          "http://yandex.ru/",
				statusError:     status.Errorf(codes.PermissionDenied, "Source ip denay"),
				dbErrorShorts:   nil,
				dbErrorUsers:    nil,
				brief:           "qwerqewr",
				numShorts:       10,
				numUsers:        20,
				ipportSource:    "192.168.2.33:4321",
				ipportConfig:    "192.168.3.1/24",
				mockCallsShorts: 0,
				mockCallsUsers:  0,
			},
			{
				name:            "Stat error source",
				origin:          "http://yandex.ru/",
				statusError:     status.Errorf(codes.Internal, "Can't parse split host and port."),
				dbErrorShorts:   nil,
				dbErrorUsers:    nil,
				brief:           "qwerqewr",
				numShorts:       10,
				numUsers:        20,
				ipportSource:    "192.168.2.33",
				ipportConfig:    "192.168.3.1/24",
				mockCallsShorts: 0,
				mockCallsUsers:  0,
			},
			{
				name:            "Stat error CIDR",
				origin:          "http://yandex.ru/",
				statusError:     status.Errorf(codes.Internal, "Can't parse CIDR IP form config."),
				dbErrorShorts:   nil,
				dbErrorUsers:    nil,
				brief:           "qwerqewr",
				numShorts:       10,
				numUsers:        20,
				ipportSource:    "192.168.2.33:4321",
				ipportConfig:    "192.168.3.1/48",
				mockCallsShorts: 0,
				mockCallsUsers:  0,
			},

			{
				name:            "Stat error CIDR",
				origin:          "http://yandex.ru/",
				statusError:     status.Errorf(codes.Internal, "Can't parse CIDR IP form config."),
				dbErrorShorts:   nil,
				dbErrorUsers:    nil,
				brief:           "qwerqewr",
				numShorts:       10,
				numUsers:        20,
				ipportSource:    "192.168.2.33:4321",
				ipportConfig:    "192.168.3.1/48",
				mockCallsShorts: 0,
				mockCallsUsers:  0,
			},
		
		{
			name:            "Db error Shorts",
			origin:          "http://yandex.ru/",
			statusError:     status.Errorf(codes.Internal, "Error during getting num of shorts (URLs)."),
			dbErrorShorts:   errors.New("Db error"),
			dbErrorUsers:    nil,
			brief:           "qwerqewr",
			numShorts:       10,
			numUsers:        20,
			ipportSource:    "192.168.2.33:4321",
			ipportConfig:    "192.168.2.1/24",
			mockCallsShorts: 1,
			mockCallsUsers:  0,
		},
		{
			name:            "Db error User",
			origin:          "http://yandex.ru/",
			statusError:     status.Errorf(codes.Internal, "Error during getting num of users."),
			dbErrorShorts:   nil,
			dbErrorUsers:    errors.New("Db error"),
			brief:           "qwerqewr",
			numShorts:       10,
			numUsers:        20,
			ipportSource:    "192.168.2.33:4321",
			ipportConfig:    "192.168.2.1/24",
			mockCallsShorts: 1,
			mockCallsUsers:  1,
		},
	}
	for _, tt := range tests {

		initCtx := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			ctx = context.WithValue(ctx, config.CtxIP{}, tt.ipportConfig)
			ctx = context.WithValue(ctx, config.CtxPassKey{}, configApp.GetPass())
			// Add IP addr.
			ip, _ := net.ResolveTCPAddr("tcp", tt.ipportSource)
			ctx = peer.NewContext(ctx, &peer.Peer{Addr: ip})
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

		t.Log("Test: ", tt.name)
		_ = storeMock.EXPECT().
			GetNumShorts(gomock.Any()).
			Times(tt.mockCallsShorts).
			Return(tt.numShorts, tt.dbErrorShorts)

		_ = storeMock.EXPECT().
			GetNumUsers(gomock.Any()).
			Times(tt.mockCallsUsers).
			Return(tt.numUsers, tt.dbErrorUsers)

		ctx := context.Background()

		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (c net.Conn, err error) {
			return lis.DialContext(ctx)
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
		assert.NoError(t, err)
		defer func() {
			err := conn.Close()
			assert.NoError(t, err)
		}()
		client := pb.NewUsersClient(conn)
		// Get request.
		_, err = client.GetStat(ctx, &pb.GetStatRequest{})
		// Check error.
		assert.Equal(t, tt.statusError, err)
	}
}
