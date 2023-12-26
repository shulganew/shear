package concurrent

type DelBrief struct {
	Briefs []string
	UserID string
}

// type ChGen struct {
// 	chans []chan DelBrief
// }

// func NewChgen(cond *sync.Cond) *ChGen {
// 	return &ChGen{chans: []chan DelBrief{}}

// }

// func (c *ChGen) AddChennel(ch chan DelBrief) chan DelBrief {
// 	var m sync.Mutex
// 	m.Lock()
// 	c.chans = append(c.chans, ch)
// 	m.Unlock()
// 	zap.S().Infoln("Channel added: Created:Nuber of channels: ", len(c.chans))
// 	return ch

// }

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
// func (c *ChGen) DeleteChannel(ch chan DelBrief) {
// 	var m sync.Mutex
// 	m.Lock()
// 	for i, localCh := range c.chans {
// 		if localCh == ch {
// 			c.chans = append(c.chans[:i], c.chans[i+1:]...)
// 			return
// 		}
// 	}
// 	m.Unlock()

// }

// func (c *ChGen) GetAllChannels() []chan DelBrief {
// 	return c.chans
// }

// fanIn объединяет несколько каналов resultChs в один.
