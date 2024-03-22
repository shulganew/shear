package service

import (
	"context"
	"database/sql"
	"net/http"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/entities"
	"github.com/shulganew/shear.git/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

const DefaultHost = "localhost:8080"

func TestShortenerURL(t *testing.T) {
	tests := []struct {
		name              string
		request           string
		responseBrief     string
		bathch            []string
		shorts            []entities.Short
		responseExist     bool
		responseIsDeleted bool
	}{
		{
			name:              "base test POTS",
			request:           "yandex1.ru/",
			responseBrief:     "dzafbfsx",
			responseExist:     true,
			responseIsDeleted: false,
			shorts:            []entities.Short{{ID: 0, UUID: sql.NullString{String: uuid.NewString(), Valid: true}, Brief: "dzafbfsx", Origin: "http://yandex1.ru/", SessionID: "qq12"}, {ID: 1, UUID: sql.NullString{String: uuid.NewString(), Valid: true}, Brief: "dzafbfsy", Origin: "http://yandex2.ru/", SessionID: "qq13"}},
			bathch:            []string{"sdfsdf", "sdfsdfsdf"},
		},
	}

	// init configApp
	configApp := config.NewConfig()

	// init config with difauls values
	configApp.SetAddress(DefaultHost)
	configApp.SetResponse(DefaultHost)
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

			_ = storeMock.EXPECT().
				GetAll(context.Background()).
				Times(1).
				Return(tt.shorts)

			// Test GetBrief.
			brief, isOk, isDel := storeMock.GetBrief(context.Background(), tt.request)
			assert.Equal(t, brief, tt.responseBrief)
			assert.Equal(t, isOk, tt.responseExist)
			assert.Equal(t, isDel, tt.responseIsDeleted)

			// Test GetALL.
			shortserv := NewService(storeMock)
			shorts := shortserv.GetAll(context.Background())
			assert.Equal(t, shorts, tt.shorts)

			mainURL, answerURL := shortserv.GetAnsURL("http", tt.request, tt.responseBrief)
			assert.Equal(t, mainURL, "http://"+tt.request)
			assert.Equal(t, answerURL.String(), "http://"+tt.request+brief)

			sortStr := GenerateShortLink()
			assert.Equal(t, len(sortStr), ShortLength)
			assert.True(t, regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(sortStr))
		})

	}
}

func TestShortenerDelButhc(t *testing.T) {
	tests := []struct {
		name   string
		bathch []string
	}{
		{
			name: "base test POTS",

			bathch: []string{"sdfsdf", "sdfsdfsdf"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storeMock := mocks.NewMockStorageURL(ctrl)

			userID, err := uuid.NewV7()
			assert.NoError(t, err)
			dbatch := DelBatch{UserID: userID.String(), Briefs: tt.bathch}

			_ = storeMock.EXPECT().
				DeleteBatch(context.Background(), userID.String(), tt.bathch).
				Times(2).
				Return(nil)

			// Test GetALL.
			shortserv := NewService(storeMock)
			err = shortserv.DeleteBatch(context.Background(), dbatch)
			assert.NoError(t, err)
			dbarray := []DelBatch{dbatch}
			shortserv.DeleteBatchArray(context.Background(), dbarray)
		})

	}
}

func TestShortenerCookie(t *testing.T) {
	t.Run("Cookie test", func(t *testing.T) {
		userID, err := uuid.NewV7()
		assert.NoError(t, err)
		pass := "myPass"

		cookie, err := EncodeCookie(userID.String(), pass)
		assert.NoError(t, err)

		decodedUserID, err := DecodeCookie(cookie, pass)
		assert.NoError(t, err)
		assert.Equal(t, userID.String(), decodedUserID)
		_, err = DecodeCookie(cookie, "error")
		assert.Error(t, err)

		_, ok := GetCodedUserID(&http.Request{}, "test")
		assert.False(t, ok)

	})

}

func BenchmarkShortener(b *testing.B) {

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

	b.Run("get URL Fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			shortener.GetAnsURLFast("http", "localhost:8080", GenerateShortLinkByte())
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

func BenchmarkGenerateShort(b *testing.B) {
	b.Run("generate short string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GenerateShortLink()
		}
	})

	b.Run("generate short byte", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GenerateShortLinkByte()
		}
	})

}
