package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
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
