package middlewares

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
	"go.uber.org/zap"
)

// Middleware function check network restriction.
func NetAccess(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("!!!!!!!!!!")
		// ip/mask from context config.
		trust := req.Context().Value(config.CtxIP{}).(string)
		_, nip, err := net.ParseCIDR(trust)
		if err != nil {
			zap.S().Error("Can't parse CIDR IP form config: ", err)
			ctx := context.WithValue(req.Context(), config.CtxAllow{}, false)
			h.ServeHTTP(res, req.WithContext(ctx))
			return
		}

		// Get IP from header.
		ipStr := req.Header.Get("X-Real-IP")
		ip := net.ParseIP(ipStr)
		if ip == nil {
			ctx := context.WithValue(req.Context(), config.CtxAllow{}, false)
			h.ServeHTTP(res, req.WithContext(ctx))
			return
		}

		// Check if ip is allow for CIDR IP/mask.
		isAllow := nip.Contains(ip)
		ctx := context.WithValue(req.Context(), config.CtxAllow{}, isAllow)
		h.ServeHTTP(res, req.WithContext(ctx))
	})

}
