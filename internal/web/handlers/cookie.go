package handlers

import (
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

// hadler for  GET and POST  short and long urls

type HandlerCookie struct {
	serviceURL *service.Shortener
	conf       *config.Config
}

func NewHandlerCookieID(conf *config.Config, stor *service.StorageURL) *HandlerCookie {

	return &HandlerCookie{serviceURL: service.NewService(stor), conf: conf}
}

func (u *HandlerCookie) GetServiceURL() service.Shortener {
	return *u.serviceURL
}

// GET and redirect by brief
func (u *HandlerCookie) GetUserURLs(res http.ResponseWriter, req *http.Request) {

	zap.S().Infoln("Hello from url api!")
	cookie, err := req.Cookie("user_id")
	if err != nil {
		//zap.S().Errorln("Cookie not found", err)
	}
	zap.S().Infoln("Cookie!!!!!!", cookie, len(req.Cookies()))

	//get long Url from storage

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	res.WriteHeader(http.StatusNotFound)

}
