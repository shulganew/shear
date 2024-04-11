package rest

import (
	"encoding/json"
	"net/http"

	"github.com/shulganew/shear.git/internal/builders"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

// DTO JSON request.
type Request struct {
	URL string `json:"url"`
}

// DTO JSON response.
type Response struct {
	Brief string `json:"result"`
}

// Handler for API:
//
//	Post "/api/shorten"
type HandlerAPI struct {
	serviceURL *service.Shorten
	conf       *config.Config
}

// Service constructor.
func NewHandlerAPI(conf *config.Config, short *service.Shorten) *HandlerAPI {
	return &HandlerAPI{serviceURL: short, conf: conf}
}

// @Summary      API add URL in JSON
// @Description  Add origin URL by JSON request, get brief URL in response.
// @Tags         api
// @Accept       json
// @Produce      json
// @Success      201 "Created"
// @Failure      401 "User unauthorized"
// @Failure      500 "Handling error"
// @Router       /api/shorten [post]
func (u *HandlerAPI) GetBrief(res http.ResponseWriter, req *http.Request) {
	// get UserID from cxt values
	ctxConfig := req.Context().Value(config.CtxConfig{}).(config.CtxConfig)

	var request Request
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Get resDTO from service.
	resDTO := u.serviceURL.AddURL(req.Context(), builders.AddRequestDTO{Origin: request.URL, CtxConfig: ctxConfig, Resp: u.conf.GetResponse()})

	response := Response{Brief: resDTO.AnwerURL}

	jsonURL, err := json.Marshal(response)
	if err != nil {
		zap.S().Errorln("Marshal URL error: ", err)
		resDTO.Status.SetStatusREST(http.StatusInternalServerError)
	}
	zap.S().Infoln("Requset URL: ", request.URL, "Ans ULR: ", string(jsonURL), "  Status: ", resDTO.Status.GetStatusREST())

	// set content type
	res.Header().Add("Content-Type", "application/json")

	// Set status code
	res.WriteHeader(resDTO.Status.GetStatusREST())
	// Send generate and saved string.
	res.Write([]byte(jsonURL))
}
