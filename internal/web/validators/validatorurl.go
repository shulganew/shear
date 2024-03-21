// Package validators has a function for URL validations.
package validators

import (
	"net"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

// Parse server and url address. If address empty or parsing error - return default URL: 8080/8443
func CheckURL(address string, isSequer bool) (host string, port string, isDefault bool) {

	if len(address) == 0 {
		if isSequer {
			return "localhost", "8443", true
		}
		return "localhost", "8080", true
	}

	// parse address
	link, err := url.Parse(strings.TrimSpace(address))
	if err != nil {
		zap.S().Info("Error parsing url: ", err, " return def localhost:8080/8443")
		if isSequer {
			return "localhost", "8443", true
		}
		return "localhost", "8080", true
	}

	// check shema
	if !(link.Scheme == "http" || link.Scheme == "https") {
		// shema not found, use http
		if isSequer {
			address = "https://" + address
		} else {
			address = "http://" + address
		}
	}

	link, err = url.Parse(strings.TrimSpace(address))
	if err != nil {
		zap.S().Info("Error parsing url: ", err, " return def localhost:8080/8443")
		if isSequer {
			return "localhost", "8443", true
		}
		return "localhost", "8080", true
	}

	// split address
	host, port, err = net.SplitHostPort(strings.TrimSpace(link.Host))
	if err != nil {
		zap.S().Errorln("Error split port: ", err, " return def localhost:8080 Host:", link.Host)
		if isSequer {
			return "localhost", "8443", true
		}
		return "localhost", "8080", true
	}
	return host, port, false
}
