package utils

import (
	"math/rand"
	"net/url"
	"strings"
)

// generate short link
func GenerateShorLink() string {

	//base charset
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//nuber of short chars in url string
	n := 8

	sb := strings.Builder{}
	sb.Grow(7)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

// get shortUrl from BDUrl
func GetShortUrl(m *map[string]url.URL, longUrl string) (shortUrl string, ok bool) {
	for k, v := range *m {
		if v.String() == longUrl {
			shortUrl = k
			ok = true
			return
		}
	}
	return
}
