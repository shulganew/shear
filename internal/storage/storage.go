package storage

import (
	"net/url"
)

type StorageURL interface {
	SetURL(sortURL string, longURL url.URL)
	GetLongURL(sortURL string) (url.URL, bool)
	GetShortURL(longURL string) (string, bool)
}

type MapStorage struct {
	StoreURLs map[string]url.URL
}

func (m *MapStorage) SetURL(sortURL string, longURL url.URL) {
	//init storage
	m.StoreURLs[sortURL] = longURL
}

func (m *MapStorage) GetLongURL(sortURL string) (longURL url.URL, exist bool) {
	longURL, exist = m.StoreURLs[sortURL]
	return longURL, exist
}

func (m *MapStorage) GetShortURL(longURL string) (shortURL string, exist bool) {
	for k, v := range m.StoreURLs {
		if v.String() == longURL {
			shortURL = k
			exist = true
			return
		}
	}
	return
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
