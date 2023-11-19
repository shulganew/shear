package service

import (
	"log"
	"net/url"

	"github.com/shulganew/shear.git/internal/shortener"
	"github.com/shulganew/shear.git/internal/storage"
)

// generate sort URL
// save short URL
// get logn URL
type Shortener struct {
	StoreURLs storage.StorageURL
}

func (s *Shortener) SetURL(sortURL string, longURL url.URL) {
	log.Printf("Store. Save URL [%s]=%s", sortURL, &longURL)
	s.StoreURLs.SetURL(sortURL, longURL)
}

func (s *Shortener) GetLongURL(sortURL string) (longURL url.URL, exist bool) {
	return s.StoreURLs.GetLongURL(sortURL)
}

func (s *Shortener) GetShortURL(longURL string) (shortURL string, exist bool) {
	return s.StoreURLs.GetShortURL(longURL)
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
func NewService(storage storage.StorageURL) *Shortener {
	return &Shortener{StoreURLs: storage}
}
