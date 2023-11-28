package storage

import "slices"

// base stract for working with storage
type Short struct {
	ID       int    `json:"uuid"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}

type StorageURL interface {
	SetURL(sortURL, longURL string) Short
	GetLongURL(sortURL string) (string, bool)
	GetShortURL(longURL string) (string, bool)
}

type MemoryStorage struct {
	StoreURLs []Short
}

func (m *MemoryStorage) SetURL(sortURL, longURL string) (short Short) {
	//init storage
	short = Short{ID: len(m.StoreURLs), ShortURL: sortURL, LongURL: longURL}
	m.StoreURLs = append(m.StoreURLs, short)
	return
}

func (m *MemoryStorage) GetLongURL(shortURL string) (longURL string, ok bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s Short) bool { return s.ShortURL == shortURL })
	if id != -1 {
		longURL = m.StoreURLs[id].LongURL
		ok = true
	}
	return
}

func (m *MemoryStorage) GetShortURL(longURL string) (shortURL string, ok bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s Short) bool { return s.LongURL == longURL })
	if id != -1 {
		shortURL = m.StoreURLs[id].ShortURL
		ok = true
	}
	return

}
