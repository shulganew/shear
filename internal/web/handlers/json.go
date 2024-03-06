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
func NewHandlerAPI(conf *config.Config, stor service.StorageURL) *HandlerAPI {
	return &HandlerAPI{serviceURL: service.NewService(stor), conf: conf}
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
	var request Request
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	origin, err := url.Parse(string(request.URL))
	if err != nil {
		http.Error(res, "Wrong URL in JSON, parse error", http.StatusInternalServerError)
	}
	brief := service.GenerateShortLinkByte()
	mainURL, answerURL := u.serviceURL.GetAnsURLFast(origin.Scheme, u.conf.Response, brief)

	// find UserID in cookies
	userID, err := req.Cookie("user_id")
	if err != nil {
		http.Error(res, "Can't find user in cookies", http.StatusUnauthorized)
	}

	// save map to storage
	err = u.serviceURL.SetURL(req.Context(), userID.Value, brief, (*origin).String())

	// set content type
	res.Header().Add("Content-Type", "application/json")

	if err != nil {

		var tagErr *storage.ErrDuplicatedURL
		if errors.As(err, &tagErr) {

			// get correct answer URL
			var answer string
			answer, err = url.JoinPath(mainURL, tagErr.Brief)
			if err != nil {
				zap.S().Errorln("Error during JoinPath", err)
			}

			// send existed string from error
			response := Response{answer}
			var jsonBrokenURL []byte
			jsonBrokenURL, err = json.Marshal(response)
			if err != nil {
				http.Error(res, "Error during Marshal answer URL", http.StatusInternalServerError)
			}
			zap.S().Infoln("Server ansver with short URL in JSON (duplicated request): ", string(jsonBrokenURL))

			// set status code 409
			res.WriteHeader(http.StatusConflict)
			res.Write([]byte(jsonBrokenURL))
			return
		}

		zap.S().Errorln(err)
		http.Error(res, "Error saving in Storage.", http.StatusInternalServerError)
	}

	response := Response{answerURL.String()}

	jsonURL, err := json.Marshal(response)
	if err != nil {
		http.Error(res, "Error during Marshal answer URL", http.StatusInternalServerError)
	}
	zap.S().Infoln("Server ansver with short URL in JSON: ", string(jsonURL))

	// set status code 201
	res.WriteHeader(http.StatusCreated)
	res.Write(jsonURL)
}
