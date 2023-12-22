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

func (m *Memory) Set(ctx context.Context, brief, origin string) (err error) {
	//init storage
	short := service.Short{ID: len(m.StoreURLs), Brief: brief, Origin: origin}
	m.StoreURLs = append(m.StoreURLs, short)
	return
}

func (m *Memory) GetOrigin(ctx context.Context, brief string) (origin string, ok bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s service.Short) bool { return s.Brief == brief })
	if id != -1 {
		origin = m.StoreURLs[id].Origin
		ok = true
	}
	return
}

func (m *Memory) GetBrief(ctx context.Context, origin string) (brief string, ok bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s service.Short) bool { return s.Origin == origin })
	if id != -1 {
		brief = m.StoreURLs[id].Brief
		ok = true
	}
	return

}

func (m *Memory) GetAll(ctx context.Context) []service.Short {
	return m.StoreURLs
}

func (m *Memory) SetAll(ctx context.Context, shotrs []service.Short) error {
	m.StoreURLs = append(m.StoreURLs, shotrs...)
	return nil
}
