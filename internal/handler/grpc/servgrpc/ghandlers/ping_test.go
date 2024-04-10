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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/shulganew/shear.git/internal/config"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"github.com/stretchr/testify/assert"

	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/service/mocks"
)

func TestGRPCPing(t *testing.T) {
	// Buffer for gRPC connection emulation.
	bufSize := 1024 * 1024
	var lis *bufconn.Listener
	// init configApp
	configApp := config.DefaultConfig(false)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// crete mock storege
	storeMock := mocks.NewMockStorageURL(ctrl)
	serviceURL := service.NewService(storeMock)

	s := grpc.NewServer()
	sql, smock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
	// Register gRPC server.
	us := NewUsersServer(serviceURL, &configApp, sql, nil)
	pb.RegisterUsersServer(s, us)

	lis = bufconn.Listen(bufSize)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	tests := []struct {
		name        string
		ping        bool
		statusError error
	}{
		{
			name:        "DB ping test available",
			ping:        true,
			statusError: nil,
		},

		{
			name:        "DB ping test not available",
			ping:        false,
			statusError: status.Errorf(codes.Internal, "Database connection failed."),
		},
	}
	for _, tt := range tests {
		t.Log("Test: ", tt.name)
		if tt.ping {
			smock.ExpectPing()
		}

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
		_, err = client.Ping(ctx, &pb.PingRequest{})

		assert.Equal(t, tt.statusError, err)
	}
}
