package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

type ResponseAuth struct {
	Brief  string `json:"short_url"`
	Origin string `json:"original_url"`
}

// Handler for API:
//
//	Get "/api/user/urls"
type HandlerAuth struct {
	serviceURL *service.Shortener
	conf       *config.Config
}

func NewHandlerAuthUser(conf *config.Config, stor service.StorageURL) *HandlerAuth {

	return &HandlerAuth{serviceURL: service.NewService(stor), conf: conf}
}

// Return all users shorts.
// @Summary      Get all user's origin URLs
// @Description  Add origin URL by JSON request, get brief URL in response
// @Tags         api
// @Accept       json
// @Produce      json
// @Success      200 {object} []ResponseAuth "OK"
// @Failure      401 "User unauthorized"
// @Failure      500 "Handling error"
// @Router       /api/user/urls [get]
func (u HandlerAuth) GetUserURLs(res http.ResponseWriter, req *http.Request) {
	// get UserID from cxt values
	ctxConfig := req.Context().Value(config.CtxConfig{}).(config.CtxConfig)

	if ctxConfig.IsNewUser() {
		http.Error(res, "Cookie not set.", http.StatusUnauthorized)

	}

	userID := ctxConfig.GetUserID()

	// get Short URLs for userID
	shorts := u.serviceURL.GetUserAll(req.Context(), userID)
	zap.S().Infof("Found: %d saved URL for User with ID: %s", len(shorts), userID)

	// if no data - 204
	if len(shorts) == 0 {
		res.WriteHeader(http.StatusNoContent)
		res.Write([]byte("Cookie not set"))
		return
	}

	resAuth := []ResponseAuth{}

	for _, short := range shorts {
		_, answerURL := u.serviceURL.GetAnsURL("http", u.conf.Response, short.Brief)
		resAuth = append(resAuth, ResponseAuth{Brief: answerURL.String(), Origin: short.Origin})
	}

	jsonURL, err := json.Marshal(resAuth)
	if err != nil {
		http.Error(res, "Error during Marshal User's URLs", http.StatusInternalServerError)
	}
	zap.S().Infoln("Server ansver with user's short URLs in JSON: ", string(jsonURL))

	// set content type
	res.Header().Add("Content-Type", "application/json")

	// set status code 200
	res.WriteHeader(http.StatusOK)
	res.Write(jsonURL)
}
