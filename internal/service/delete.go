package service

import (
	"sync"

	"github.com/shulganew/shear.git/internal/config"
)

// Main struct for aggregation and delete user's URL in goroutine.
type Delete struct {
	finalCh chan DelBatch
	waitDel *sync.WaitGroup
	conf    *config.Config
}

// Constructor, create Delete service.
func NewDelete(finalCh chan DelBatch, waitDel *sync.WaitGroup, conf *config.Config) *Delete {
	return &Delete{finalCh: finalCh, waitDel: waitDel, conf: conf}
}

// Struct for working with concurrent requests for delete update with channels - fanIn pattern.
type DelBatch struct {
	UserID string
	Briefs []string
}

// Async delete user's URL in goroutine.
func (d *Delete) AsyncDelete(userID string, shorts []string) {
	d.waitDel.Add(1)
	go WriteFinal(Generator(shorts), userID, d.finalCh, d.waitDel)
}

// Write data from handlers to final channel.
func WriteFinal(input chan string, userID string, finalCh chan DelBatch, waitDel *sync.WaitGroup) {
	buff := make([]string, 0)
	//read to buffer from generator channel
	for data := range input {
		buff = append(buff, data)
	}
	finalCh <- DelBatch{UserID: userID, Briefs: buff}
	waitDel.Done()
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
