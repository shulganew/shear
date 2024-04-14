package ghandlers

import (
	"context"
	"database/sql"
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
	"github.com/shulganew/shear.git/internal/entities"
	pb "github.com/shulganew/shear.git/internal/handler/grpc/proto"
	"github.com/stretchr/testify/assert"

	"github.com/shulganew/shear.git/internal/handler/grpc/servgrpc/interceptors"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/service/mocks"
)

func TestGRPCGetUSers(t *testing.T) {
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
		userID            string
		statusError       error
		responseExist     bool
		responseIsDeleted bool
		userAuth          bool
		result            []entities.Short
		mockCalls         int
	}{
		{
			name:              "GetAll URL gRPC Ok",
			origin:            "http://yandex.ru/",
			userID:            "018dea9b-7085-75f5-91c5-2ba674052348",
			statusError:       nil,
			brief:             "qwerqewr",
			responseExist:     true,
			responseIsDeleted: false,
			userAuth:          true,
			result:            []entities.Short{{Brief: "sdfsf", Origin: "http://ya.ru", SessionID: "", UserID: sql.NullString{String: "018dea9b-7085-75f5-91c5-2ba674052348", Valid: true}, IsDeleted: false}},
			mockCalls:         1,
		},
		{
			name:              "GetAll empty",
			origin:            "http://yandex.ru/",
			userID:            "018dea9b-7085-75f5-91c5-2ba674052348",
			statusError:       status.Errorf(codes.Unknown, "No contenet for user."),
			brief:             "qwerqewr",
			responseExist:     true,
			responseIsDeleted: false,
			userAuth:          true,
			result:            []entities.Short{},
			mockCalls:         1,
		},
		{
			name:              "GetAll User auth",
			origin:            "http://yandex.ru/",
			userID:            "018dea9b-7085-75f5-91c5-2ba674052348",
			statusError:       status.Errorf(codes.PermissionDenied, "User not athorized"),
			brief:             "qwerqewr",
			responseExist:     true,
			responseIsDeleted: false,
			userAuth:          false,
			result:            []entities.Short{{Brief: "sdfsf", Origin: "http://ya.ru", SessionID: "", UserID: sql.NullString{String: "018dea9b-7085-75f5-91c5-2ba674052348", Valid: true}, IsDeleted: false}},
			mockCalls:         0,
		},
	}
	for _, tt := range tests {
		t.Log("Test: ", tt.name)
		_ = storeMock.EXPECT().
			GetUserAll(gomock.Any(), gomock.Any()).
			Times(tt.mockCalls).
			Return(tt.result)

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
		_, err = client.GetUserURLs(ctx, &pb.GetURLs{})
		// Check error.
		assert.Equal(t, tt.statusError, err)
	}
}
