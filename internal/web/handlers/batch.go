package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/entities"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

// Handler for API:
//
//	Post "/api/shorten/batch"
type HandlerBatch struct {
	serviceURL *service.Shorten
	conf       *config.Config
}

// Service constructor.
func NewHandlerBatch(conf *config.Config, stor service.StorageURL) *HandlerBatch {

	return &HandlerBatch{serviceURL: service.NewService(stor), conf: conf}
}

// @Summary      Set several user's URLs in body in JSON format
// @Description  Set json URLs
// @Tags         api
// @Accept       json
// @Produce      json
// @Success      201 {object}  []entities.BatchResponse  "Created"
// @Failure      400 "Error JSON Unmarshal"
// @Failure      401 "User unauthorized"
// @Failure      404 "Conflict. URL existed."
// @Failure      500 "Handling error"
// @Router       /api/shorten/batch [post]
func (u *HandlerBatch) BatchSet(res http.ResponseWriter, req *http.Request) {
	// find UserID in cookies
	userID, err := req.Cookie("user_id")
	if err != nil {
		http.Error(res, "Can't find user in cookies", http.StatusUnauthorized)
	}
	// handle bach requests
	var requests []entities.BatchRequest
	if err := json.NewDecoder(req.Body).Decode(&requests); err != nil {
		zap.S().Errorln("Get batch: ", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	batches := []entities.BatchResponse{}
	shorts := []entities.Short{}
	for i, r := range requests {

		origin, err := url.Parse(string(r.Origin))
		if err != nil {
			http.Error(res, "Wrong URL in JSON, parse error", http.StatusInternalServerError)
		}
		// get short brief and full answer URL
		brief := service.GenerateShortLinkByte()
		_, answerURL := u.serviceURL.GetAnsURL(origin.Scheme, u.conf.Response, brief)
		// get batch for answer
		batch := entities.BatchResponse{SessionID: r.SessionID, Answer: answerURL.String()}
		// add batches
		batches = append(batches, batch)

		if err != nil {
			zap.S().Errorln("Brocken sessionID", batch.SessionID)
		}
		shortSession := entities.NewShort(i, userID.Value, brief, (*origin).String(), batch.SessionID)
		shorts = append(shorts, *shortSession)

	}
	// save to storage
	err = u.serviceURL.SetAll(req.Context(), shorts)

	// check duplicated strings
	var tagErr *storage.ErrDuplicatedShort
	if err != nil {
		if errors.As(err, &tagErr) {
			// set status code 409 Conflict
			res.WriteHeader(http.StatusConflict)
			// send existed URL to response
			broken := []entities.BatchResponse{}
			batch := entities.BatchResponse{SessionID: tagErr.Short.SessionID, Answer: tagErr.Short.Brief}
			broken = append(broken, batch)
			jsonBrokenBatch, err := json.Marshal(broken)
			if err != nil {
				http.Error(res, "Error during Marshal answer URL", http.StatusInternalServerError)
			}

			// set content type
			res.Header().Add("Content-Type", "application/json")
			zap.S().Infoln("Broken: ", string(jsonBrokenBatch))
			res.Write(jsonBrokenBatch)
			return
		}
	}
	// create Ok answer
	jsonBatch, err := json.Marshal(batches)
	if err != nil {
		http.Error(res, "Error during Marshal answer URL", http.StatusInternalServerError)
	}

	zap.S().Infoln("Batch saved size: ", len(batches))
	// set content type
	res.Header().Add("Content-Type", "application/json")
	// set status code 201
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonBatch)
}
