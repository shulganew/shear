package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
)

const BATCHSIZE int = 10

type DelShorts struct {
	serviceURL *service.Shortener
	conf       *config.Config
	cond       *sync.Cond
}

func NewHandlerDelShorts(conf *config.Config, stor *service.StorageURL) *DelShorts {

	return &DelShorts{serviceURL: service.NewService(stor), conf: conf}
}

func (u *DelShorts) GetServiceURL() service.Shortener {
	return *u.serviceURL
}

// Delete User's URLs from json array in request (mark as deleted with saving in DB)
func (u *DelShorts) DelUserURLs(res http.ResponseWriter, req *http.Request) {

	if userID, ok := service.GetCodedUserID(req, u.conf.Pass); ok {
		//cookie iser_id is set
		cookies := req.Cookies()

		//clean cookie data
		req.Header["Cookie"] = make([]string, 0)
		for _, cookie := range cookies {

			if cookie.Name == "user_id" {
				cookie.Value = userID
			}
			req.AddCookie(cookie)
		}

		//read the body and UPDATE DB in gorutine
		DeleteGorutine(req, userID, u.serviceURL)

		// set content type
		res.Header().Add("Content-Type", "plain/text")

		//set status code 202
		res.WriteHeader(http.StatusAccepted)

		res.Write([]byte("Done."))

	} else {
		http.Error(res, "Cookie not set or can't Open UserID Seal", http.StatusUnauthorized)
	}

}

func DeleteGorutine(req *http.Request, userID string, stor *service.Shortener) {

	//read body as buffer
	dec := json.NewDecoder(req.Body)

	// read open bracket
	_, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	breifs := make([]string, 0)

	doneCh := make(chan struct{})

	for dec.More() {

		var brief string
		// // decode an array value (Message)
		err := dec.Decode(&brief)
		if err != nil {
			log.Fatal(err)
		}
		//check end of json array
		if brief != "]" {

			breifs = append(breifs, brief)
		}

		// send butch of short ULS (briefs) in channel, zise BATCHSIZE
		if len(breifs) == BATCHSIZE {

			tmp := make([]string, len(breifs))
			copy(tmp, breifs)
			go writeDB(doneCh, generator(doneCh, tmp), userID, stor)
			breifs = breifs[:0]
		}

	}

	go writeDB(doneCh, generator(doneCh, breifs), userID, stor)

}

func writeDB(doneCh chan struct{}, input chan string, userID string, stor *service.Shortener) {

	buff := make([]string, 0)
	//read to buffer from generator channel
	for data := range input {
		buff = append(buff, data)
	}
	stor.DelelteBatch(context.Background(), userID, buff)
}

func generator(doneCh chan struct{}, input []string) chan string {
	inputCh := make(chan string, BATCHSIZE)

	go func() {
		defer close(inputCh)

		for _, data := range input {

			select {

			case <-doneCh:
				return

			case inputCh <- data:
			}
		}
	}()

	return inputCh
}
