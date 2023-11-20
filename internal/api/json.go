package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

const CONTENT_TYPE_JSON = "application/json"

type RequestJSON struct {
	LongURL string `json:"url"`
}

type ResonseJSON struct {
	ShortURL string `json:"result"`
}

type HandlerAPI struct {
	serviceURL *service.Shortener
	conf       *config.Shear
	logz       zap.SugaredLogger
}

func (u *HandlerAPI) GetService() service.Shortener {
	return *u.serviceURL
}

func (u *HandlerAPI) SetAPI(res http.ResponseWriter, req *http.Request) {

	var request RequestJSON

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	u.logz.Infoln("Request long URL from JSON: ", request.LongURL)
	longURL, err := url.Parse(string(request.LongURL))
	if err != nil {
		http.Error(res, "Wrong URL in JSON, parse error", http.StatusInternalServerError)
	}

	shortURL, answerURL := u.serviceURL.GetAnsURL(longURL.Scheme, u.conf.ResultAddress)

	//save map to storage
	u.serviceURL.SetURL(shortURL, *longURL)

	response := ResonseJSON{answerURL.String()}

	jsonURL, err := json.Marshal(response)
	if err != nil {
		http.Error(res, "Error during Marshal answer URL", http.StatusInternalServerError)
	}
	u.logz.Infoln("Server ansver with short URL in JSON: ", string(jsonURL))

	//set content type
	res.Header().Add("Content-Type", CONTENT_TYPE_JSON)

	//set status code 201
	res.WriteHeader(http.StatusCreated)

	res.Write(jsonURL)

}

func NewHandler(configApp *config.Shear) *HandlerAPI {

	return &HandlerAPI{serviceURL: service.NewService(configApp.Storage), conf: configApp, logz: configApp.Applog}
}
