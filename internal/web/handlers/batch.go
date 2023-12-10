package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

type BatchRequest struct {
	SessionID int    `json:"correlation_id"`
	Origin    string `json:"original_url"`
}

type BatchResonse struct {
	SessionID int    `json:"correlation_id"`
	Answer    string `json:"short_url"`
}

type HandlerBatch struct {
	serviceURL *service.Shortener
	conf       *config.App
}

func (u *HandlerBatch) GetService() service.Shortener {
	return *u.serviceURL
}

func (u *HandlerBatch) Batch(res http.ResponseWriter, req *http.Request) {

	//handle bach request
	var request []BatchRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		zap.S().Errorln("Get batch: ", err)

		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	batches := []BatchResonse{}
	shorts := []storage.Short{}
	for i := 0; i < len(request); i++ {

		origin, err := url.Parse(string(request[i].Origin))
		if err != nil {
			http.Error(res, "Wrong URL in JSON, parse error", http.StatusInternalServerError)
		}

		//get short brief and full answer URL
		brief, answerURL := u.serviceURL.GetAnsURL(origin.Scheme, u.conf.Response)
		//get batch for answer
		batch := BatchResonse{SessionID: request[i].SessionID, Answer: answerURL.String()}
		//add batches
		batches = append(batches, batch)
		//add short
		short := storage.Short{Brief: brief, Origin: (*origin).String()}
		shorts = append(shorts, short)

	}
	//save to storage
	u.serviceURL.SetAll(req.Context(), shorts)

	jsonBatch, err := json.Marshal(batches)
	if err != nil {
		http.Error(res, "Error during Marshal answer URL", http.StatusInternalServerError)
	}

	//set content type
	res.Header().Add("Content-Type", "application/json")

	//set status code 201
	res.WriteHeader(http.StatusCreated)

	res.Write(jsonBatch)

}

func NewHandlerBatch(configApp *config.App) *HandlerBatch {

	return &HandlerBatch{serviceURL: service.NewService(configApp.Storage, configApp.Backup), conf: configApp}
}
