package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

type BatchRequest struct {
	SessionID string `json:"correlation_id"`
	Origin    string `json:"original_url"`
}

type BatchResonse struct {
	SessionID string `json:"correlation_id"`
	Answer    string `json:"short_url"`
}

type HandlerBatch struct {
	serviceURL *service.Shortener
	conf       *config.Config
}

func NewHandlerBatch(conf *config.Config, stor *service.StorageURL) *HandlerBatch {

	return &HandlerBatch{serviceURL: service.NewService(stor), conf: conf}
}

func (u *HandlerBatch) GetService() service.Shortener {
	return *u.serviceURL
}

func (u *HandlerBatch) BatchSet(res http.ResponseWriter, req *http.Request) {

	//handle bach requests
	var requests []BatchRequest

	if err := json.NewDecoder(req.Body).Decode(&requests); err != nil {
		zap.S().Errorln("Get batch: ", err)

		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	batches := []BatchResonse{}
	shorts := []service.Short{}
	for _, r := range requests {

		origin, err := url.Parse(string(r.Origin))
		if err != nil {
			http.Error(res, "Wrong URL in JSON, parse error", http.StatusInternalServerError)
		}

		//get short brief and full answer URL
		brief, _, answerURL := u.serviceURL.GetAnsURL(origin.Scheme, u.conf.Response)
		//get batch for answer
		batch := BatchResonse{SessionID: r.SessionID, Answer: answerURL.String()}
		//add batches
		batches = append(batches, batch)

		if err != nil {
			zap.S().Errorln("Brocken sessionID", batch.SessionID)
		}

		shortSession := service.Short{Brief: brief, Origin: (*origin).String(), SessionID: batch.SessionID}
		shorts = append(shorts, shortSession)

	}
	//save to storage
	err := u.serviceURL.SetAll(req.Context(), shorts)

	//check duplicated strings
	var tagErr *storage.ErrDuplicatedShort
	if err != nil {
		if errors.As(err, &tagErr) {
			//set status code 409 Conflict
			res.WriteHeader(http.StatusConflict)

			//send existed URL to response
			broken := []BatchResonse{}

			batch := BatchResonse{SessionID: tagErr.Short.SessionID, Answer: tagErr.Short.Brief}
			broken = append(broken, batch)
			jsonBrokenBatch, err := json.Marshal(broken)
			if err != nil {
				http.Error(res, "Error during Marshal answer URL", http.StatusInternalServerError)
			}

			//set content type
			res.Header().Add("Content-Type", "application/json")

			zap.S().Infoln("Broken: ", string(jsonBrokenBatch))
			res.Write(jsonBrokenBatch)
			return
		}
	}
	//create Ok answer
	jsonBatch, err := json.Marshal(batches)
	if err != nil {
		http.Error(res, "Error during Marshal answer URL", http.StatusInternalServerError)
	}

	zap.S().Infoln("Batch saved size: ", len(batches))
	//set content type
	res.Header().Add("Content-Type", "application/json")
	//set status code 201
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonBatch)

}
