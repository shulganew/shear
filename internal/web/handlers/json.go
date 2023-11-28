package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

type RequestJSON struct {
	LongURL string `json:"url"`
}

type ResonseJSON struct {
	ShortURL string `json:"result"`
}

type HandlerAPI struct {
	serviceURL *service.Shortener
	conf       *config.Shear
}

func (u *HandlerAPI) GetService() service.Shortener {
	return *u.serviceURL
}

func (u *HandlerAPI) GetShortURL(res http.ResponseWriter, req *http.Request) {

	var request RequestJSON

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	longURL, err := url.Parse(string(request.LongURL))
	if err != nil {
		http.Error(res, "Wrong URL in JSON, parse error", http.StatusInternalServerError)
	}

	shortURL, answerURL := u.serviceURL.GetAnsURL(longURL.Scheme, u.conf.Response)

	//save map to storage
	u.serviceURL.SetURL(shortURL, (*longURL).String())

	response := ResonseJSON{answerURL.String()}

	jsonURL, err := json.Marshal(response)
	if err != nil {
		http.Error(res, "Error during Marshal answer URL", http.StatusInternalServerError)
	}
	zap.S().Infoln("Server ansver with short URL in JSON: ", string(jsonURL))

	//set content type
	res.Header().Add("Content-Type", "application/json")

	//set status code 201
	res.WriteHeader(http.StatusCreated)

	res.Write(jsonURL)

}

func NewHandlerAPI(configApp *config.Shear) *HandlerAPI {

	return &HandlerAPI{serviceURL: service.NewService(configApp.Storage, configApp.Backup), conf: configApp}
}
