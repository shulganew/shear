package service

import (
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

func (d *Deleter) AsyncDelete(userID string, shorts []string) {
	d.waitDel.Add(1)
	go WriteFinal(Generator(shorts), userID, d.finalCh, d.waitDel)
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
func Generator(input []string) chan string {
	inputCh := make(chan string)
	go func() {
		defer close(inputCh)
		for _, data := range input {
			inputCh <- data
		}
	}()
	return inputCh
}
