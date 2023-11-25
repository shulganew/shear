package validators

import (
	"net"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

// Parse server and url address
func CheckURL(address string) (host string, port string) {

	zap.S().Infoln("Parse address: ", address)
	link, err := url.Parse(strings.TrimSpace(address))
	if err != nil {
		zap.S().Infoln("Error parsing url: ", err, " return def localhost:8080")
		return "localhost", "8080"
	}

	//check shema
	if link.Scheme != "http" {
		zap.S().Infoln("Shema not found, use http")
		address = "http://" + address
	}

	link, err = url.Parse(strings.TrimSpace(address))
	if err != nil {
		zap.S().Infoln("Error parsing url whis shema: ", err, " return def localhost:8080")
		return "localhost", "8080"
	}

	zap.S().Infoln("Split address: ", link)
	host, port, err2 := net.SplitHostPort(strings.TrimSpace(link.Host))
	if err2 != nil {
		zap.S().Infoln("Error split port: ", err, " return def localhost:8080 Host:", link.Host)
		return "", ""
	}
	return host, port
}