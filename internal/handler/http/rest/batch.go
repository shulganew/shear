package rest

import (
	"encoding/json"
	"net/http"

	"github.com/shulganew/shear.git/internal/builders"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/entities"
	"github.com/shulganew/shear.git/internal/service"
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
func NewHandlerBatch(conf *config.Config, short *service.Shorten) *HandlerBatch {

	return &HandlerBatch{serviceURL: short, conf: conf}
}

// @Summary      Add several user's URLs in body in JSON format
// @Description  Add json URLs. No User auth.
// @Tags         api
// @Accept       json
// @Produce      json
// @Success      201 {object}  []entities.BatchResponse  "Created"
// @Failure      400 "Error JSON Unmarshal"
// @Failure      401 "User unauthorized"
// @Failure      404 "Conflict. URL existed."
// @Failure      500 "Handling error"
// @Router       /api/shorten/batch [post]
func (u *HandlerBatch) BatchAdd(res http.ResponseWriter, req *http.Request) {
	// get UserID from cxt values
	ctxConfig := req.Context().Value(config.CtxConfig{}).(config.CtxConfig)

	// handle bach requests
	var requests []entities.BatchRequest
	if err := json.NewDecoder(req.Body).Decode(&requests); err != nil {
		zap.S().Errorln("Get batch: ", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to string array of origins.
	resDTO := u.serviceURL.AddBatch(req.Context(), builders.BatchRequestDTO{Origins: requests, CtxConfig: ctxConfig, Resp: u.conf.GetResponse()})

	// Check errors in service.
	if resDTO.Err != nil {
		http.Error(res, resDTO.Err.Error(), resDTO.Status.GetStatusREST())
		return
	}

	// create Ok answer
	jsonBatch, err := json.Marshal(resDTO.AnwerURLs)
	if err != nil {
		http.Error(res, "Error during Marshal answer URL", http.StatusInternalServerError)
		return
	}

	zap.S().Infoln("Batch saved size: ", len(resDTO.AnwerURLs))
	// set content type
	res.Header().Add("Content-Type", "application/json")
	// set status code 201
	res.WriteHeader(resDTO.Status.GetStatusREST())
	res.Write(jsonBatch)
}
