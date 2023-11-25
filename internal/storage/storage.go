package storage

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
	for _, short := range m.StoreURLs {
		if short.ShortURL == shortURL {
			return short.LongURL, true
		}
	}

	return
}

func (m *MemoryStorage) GetShortURL(longURL string) (shortURL string, ok bool) {
	for _, short := range m.StoreURLs {
		if short.LongURL == longURL {
			return short.ShortURL, true
		}
	}
	return
}

// get shortUrl from BDUrl
// func GetShortURL(m *map[string]url.URL, longURL string) (shortURL string, ok bool) {
// 	for k, v := range *m {
// 		if v.String() == longURL {
// 			shortURL = k
// 			ok = true
// 			return
// 		}
// 	}
// 	return
// }
