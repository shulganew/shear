package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/shulganew/shear.git/internal/entities"
	"go.uber.org/zap"
)

// Length of short URL random sequence length.
const ShortLength = 8

type Shortener struct {
	storeURLs StorageURL
}

// Interface for universal data storage, witch contain short and full URL.
type StorageURL interface {
	Set(ctx context.Context, userID string, brief, origin string) error
	SetAll(ctx context.Context, short []entities.Short) error
	GetOrigin(ctx context.Context, brief string) (string, bool, bool)
	GetBrief(ctx context.Context, origin string) (string, bool, bool)
	GetAll(ctx context.Context) []entities.Short
	GetUserAll(ctx context.Context, userID string) []entities.Short
	DeleteBatch(ctx context.Context, userID string, briefs []string)
}

// Service constructor.
func NewService(storage StorageURL) *Shortener {
	return &Shortener{storeURLs: storage}
}

func (s *Shortener) SetURL(ctx context.Context, userID, brief, origin string) (err error) {
	err = s.storeURLs.Set(ctx, userID, brief, origin)
	if err != nil {
		return err
	}
	return nil
}

func (s *Shortener) SetAll(ctx context.Context, short []entities.Short) (err error) {
	err = s.storeURLs.SetAll(ctx, short)
	if err != nil {
		return fmt.Errorf("error during save URL to Store: %w", err)
	}
	return nil
}

func (s *Shortener) GetOrigin(ctx context.Context, brief string) (origin string, exist bool, isDeleted bool) {
	return s.storeURLs.GetOrigin(ctx, brief)
}

func (s *Shortener) GetBrief(ctx context.Context, origin string) (brief string, exist bool, isDeleted bool) {
	return s.storeURLs.GetBrief(ctx, origin)
}

func (s *Shortener) GetUserAll(ctx context.Context, userID string) (short []entities.Short) {
	return s.storeURLs.GetUserAll(ctx, userID)
}

// Delete
func (s *Shortener) DeleteBatchArray(ctx context.Context, delBatchs []DelBatch) {
	for _, del := range delBatchs {
		s.storeURLs.DeleteBatch(ctx, del.UserID, del.Briefs)
	}
}

func (s *Shortener) DeleteBatch(ctx context.Context, delBatch DelBatch) {
	s.storeURLs.DeleteBatch(ctx, delBatch.UserID, delBatch.Briefs)
}

// return answer url: "schema + response server address from config + brief"
func (s *Shortener) GetAnsURL(schema, resultaddr string, brief string) (mainURL string, answerURL *url.URL) {
	// main URL = Schema + hostname + port (from result add -flag cmd -b)
	mainURL = schema + "://" + resultaddr

	// join full long URL
	longStrURL, err := url.JoinPath(mainURL, brief)
	if err != nil {
		zap.S().Errorln("Error during JoinPath", err)
	}
	answerURL, err = url.Parse(longStrURL)
	if err != nil {
		zap.S().Errorln("Error during Parse URL", err)
	}
	return
}

// Generate short link.
func GenerateShortLink() string {
	b := []byte{97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90}
	s := make([]byte, ShortLength)
	for i := 0; i < ShortLength; i++ {
		s[i] = b[rand.Intn(len(b))]
	}
	return string(s)
}

// Crypto function for getting crypted user id from cookies.
func GetCodedUserID(req *http.Request, pass string) (userID string, ok bool) {
	// find UserID in cookies
	cookie, err := req.Cookie("user_id")
	if err != nil {
		return "", false
	}

	userID, err = DecodeCookie(cookie.Value, pass)
	if err != nil {
		return "", false
	}

	// check correct UUID
	_, err = uuid.Parse(userID)
	if err != nil {
		return "", false
	}
	return userID, true

}

// Crypto function - decode secret string with password.
func DecodeCookie(secret string, password string) (uuid string, err error) {
	nonce, aesgcm, err := getCryptData(password)
	if err != nil {
		zap.S().Errorln("Encription Error: get enctypt data")
		return
	}

	msg, err := base64.StdEncoding.DecodeString(secret)

	if err != nil {
		return "", err
	}

	binary := []byte(msg)

	decrypted, err := aesgcm.Open(nil, nonce, binary, nil)
	if err != nil {
		zap.S().Errorln("Encryption Error: Open seal")
		return
	}
	return string(decrypted), nil
}

// Crypto function - Encode secret string with password.
func EncodeCookie(uuid string, password string) (secret string, err error) {
	binary := []byte(uuid)
	nonce, aesgcm, err := getCryptData(password)
	if err != nil {
		zap.S().Errorln("Encription Error: get enctypt data")
		return
	}

	coded := aesgcm.Seal(nil, nonce, binary, nil)
	return base64.StdEncoding.EncodeToString(coded), nil
}

// Get nonce and cipher from string.
func getCryptData(password string) (nonce []byte, aesgcm cipher.AEAD, err error) {
	key := sha256.Sum256([]byte(password))

	aesblock, err := aes.NewCipher(key[:32])
	if err != nil {
		zap.S().Errorln("Encryption Error: aesblock")
		return
	}

	aesgcm, err = cipher.NewGCM(aesblock)
	if err != nil {
		zap.S().Errorln("Encryption Error: aesgcm")
		return
	}

	length := aesgcm.NonceSize()
	nonceSize := len(key) - length
	nonce = key[nonceSize:]
	return
}
