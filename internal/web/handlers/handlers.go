package handlers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

// hadler for  GET and POST  hor and log urls

type URLHandler struct {
	serviceURL *service.Shortener
	conf       *config.Shear
	logz       zap.SugaredLogger
}

func (u *URLHandler) GetServiceURL() service.Shortener {
	return *u.serviceURL
}

// GET and redirect by shortUrl
func (u *URLHandler) GetURL(res http.ResponseWriter, req *http.Request) {

	shortURL := chi.URLParam(req, "id")

	//get long Url from storage
	longURL, exist := u.serviceURL.GetLongURL(shortURL)

	//set content type
	res.Header().Add("Content-Type", "text/plain")
	u.logz.Infoln("Redirect to: ", longURL)

	if exist {
		res.Header().Set("Location", longURL.String())
		//set status code 307
		res.WriteHeader(http.StatusTemporaryRedirect)

		return
	}

	res.WriteHeader(http.StatusNotFound)

}

// POTS and set generate short Url
func (u *URLHandler) SetURL(res http.ResponseWriter, req *http.Request) {

	readBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Body not found", http.StatusInternalServerError)
	}

	redirectURL, err := url.Parse(string(readBody))
	if err != nil {
		http.Error(res, "Wrong URL in body, parse error", http.StatusInternalServerError)
	}

	shortURL, answerURL := u.serviceURL.GetAnsURL(redirectURL.Scheme, u.conf.ResultAddress)

	//save map to storage
	u.serviceURL.SetURL(shortURL, *redirectURL)

	u.logz.Infoln("Server ansver with short URL: ", answerURL)

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	//set status code 201
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(answerURL.String()))
}

func NewHandler(configApp *config.Shear) *URLHandler {

	return &URLHandler{serviceURL: service.NewService(configApp.Storage), conf: configApp, logz: configApp.Applog}
}
