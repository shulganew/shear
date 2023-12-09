package service

import (
	"context"
	"math/rand"
	"net/url"
	"strings"

	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

const ShortLength = 8

// generate sort URL
// save short URL
// get logn URL
type Shortener struct {
	storeURLs storage.StorageURL
	backup    Backup
}

func (s *Shortener) SetURL(ctx context.Context, brief, origin string) {
	zap.S().Infof("Store. Save URL [%s]=%s", brief, origin)
	s.storeURLs.Set(ctx, brief, origin)
}

func (s *Shortener) GetOrigin(ctx context.Context, brief string) (origin string, exist bool) {
	return s.storeURLs.GetOrigin(ctx, brief)
}

func (s *Shortener) GetBrief(ctx context.Context, origin string) (brief string, exist bool) {
	return s.storeURLs.GetBrief(ctx, origin)
}

// return anwwer url: "shema + respose server addres from config + brief"
func (s *Shortener) GetAnsURL(shema, resultaddr string) (brief string, answerURL *url.URL) {
	//main URL = Shema + hostname + port (from result add -flag cmd -b)
	mainURL := shema + "://" + resultaddr

	brief = GenerateShorLink()

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

// return service
func NewService(storage storage.StorageURL, backup Backup) *Shortener {
	return &Shortener{storeURLs: storage, backup: backup}
}
