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

// hadler for  GET and POST  short and long urls

type HandlerURL struct {
	serviceURL *service.Shortener
	conf       *config.Config
}

func NewHandlerGetURL(conf *config.Config, stor service.StorageURL) *HandlerURL {

	return &HandlerURL{serviceURL: service.NewService(stor), conf: conf}
}

// GET and redirect by brief
func (u *HandlerURL) GetURL(res http.ResponseWriter, req *http.Request) {
	brief := chi.URLParam(req, "id")

	//get long Url from storage
	zap.S().Infoln("ID: ", brief)
	origin, exist, isDeleted := u.serviceURL.GetOrigin(req.Context(), brief)

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	if exist {
		if isDeleted {
			//set status code 410
			res.WriteHeader(http.StatusGone)
			return
		}
		res.Header().Set("Location", origin)
		//set status code 307
		res.WriteHeader(http.StatusTemporaryRedirect)

		return
	}

	res.WriteHeader(http.StatusNotFound)
}

// POTS and set generate short Url
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
	brief := service.GenerateShorLink()
	mainURL, answerURL := u.serviceURL.GetAnsURL(redirectURL.Scheme, u.conf.Response, brief)

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	//find UserID in cookies
	userID, err := req.Cookie("user_id")
	if err != nil {
		http.Error(res, "Can't find user in cookies", http.StatusUnauthorized)
	}

	// Save map to storage.
	err = u.serviceURL.SetURL(req.Context(), userID.Value, brief, (*redirectURL).String())

	if err != nil {
		var tagErr *storage.ErrDuplicatedURL
		if errors.As(err, &tagErr) {
			//set status code 409 Conflict
			res.WriteHeader(http.StatusConflict)

			//send existed string from error
			answer, err := url.JoinPath(mainURL, tagErr.Brief)
			if err != nil {
				zap.S().Errorln("Error during JoinPath", err)
			}
			res.Write([]byte(answer))
			return
		}

		zap.S().Errorln(err)
		http.Error(res, "Error saving in Storage.", http.StatusInternalServerError)
	}
	//set status code 201
	res.WriteHeader(http.StatusCreated)
	//send generate and saved string
	res.Write([]byte(answerURL.String()))
}
