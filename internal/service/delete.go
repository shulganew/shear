package service

import "sync"

// Write data from handlers to final channel
func WriteFinal(input chan string, userID string, stor *Shortener, finalCh chan DelBatch, waitDel *sync.WaitGroup) {

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
