package service

import (
	"log"
	"net/url"

	"github.com/shulganew/shear.git/internal/shortener"
	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

// generate sort URL
// save short URL
// get logn URL
type Shortener struct {
	storeURLs storage.StorageURL
	backup    Backup
}

func (s *Shortener) SetURL(sortURL, longURL string) {
	zap.S().Infof("Store. Save URL [%s]=%s", sortURL, longURL)

	short := s.storeURLs.SetURL(sortURL, longURL)
	//save Short if backup is enable
	if s.backup.IsActive {
		s.backup.Save(short)
	}

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

	shortURL = shortener.GenerateShorLink()

	//join full long URL
	longStrURL, _ := url.JoinPath(mainURL, shortURL)
	answerURL, _ = url.Parse(longStrURL)

	log.Println("Save long url: ", answerURL)
	return
}

// return service
func NewService(storage storage.StorageURL, backup Backup) *Shortener {
	return &Shortener{storeURLs: storage, backup: backup}
}
