package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/shulganew/shear.git/internal/concurrent"
	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
)

const BATCHSIZE int = 30

type DelShorts struct {
	serviceURL *service.Shortener
	conf       *config.Config
	chgen      *concurrent.ChGen
	cond       *sync.Cond
}

func NewHandlerDelShorts(conf *config.Config, stor *service.StorageURL, cg *concurrent.ChGen, cd *sync.Cond) *DelShorts {

	return &DelShorts{serviceURL: service.NewService(stor), conf: conf, chgen: cg, cond: cd}
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

		DeleteGorutine(req, u.chgen, userID, u.cond)
		// while the array contains values

		// set content type
		res.Header().Add("Content-Type", "plain/text")

		//set status code 202
		res.WriteHeader(http.StatusAccepted)

		res.Write([]byte("Done."))

	} else {
		http.Error(res, "Cookie not set or can't Open UserID Seal", http.StatusUnauthorized)
	}

}

func DeleteGorutine(req *http.Request, chgen *concurrent.ChGen, userID string, cond *sync.Cond) {

	//create slice of Breifs size of int BATCHSIZE
	inputCh := make(chan concurrent.DelBrief)
	chgen.AddChennel(inputCh)
	cond.Broadcast()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(req *http.Request, inputCh chan concurrent.DelBrief, userID string) {

		//read body as buffer
		dec := json.NewDecoder(req.Body)

		// read open bracket
		_, err := dec.Token()
		if err != nil {
			log.Fatal(err)
		}
		breifs := make([]string, 0, BATCHSIZE)
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
			} else {
				break
			}

			// send butch of short ULS (briefs) in channel, zise BATCHSIZE
			//if len(breifs) == BATCHSIZE {

			inputCh <- concurrent.DelBrief{UserID: userID, Briefs: breifs}
			//clean slice for next briefs from buffer
			//	breifs = breifs[:0]
			//	}

		}

		//send last breifs, that wasn't sent with batch BATCHSIZE (less BATCHSIZE)
		//if len(breifs) != 0 {
		//	inputCh <- concurrent.DelBrief{UserID: userID, Briefs: breifs}
		//}
		wg.Done()

	}(req, inputCh, userID)

	wg.Wait()
	close(inputCh)
	inputCh = nil
	

}
