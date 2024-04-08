package ghandlers

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shulganew/shear.git/internal/app"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/entities"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"github.com/shulganew/shear.git/internal/handler/grpc/servgrpc/interceptors"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/service/mocks"
	"github.com/shulganew/shear.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestGRPCBatch(t *testing.T) {
	app.InitLog()
	// Buffer for gRPC connection emulation.
	bufSize := 1024 * 1024
	var lis *bufconn.Listener
	// init configApp
	configApp := config.DefaultConfig(false)

	initCtx := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
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
		origins           []string
		userID            string
		statusError       error
		responseIsDeleted bool
		mockCalls         int
		errDB             error
		isDuplicated      bool
		userAuth          bool
	}{
		{
			name:              "Batch URL gRPC Ok",
			origins:           []string{"http://yandex1.ru/", "http://yandex1.ru/", "http://yandex3.ru/"},
			statusError:       nil,
			userID:            "018dea9b-7085-75f5-91c5-2ba674052348",
			responseIsDeleted: false,
			mockCalls:         1,
			errDB:             nil,
			isDuplicated:      false,
			userAuth:          true,
		},
		{
			name:              "Batch auth not ok",
			origins:           []string{"http://yandex1.ru/", "http://yandex1.ru/", "http://yandex3.ru/"},
			statusError:       status.Errorf(codes.PermissionDenied, "User not athorized"),
			userID:            "018dea9b-7085-75f5-91c5-2ba674052348",
			responseIsDeleted: false,
			mockCalls:         0,
			errDB:             nil,
			isDuplicated:      false,
			userAuth:          false,
		},
		{
			name:              "Batch Duplicated URL",
			origins:           []string{"http://yandex1.ru/", "http://yandex1.ru/", "http://yandex3.ru/"},
			statusError:       status.Errorf(codes.AlreadyExists, "Has existed original URL"),
			userID:            "018dea9b-7085-75f5-91c5-2ba674052348",
			responseIsDeleted: false,
			mockCalls:         1,
			errDB:             &storage.ErrDuplicatedShort{Err: errors.New("Database duplicated"), Label: "Duplicated", Short: entities.Short{}},
			isDuplicated:      true,
			userAuth:          true,
		},
		{
			name:              "Database error",
			origins:           []string{"http://yandex1.ru/", "http://yandex1.ru/", "http://yandex3.ru/"},
			statusError:       status.Errorf(codes.Internal, "Error saving in Storage"),
			userID:            "018dea9b-7085-75f5-91c5-2ba674052348",
			responseIsDeleted: false,
			mockCalls:         1,
			errDB:             errors.New("Database error"),
			isDuplicated:      true,
			userAuth:          true,
		},
	}
	for _, tt := range tests {
		t.Log("Test: ", tt.name)
		_ = storeMock.EXPECT().
			AddAll(gomock.Any(), gomock.Any()).
			Times(tt.mockCalls).
			Return(tt.errDB)

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

		// Get request.
		_, err = client.Batch(ctx, &pb.BatchRequest{Origins: tt.origins})

		// Check error.
		assert.Equal(t, tt.statusError, err)
	}
}