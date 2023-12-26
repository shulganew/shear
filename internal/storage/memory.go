package storage

import (
	"context"
	"slices"

	"github.com/shulganew/shear.git/internal/service"
)

type Memory struct {
	StoreURLs []service.Short
}

func NewMemory() *Memory {
	return &Memory{StoreURLs: make([]service.Short, 0)}
}

func (m *Memory) Set(ctx context.Context, userID, brief, origin string) (err error) {
	//init storage
	short := service.NewShort(len(m.StoreURLs), userID, brief, origin, "")
	m.StoreURLs = append(m.StoreURLs, *short)
	return
}

func (m *Memory) SetAll(ctx context.Context, shotrs []service.Short) error {
	m.StoreURLs = append(m.StoreURLs, shotrs...)
	return nil
}

func (m Memory) GetOrigin(ctx context.Context, brief string) (origin string, existed bool, isDeleted bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s service.Short) bool { return s.Brief == brief })
	if id != -1 {
		origin = m.StoreURLs[id].Origin
		existed = true
		isDeleted = m.StoreURLs[id].IsDeleted
	}

	return
}

func (m Memory) GetBrief(ctx context.Context, origin string) (brief string, existed bool, isDeleted bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s service.Short) bool { return s.Origin == origin })
	if id != -1 {
		brief = m.StoreURLs[id].Brief
		existed = true
		isDeleted = m.StoreURLs[id].IsDeleted
	}
	return

}

func (m *Memory) GetAll(ctx context.Context) []service.Short {
	return m.StoreURLs
}

func (m Memory) GetUserAll(ctx context.Context, userID string) []service.Short {
	slices.DeleteFunc(m.StoreURLs, func(s service.Short) bool { return s.UUID.String == userID })
	return m.StoreURLs
}

func (m *Memory) DelelteBatch(ctx context.Context, userID string, briefs []string) {
	for _, brief := range briefs {
		id := slices.IndexFunc(m.StoreURLs, func(s service.Short) bool { return s.Brief == brief && s.UUID.String == userID })
		if id != -1 {
			m.StoreURLs[id].IsDeleted = true

		}
	}
}
