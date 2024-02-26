package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
)

// Handler for API:
//
//	Delete "/api/user/urls"
type DelShorts struct {
	servDelete *service.Delete
}

func NewHandlerDelShorts(serviceDel *service.Delete) *DelShorts {

	return &DelShorts{servDelete: serviceDel}
}

// Delete User's URLs from json array in request (mark as deleted with saving in DB)
func (d *DelShorts) DelUserURLs(res http.ResponseWriter, req *http.Request) {
	// get UserID from cxt values
	ctxConfig := req.Context().Value(config.CtxConfig{}).(config.CtxConfig)
	if ctxConfig.IsNewUser() {
		http.Error(res, "Cookie not set.", http.StatusUnauthorized)

	}

	userID := ctxConfig.GetUserID()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Can't read body. ", http.StatusInternalServerError)
	}
	// read body as buffer
	var breifs []string
	err = json.Unmarshal(body, &breifs)
	if err != nil {
		http.Error(res, "Can't parse JSON delete short's array. ", http.StatusInternalServerError)
	}

	// async delete Shorts from body
	d.servDelete.AsyncDelete(userID, breifs)

	// set content type
	res.Header().Add("Content-Type", "plain/text")

	// set status code 202
	res.WriteHeader(http.StatusAccepted)
	res.Write([]byte("Done."))
}
