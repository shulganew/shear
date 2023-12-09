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
	conf       *config.App
}

func (u *HandlerURL) GetServiceURL() service.Shortener {
	return *u.serviceURL
}

// GET and redirect by brief
func (u *HandlerURL) GetURL(res http.ResponseWriter, req *http.Request) {
	brief := chi.URLParam(req, "id")

	//get long Url from storage
	origin, exist := u.serviceURL.GetOrigin(req.Context(), brief)

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	if exist {
		res.Header().Set("Location", origin)
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

	brief, answerURL := u.serviceURL.GetAnsURL(redirectURL.Scheme, u.conf.Response)

	//save map to storage
	u.serviceURL.SetURL(req.Context(), brief, (*redirectURL).String())

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	//set status code 201
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(answerURL.String()))
}

func NewHandlerWeb(configApp *config.App) *HandlerURL {

	return &HandlerURL{serviceURL: service.NewService(configApp.Storage, configApp.Backup), conf: configApp}
}
