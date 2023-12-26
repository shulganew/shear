package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/shulganew/shear.git/internal/config"
	"github.com/shulganew/shear.git/internal/service"
	"go.uber.org/zap"
)

const BATCHSIZE int = 2

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

		DeleteGorutine(req, userID, u.serviceURL)
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

func DeleteGorutine(req *http.Request, userID string, stor *service.Shortener) {

	//go func(req *http.Request, userID string, stor *service.Shortener) {

	//read body as buffer
	dec := json.NewDecoder(req.Body)

	// read open bracket
	_, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	breifs := make([]string, 0)

	//resultChs := make([]chan string, 0)

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
		//if len(breifs) == BATCHSIZE {
		//zap.S().Infoln("Breifs : ", breifs)
		//zap.S().Infoln("append breifs: ", len(breifs))
		//resultChs = append(resultChs, generator(doneCh, breifs))

		//tmp := make([]string, len(breifs))
		//copy(tmp, breifs)

		//wreteDB(doneCh, generator(doneCh, tmp), userID, stor)
		//breifs = breifs[:0]
		//}

	}
	zap.S().Infoln("Breifs append: ", breifs)

	writeDB(doneCh, generator(doneCh, breifs), userID, stor)

	//zap.S().Infoln("append breifs tail : ", len(breifs))
	//resultChs = append(resultChs, generator(doneCh, breifs))
	//zap.S().Infoln("Num of channels: ", len(resultChs))
	//finalCh := FanIn(doneCh, resultChs)
	//wreteDB(doneCh, finalCh, userID, stor)

	//}(req, userID, stor)

}

func writeDB(doneCh chan struct{}, input chan string, userID string, stor *service.Shortener) {
	// канал, в который будем отправлять данные из слайса
	buff := make([]string, 0)
	// горутина, в которой отправляем в канал  inputCh данные
	var wg sync.WaitGroup
	wg.Add(1)
	go func(doneCh chan struct{}, input chan string, userID string, st *service.Shortener, buff *[]string, wg *sync.WaitGroup) {

		zap.S().Infoln("Write")
		// перебираем все данные в слайсе
		for data := range input {
			zap.S().Infoln("Get form channel: ", data)
			*buff = append(*buff, data)

		}
		wg.Done()
	}(doneCh, input, userID, stor, &buff, &wg)

	wg.Wait()

	zap.S().Infoln("Buff: ", buff)
	stor.DelelteBatch(context.Background(), userID, buff)
	zap.S().Infoln("CLOSE input channel")

}

func generator(doneCh chan struct{}, input []string) chan string {
	// канал, в который будем отправлять данные из слайса
	inputCh := make(chan string)

	// горутина, в которой отправляем в канал  inputCh данные
	go func() {
		// как отправители закрываем канал, когда всё отправим
		defer close(inputCh)

		// перебираем все данные в слайсе
		for i, data := range input {
			zap.S().Infoln("Input ", data, " ", i)
			select {
			// если doneCh закрыт, сразу выходим из горутины
			case <-doneCh:
				zap.S().Infoln("Generator closed")
				return
			// если doneCh не закрыт, кидаем в канал inputCh данные data
			case inputCh <- data:
				zap.S().Infoln("Send data ", data)
			}
		}
	}()

	// возвращаем канал для данных
	return inputCh
}

func FanIn(doneCh chan struct{}, resultChs []chan string) chan string {
	// конечный выходной канал в который отправляем данные из всех каналов из слайса, назовём его результирующим
	finalCh := make(chan string)

	// понадобится для ожидания всех горутин
	var wg sync.WaitGroup

	for _, ch := range resultChs {
		// в горутину передавать переменную цикла нельзя, поэтому делаем так
		chClosure := ch

		// инкрементируем счётчик горутин, которые нужно подождать
		wg.Add(1)

		go func() {
			// откладываем сообщение о том, что горутина завершилась
			defer wg.Done()

			// получаем данные из канала
			for data1 := range chClosure {
				//zap.S().Infoln("FANIN channel range")
				select {
				// // выходим из горутины, если канал закрылся
				case <-doneCh:
					return
					// // если не закрылся, отправляем данные в конечный выходной канал
				case finalCh <- data1:
				}

			}
		}()
	}

	go func() {
		// ждём завершения всех горутин
		wg.Wait()
		// когда все горутины завершились, закрываем результирующий канал
		zap.S().Infoln("CLOSE final channel")
		close(finalCh)
	}()

	// возвращаем результирующий канал
	return finalCh
}
