package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
)

const BATCHSIZE int = 20

type DelShorts struct {
	serviceURL *service.Shortener
	conf       *config.Config
	finalCh    chan service.DelBatch
	waitDel    *sync.WaitGroup
}

func NewHandlerDelShorts(conf *config.Config, stor *service.StorageURL, finalCh chan service.DelBatch, waitDel *sync.WaitGroup) *DelShorts {

	return &DelShorts{serviceURL: service.NewService(stor), conf: conf, finalCh: finalCh, waitDel: waitDel}
}

func (d *DelShorts) GetServiceURL() service.Shortener {
	return *d.serviceURL
}

// Delete User's URLs from json array in request (mark as deleted with saving in DB)
func (d *DelShorts) DelUserURLs(res http.ResponseWriter, req *http.Request) {

	if userID, ok := service.GetCodedUserID(req, d.conf.Pass); ok {
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
		DeleteGorutine(req, userID, d.serviceURL, d.finalCh, d.waitDel)

		// set content type
		res.Header().Add("Content-Type", "plain/text")

		//set status code 202
		res.WriteHeader(http.StatusAccepted)

		res.Write([]byte("Done."))

	} else {
		http.Error(res, "Cookie not set or can't Open UserID Seal", http.StatusUnauthorized)
	}

}

func DeleteGorutine(req *http.Request, userID string, stor *service.Shortener, finalCh chan service.DelBatch, waitDel *sync.WaitGroup) {

	//read body as buffer
	dec := json.NewDecoder(req.Body)

	// read open bracket
	_, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	breifs := make([]string, 0)

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
			waitDel.Add(1)
			go writeOut(generator(tmp), userID, stor, finalCh, waitDel)
			breifs = breifs[:0]
		}
	}
	if len(breifs) != 0 {
		waitDel.Add(1)
		go writeOut(generator(breifs), userID, stor, finalCh, waitDel)
	}

}

func writeOut(input chan string, userID string, stor *service.Shortener, finalCh chan service.DelBatch, waitDel *sync.WaitGroup) {

	buff := make([]string, 0)
	//read to buffer from generator channel
	for data := range input {
		buff = append(buff, data)
	}
	//stor.DelelteBatch(context.Background(), userID, buff)
	finalCh <- service.DelBatch{UserID: userID, Briefs: buff}
	waitDel.Done()
}

// return channel with useres briefs for sending in final channel (fan-in)
func generator(input []string) chan string {

	inputCh := make(chan string, BATCHSIZE)
	go func() {
		defer close(inputCh)
		for _, data := range input {
			inputCh <- data
		}
	}()
	return inputCh
}
