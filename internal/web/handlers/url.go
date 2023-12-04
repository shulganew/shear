package handlers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
)

// hadler for  GET and POST  short and long urls

type HandlerURL struct {
	serviceURL *service.Shortener
	conf       *config.Shear
}

func (u *HandlerURL) GetServiceURL() service.Shortener {
	return *u.serviceURL
}

// GET and redirect by shortUrl
func (u *HandlerURL) GetURL(res http.ResponseWriter, req *http.Request) {
	shortURL := chi.URLParam(req, "id")

	//get long Url from storage
	longURL, exist := u.serviceURL.GetLongURL(shortURL)

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	if exist {
		res.Header().Set("Location", longURL)
		//set status code 307
		res.WriteHeader(http.StatusTemporaryRedirect)

		return
	}

	res.WriteHeader(http.StatusNotFound)

}

// POTS and set generate short Url
func (u *HandlerURL) SetURL(res http.ResponseWriter, req *http.Request) {

	readBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Body not found", http.StatusInternalServerError)
	}

	redirectURL, err := url.Parse(string(readBody))
	if err != nil {
		http.Error(res, "Wrong URL in body, parse error", http.StatusInternalServerError)
	}

	shortURL, answerURL := u.serviceURL.GetAnsURL(redirectURL.Scheme, u.conf.Response)

	//save map to storage
	u.serviceURL.SetURL(shortURL, (*redirectURL).String())

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	//set status code 201
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(answerURL.String()))
}

func NewHandlerWeb(configApp *config.Shear) *HandlerURL {

	return &HandlerURL{serviceURL: service.NewService(configApp.Storage, configApp.Backup), conf: configApp}
}
