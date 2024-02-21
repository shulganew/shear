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
	"strings"

	"github.com/google/uuid"
	"github.com/shulganew/shear.git/internal/entities"
	"go.uber.org/zap"
)

const ShortLength = 8

// generate sort URL
// save short URL
// get logn URL
type Shortener struct {
	storeURLs StorageURL
}

// intarface for universal data storage
type StorageURL interface {
	Set(ctx context.Context, userID string, brief, origin string) error
	SetAll(ctx context.Context, short []entities.Short) error
	GetOrigin(ctx context.Context, brief string) (string, bool, bool)
	GetBrief(ctx context.Context, origin string) (string, bool, bool)
	GetAll(ctx context.Context) []entities.Short
	GetUserAll(ctx context.Context, userID string) []entities.Short
	DelelteBatch(ctx context.Context, userID string, briefs []string)
}

// return service
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

func (s *Shortener) DelelteBatchArray(ctx context.Context, delBatchs []DelBatch) {

	for _, del := range delBatchs {
		s.storeURLs.DelelteBatch(ctx, del.UserID, del.Briefs)
	}

}

func (s *Shortener) DelelteBatch(ctx context.Context, delBatch DelBatch) {

	s.storeURLs.DelelteBatch(ctx, delBatch.UserID, delBatch.Briefs)

}

// return anwwer url: "shema + respose server addres from config + brief"
func (s *Shortener) GetAnsURL(shema, resultaddr string, brief string) (mainURL string, answerURL *url.URL) {
	//main URL = Shema + hostname + port (from result add -flag cmd -b)
	mainURL = shema + "://" + resultaddr

	//join full long URL
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

// generate short link
func GenerateShorLink() string {

	//base charset
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//nuber of short chars in url string

	sb := strings.Builder{}
	sb.Grow(ShortLength)
	for i := 0; i < ShortLength; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

//Crypto functions

func GetCodedUserID(req *http.Request, pass string) (userID string, ok bool) {
	//find UserID in cookies
	cookie, err := req.Cookie("user_id")
	if err != nil {
		return "", false
	}

	userID, err = DecodeCookie(cookie.Value, pass)
	if err != nil {
		return "", false
	}

	//check correct UUID
	_, err = uuid.Parse(userID)
	if err != nil {
		return "", false
	}
	return userID, true

}

func DecodeCookie(secret string, password string) (uuid string, err error) {

	nonce, aesgcm, err := GetCryptData(password)
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
		zap.S().Errorln("Encription Error: Open seal")
		return
	}
	return string(decrypted), nil

}

func EncodeCookie(uuid string, password string) (secret string, err error) {

	binary := []byte(uuid)

	nonce, aesgcm, err := GetCryptData(password)
	if err != nil {
		zap.S().Errorln("Encription Error: get enctypt data")
		return
	}

	coded := aesgcm.Seal(nil, nonce, binary, nil)

	return base64.StdEncoding.EncodeToString(coded), nil

}

func GetCryptData(password string) (nonce []byte, aesgcm cipher.AEAD, err error) {

	key := sha256.Sum256([]byte(password))

	aesblock, err := aes.NewCipher(key[:32])
	if err != nil {
		zap.S().Errorln("Encription Error: aesblock")
		return
	}

	aesgcm, err = cipher.NewGCM(aesblock)
	if err != nil {
		zap.S().Errorln("Encription Error: aesgcm")
		return
	}

	lenth := aesgcm.NonceSize()
	nonceSize := len(key) - lenth
	nonce = key[nonceSize:]
	return

}
