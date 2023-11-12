package utils

import (
	"log"
	"math/rand"
	"net"
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

// Parse server and url address
func CheckAddress(address string) (host string, port string) {
	log.Println("Parse address: ", address)
	link, err := url.Parse(strings.TrimSpace(address))
	if err != nil {
		log.Printf("Error parsing url: ", err, " return def localhost:8080")
		return "localhost", "8080"
	}

	//check shema
	if link.Scheme != "http" {
		log.Println("Shema not found, use http")
		address = "http://" + address
	}

	link, err = url.Parse(strings.TrimSpace(address))
	if err != nil {
		log.Printf("Error parsing url whis shema: ", err, " return def localhost:8080")
		return "localhost", "8080"
	}

	log.Println("Split address: ", link)
	host, port, err2 := net.SplitHostPort(strings.TrimSpace(link.Host))
	if err2 != nil {
		log.Printf("Error split port: ", err, " return def localhost:8080 Host:", link.Host)
		return "", ""
	}
	return host, port
}
