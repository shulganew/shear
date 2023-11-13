package storage

import "net/url"

var urldb map[string]url.URL = make(map[string]url.URL)

func GetURLdb() *map[string]url.URL {

	return &urldb
}
