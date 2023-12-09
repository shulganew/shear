package storage

import (
	"context"
	"slices"
)

// base stract for working with storage
type Short struct {
	ID int `json:"uuid"`
	//short URL (cache)
	Brief string `json:"short_url"`
	//Long full URL
	Origin string `json:"original_url"`
}

// intarface for universal data storage
type StorageURL interface {
	Set(ctx context.Context, brief, origin string) Short
	GetOrigin(ctx context.Context, brief string) (string, bool)
	GetBrief(ctx context.Context, origin string) (string, bool)
	GetAll(ctx context.Context) []Short
	SetAll(ctx context.Context, short []Short)
}

type Memory struct {
	StoreURLs []Short
}

func (m *Memory) Set(ctx context.Context, brief, origin string) (short Short) {
	//init storage
	short = Short{ID: len(m.StoreURLs), Brief: brief, Origin: origin}
	m.StoreURLs = append(m.StoreURLs, short)
	return
}

func (m *Memory) GetOrigin(ctx context.Context, brief string) (origin string, ok bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s Short) bool { return s.Brief == brief })
	if id != -1 {
		origin = m.StoreURLs[id].Origin
		ok = true
	}
	return
}

func (m *Memory) GetBrief(ctx context.Context, origin string) (brief string, ok bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s Short) bool { return s.Origin == origin })
	if id != -1 {
		brief = m.StoreURLs[id].Brief
		ok = true
	}
	return

}

func (m *Memory) GetAll(ctx context.Context) []Short {
	return m.StoreURLs
}

func (m *Memory) SetAll(ctx context.Context, shotrs []Short) {
	m.StoreURLs = append(m.StoreURLs, shotrs...)
}
