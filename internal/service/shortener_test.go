package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShortener(t *testing.T) {
	tests := []struct {
		name              string
		request           string
		responseBrief     string
		responseExist     bool
		responseIsDeleted bool
	}{
		{
			name:              "base test POTS",
			request:           "http://yandex1.ru/",
			responseBrief:     "dzafbfsx",
			responseExist:     true,
			responseIsDeleted: false,
		},
	}

	// init configApp
	configApp := config.InitConfig()

	// init config with difauls values
	configApp.Address = config.DefaultHost
	configApp.Response = config.DefaultHost

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storeMock := mocks.NewMockStorageURL(ctrl)

			_ = storeMock.EXPECT().
				GetBrief(context.Background(), tt.request).
				Times(1).
				Return(tt.responseBrief, tt.responseExist, tt.responseIsDeleted)
			brief, isOk, isDel := storeMock.GetBrief(context.Background(), tt.request)
			assert.Equal(t, brief, tt.responseBrief)
			assert.Equal(t, isOk, tt.responseExist)
			assert.Equal(t, isDel, tt.responseIsDeleted)

		})

	}
}
func BenchmarkShortener(b *testing.B) {
	b.Run("generate short", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GenerateShortLinkByte()
		}
	})

	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// crete mock storege
	storeMock := mocks.NewMockStorageURL(ctrl)

	// init storage
	shortener := NewService(storeMock)

	_ = storeMock.EXPECT().
		GetOrigin(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return("yandex.ru", true, false)

	b.Run("get URL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			shortener.GetAnsURL("http", "localhost:8080", GenerateShortLinkByte())
		}

	})

	pass := "mypassword"
	b.Run("Encode and decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Init user's UUID.
			b.StopTimer()
			uuid, err := uuid.NewV7()
			assert.NoError(b, err)
			b.StartTimer()

			secret, err := EncodeCookie(uuid.String(), pass)
			assert.NoError(b, err)
			_, err = DecodeCookie(secret, pass)
			assert.NoError(b, err)
		}
	})

}
