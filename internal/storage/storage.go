package storage

import (
	"log"
	"net/url"
)

type URLSetGet interface {
	SetURL(sortURL string, longURL url.URL)
	GetLongURL(sortURL string) (url.URL, bool)
	GetShortURL(longURL string) (string, bool)
}

type Storage struct {
	storage URLSetGet
}

type MapStorage struct {
	storeURLs map[string]url.URL
}

func (m *MapStorage) SetURL(sortURL string, longURL url.URL) {
	//init storage
	log.Printf("Store. Save URL [%s]=%s", sortURL, &longURL)
	m.storeURLs[sortURL] = longURL
}

func (m *MapStorage) GetLongURL(sortURL string) (longURL url.URL, exist bool) {
	longURL, exist = m.storeURLs[sortURL]
	return longURL, exist
}

func (m *MapStorage) GetShortURL(longURL string) (shortURL string, exist bool) {
	for k, v := range m.storeURLs {
		if v.String() == longURL {
			shortURL = k
			exist = true
			return
		}
	}
	return
}

func NewMapStorage() URLSetGet {
	return &MapStorage{storeURLs: make(map[string]url.URL)}
}

// get shortUrl from BDUrl
func GetShortURL(m *map[string]url.URL, longURL string) (shortURL string, ok bool) {
	for k, v := range *m {
		if v.String() == longURL {
			shortURL = k
			ok = true
			return
		}
	}
	return
}
