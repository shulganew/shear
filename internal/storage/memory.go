package storage

import (
	"context"
	"slices"

	"github.com/shulganew/shear.git/internal/entities"
)

type Memory struct {
	StoreURLs []entities.Short
}

func NewMemory() *Memory {
	return &Memory{StoreURLs: make([]entities.Short, 0)}
}

func (m *Memory) Set(ctx context.Context, userID, brief, origin string) (err error) {
	//init storage
	short := entities.NewShort(len(m.StoreURLs), userID, brief, origin, "")
	m.StoreURLs = append(m.StoreURLs, *short)
	return
}

func (m *Memory) SetAll(ctx context.Context, shotrs []entities.Short) error {
	m.StoreURLs = append(m.StoreURLs, shotrs...)
	return nil
}

func (m *Memory) GetOrigin(ctx context.Context, brief string) (origin string, existed bool, isDeleted bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s entities.Short) bool { return s.Brief == brief })
	if id != -1 {
		origin = m.StoreURLs[id].Origin
		existed = true
		isDeleted = m.StoreURLs[id].IsDeleted
	}
	return
}

func (m *Memory) GetBrief(ctx context.Context, origin string) (brief string, existed bool, isDeleted bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s entities.Short) bool { return s.Origin == origin })
	if id != -1 {
		brief = m.StoreURLs[id].Brief
		existed = true
		isDeleted = m.StoreURLs[id].IsDeleted
	}
	return

}

func (m *Memory) GetAll(ctx context.Context) []entities.Short {
	return m.StoreURLs
}

func (m *Memory) GetUserAll(ctx context.Context, userID string) []entities.Short {
	shorts := make([]entities.Short, 0)
	for _, short := range m.StoreURLs {
		id, err := short.UUID.Value()
		if err == nil {
			if id == userID {
				shorts = append(shorts, short)
			}
		}
	}
	return shorts
}

func (m *Memory) DelelteBatch(ctx context.Context, userID string, briefs []string) {
	for _, brief := range briefs {
		id := slices.IndexFunc(m.StoreURLs, func(s entities.Short) bool { return s.Brief == brief && s.UUID.String == userID })
		if id != -1 {
			m.StoreURLs[id].IsDeleted = true
		}
	}
}
