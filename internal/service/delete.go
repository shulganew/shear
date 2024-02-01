package service

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/shulganew/shear.git/internal/config"
)

type Deleter struct {
	storeURLs StorageURL
	finalCh   chan DelBatch
	waitDel   *sync.WaitGroup
	conf      *config.Config
}

// return service
func NewDelete(storage *StorageURL, finalCh chan DelBatch, waitDel *sync.WaitGroup, conf *config.Config) *Deleter {
	return &Deleter{storeURLs: *storage, finalCh: finalCh, waitDel: waitDel, conf: conf}
}

// stuct for working with concurrent requests for delete update with channes - fanIn pattern
type DelBatch struct {
	UserID string
	Briefs []string
}

func (d *Deleter) AsyncDelete(userID string, dec *json.Decoder) {
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
		if len(breifs) == d.conf.BatchSize {

			tmp := make([]string, len(breifs))
			copy(tmp, breifs)
			d.waitDel.Add(1)
			go WriteFinal(Generator(tmp, d.conf.BatchSize), userID, d.finalCh, d.waitDel)
			breifs = breifs[:0]
		}
	}

	if len(breifs) != 0 {
		d.waitDel.Add(1)
		go WriteFinal(Generator(breifs, d.conf.BatchSize), userID, d.finalCh, d.waitDel)
	}

}

// funcktions for anync mark delete users URL

// Write data from handlers to final channel
func WriteFinal(input chan string, userID string, finalCh chan DelBatch, waitDel *sync.WaitGroup) {

	buff := make([]string, 0)
	//read to buffer from generator channel
	for data := range input {
		buff = append(buff, data)
	}
	finalCh <- DelBatch{UserID: userID, Briefs: buff}
	waitDel.Done()
}

// return channel with useres briefs for sending in final channel (fan-in)
func Generator(input []string, batchsize int) chan string {

	inputCh := make(chan string, batchsize)
	go func() {
		defer close(inputCh)
		for _, data := range input {
			inputCh <- data
		}
	}()
	return inputCh
}
