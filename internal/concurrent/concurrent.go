package concurrent

import (
	"sync"

	"go.uber.org/zap"
)

type DelBrief struct {
	Briefs []string
	UserID string
}

type ChGen struct {
	chans []chan DelBrief
}

func NewChgen(cond *sync.Cond) *ChGen {
	return &ChGen{chans: []chan DelBrief{}}

}

func (c *ChGen) AddChennel(ch chan DelBrief) chan DelBrief {
	var m sync.Mutex
	m.Lock()
	c.chans = append(c.chans, ch)
	m.Unlock()
	zap.S().Infoln("Channel added: Created:Nuber of channels: ", len(c.chans))
	return ch

}

// func (c *ChGen) CreateChannel() chan DelBrief {
// 	var m sync.Mutex
// 	ch := make(chan DelBrief)
// 	m.Lock()
// 	c.chans = append(c.chans, ch)
// 	m.Unlock()
// 	zap.S().Infoln("Created:Nuber of channels: ", len(c.chans))
// 	return ch

// }

// delete channel form slice
func (c *ChGen) DeleteChannel(ch chan DelBrief) {
	var m sync.Mutex
	m.Lock()
	for i, localCh := range c.chans {
		if localCh == ch {
			c.chans = append(c.chans[:i], c.chans[i+1:]...)
			return
		}
	}
	m.Unlock()

}

func (c *ChGen) GetAllChannels() []chan DelBrief {
	return c.chans
}

// fanIn объединяет несколько каналов resultChs в один.
func FanIn(doneCh chan struct{}, resultChs []chan DelBrief) chan DelBrief {
	// конечный выходной канал в который отправляем данные из всех каналов из слайса, назовём его результирующим
	finalCh := make(chan DelBrief)

	// понадобится для ожидания всех горутин
	var wg sync.WaitGroup

	// перебираем все входящие каналы
	zap.S().Infoln("FANIN Channels number: ", len(resultChs))

	for _, ch := range resultChs {
		// в горутину передавать переменную цикла нельзя, поэтому делаем так
		chClosure := ch

		// инкрементируем счётчик горутин, которые нужно подождать
		wg.Add(1)

		go func() {
			// откладываем сообщение о том, что горутина завершилась
			defer wg.Done()

			// получаем данные из канала
			for data := range chClosure {
				//zap.S().Infoln("FANIN channel range")
				select {
				// // выходим из горутины, если канал закрылся
				case <-doneCh:
					return
					// // если не закрылся, отправляем данные в конечный выходной канал
				case finalCh <- data:
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
