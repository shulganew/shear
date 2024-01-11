package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
)

const BATCHSIZE int = 20

type DelShorts struct {
	serviceDel *service.Delete
}

func NewHandlerDelShorts(serviceDel *service.Delete) *DelShorts {

	return &DelShorts{serviceDel: serviceDel}
}

func (d *DelShorts) GetServiceURL() service.Delete {
	return *d.serviceDel
}

// Delete User's URLs from json array in request (mark as deleted with saving in DB)
func (d *DelShorts) DelUserURLs(res http.ResponseWriter, req *http.Request) {

	//get UserID from cxt values
	ctxConfig := req.Context().Value(config.CtxConfig{}).(config.CtxConfig)

	if ctxConfig.IsNewUser() {
		http.Error(res, "Cookie not set.", http.StatusUnauthorized)

	}

	userID := ctxConfig.GetUserID()

	//read body as buffer
	dec := json.NewDecoder(req.Body)

	//async delete Shorts from body
	d.serviceDel.AsyncDelete(userID, dec)

	// set content type
	res.Header().Add("Content-Type", "plain/text")

	//set status code 202
	res.WriteHeader(http.StatusAccepted)

	res.Write([]byte("Done."))

}
