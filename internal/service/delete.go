package service

import (
	"sync"

	"github.com/shulganew/shear.git/internal/config"
)

// Main struct for aggregation and delete user's URL in goroutine.
type Delete struct {
	delCh chan DelBatch
	conf  *config.Config
	wg    *sync.WaitGroup
}

// Constructor, create Delete service.
func NewDelete(delCh chan DelBatch, conf *config.Config) *Delete {
	w := &sync.WaitGroup{}
	return &Delete{delCh: delCh, conf: conf, wg: w}
}

// Struct for working with concurrent requests for delete update with channels - fanIn pattern.
type DelBatch struct {
	UserID string
	Briefs []string
}

// Async delete user's URL in goroutine.
func (d *Delete) AsyncDelete(userID string, shorts []string) {
	d.wg.Add(1)
	WriteFinal(Generator(shorts), userID, d.delCh, d.wg)
}

// Method for waiting delete service.
func (d *Delete) Stop() (delServDone chan struct{}) {
	delServDone = make(chan struct{})
	go func() {
		d.wg.Wait()
		close(delServDone)
	}()
	return delServDone
}

// Write data from handlers to final channel.
func WriteFinal(input chan string, userID string, delCh chan DelBatch, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		buff := make([]string, 0)
		//read to buffer from generator channel
		for data := range input {
			buff = append(buff, data)
		}
		delCh <- DelBatch{UserID: userID, Briefs: buff}
	}()
}

// Return channel with user's briefs for sending in final channel (fan-in).
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
