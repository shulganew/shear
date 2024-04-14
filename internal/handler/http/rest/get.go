package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
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
