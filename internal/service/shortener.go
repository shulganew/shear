package service

import (
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

func (s *Shortener) SetURL(sortURL, longURL string) {
	zap.S().Infof("Store. Save URL [%s]=%s", sortURL, longURL)
	s.storeURLs.SetURL(sortURL, longURL)
}

func (s *Shortener) GetLongURL(sortURL string) (longURL string, exist bool) {
	return s.storeURLs.GetLongURL(sortURL)
}

func (s *Shortener) GetShortURL(longURL string) (shortURL string, exist bool) {
	return s.storeURLs.GetShortURL(longURL)
}

// return anwwer url: "shema + respose server addres from config + shortURL"
func (s *Shortener) GetAnsURL(shema, resultaddr string) (shortURL string, answerURL *url.URL) {
	//main URL = Shema + hostname + port (from result add -flag cmd -b)
	mainURL := shema + "://" + resultaddr

	shortURL = GenerateShorLink()

	//join full long URL
	longStrURL, err := url.JoinPath(mainURL, shortURL)
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
	sb.Grow(7)
	for i := 0; i < ShortLength; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

// return service
func NewService(storage storage.StorageURL, backup Backup) *Shortener {
	return &Shortener{storeURLs: storage, backup: backup}
}
