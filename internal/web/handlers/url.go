package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/shulganew/shear.git/internal/storage"
	"go.uber.org/zap"
)

// API handler for  GET and POST  short and long urls:
//
// Post "/"
//
// Get  "/{id}"
type HandlerURL struct {
	serviceURL *service.Shorten
	conf       *config.Config
}

// Service constructor.
func NewHandlerGetURL(conf *config.Config, short *service.Shorten) *HandlerURL {
	return &HandlerURL{serviceURL: short, conf: conf}
}

// GET and redirect by brief.
// @Summary      Get origin URL by brief (short) URL
// @Description  get short by id
// @Tags         api
// @Param        id   path  string  true  "brief URL"
// @Success      307
// @Failure      410
// @Failure      404
// @Router       /{id} [get]
func (u *HandlerURL) GetURL(res http.ResponseWriter, req *http.Request) {
	brief := chi.URLParam(req, "id")

	// get long Url from storage
	zap.S().Infoln("ID: ", brief)
	origin, exist, isDeleted := u.serviceURL.GetOrigin(req.Context(), brief)

	// set content type
	res.Header().Add("Content-Type", "text/plain")

	if exist {
		if isDeleted {
			// set status code 410
			res.WriteHeader(http.StatusGone)
			return
		}
		res.Header().Set("Location", origin)
		// set status code 307
		res.WriteHeader(http.StatusTemporaryRedirect)

		return
	}

	res.WriteHeader(http.StatusNotFound)
}

// POTS and set generate short URL.
// @Summary      Set origin URL
// @Description  set URL in body POST
// @Tags         api
// @Accept       plain
// @Produce      plain
// @Success      201 {string}  string  "Created"
// @Failure      401 "User unauthorized"
// @Failure      404 "Conflict. URL existed."
// @Failure      500 "Handling error"
// @Router       / [post]
func (u *HandlerURL) SetURL(res http.ResponseWriter, req *http.Request) {
	readBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Body not found", http.StatusInternalServerError)
	}

	redirectURL, err := url.Parse(string(readBody))
	if err != nil {
		http.Error(res, "Wrong URL in body, parse error", http.StatusInternalServerError)
	}
	zap.S().Infoln("redirectURL ", redirectURL)
	brief := service.GenerateShortLinkByte()
	mainURL, answerURL, err := u.serviceURL.GetAnsURLFast(redirectURL.Scheme, u.conf.Response, brief)
	if err != nil {
		http.Error(res, "Error parse URL", http.StatusInternalServerError)
		return
	}

	// set content type
	res.Header().Add("Content-Type", "text/plain")

	// find UserID in cookies
	userID, err := req.Cookie("user_id")
	if err != nil {
		http.Error(res, "Can't find user in cookies", http.StatusUnauthorized)
	}

	// save map to storage
	err = u.serviceURL.SetURL(req.Context(), userID.Value, brief, (*redirectURL).String())

	if err != nil {
		var tagErr *storage.ErrDuplicatedURL
		if errors.As(err, &tagErr) {
			// set status code 409 Conflict
			res.WriteHeader(http.StatusConflict)

			//send existed string from error
			var answer string
			answer, err = url.JoinPath(mainURL, tagErr.Brief)
			if err != nil {
				zap.S().Errorln("Error during JoinPath", err)
			}
			res.Write([]byte(answer))
			return
		}

		zap.S().Errorln(err)
		http.Error(res, "Error saving in Storage.", http.StatusInternalServerError)
	}
	// set status code 201
	res.WriteHeader(http.StatusCreated)
	// send generate and saved string
	res.Write([]byte(answerURL.String()))
}
