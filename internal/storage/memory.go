package storage

import (
	"context"
	"slices"

	"github.com/shulganew/shear.git/internal/entities"
)

// In-memory storage.
type Memory struct {
	StoreURLs []entities.Short
}

// Storage constructor.
func NewMemory() *Memory {
	return &Memory{StoreURLs: make([]entities.Short, 0)}
}

// Set short and original URL to storage.
func (m *Memory) Set(ctx context.Context, userID, brief, origin string) (err error) {
	// init storage
	short := entities.NewShort(len(m.StoreURLs), userID, brief, origin, "", "")
	m.StoreURLs = append(m.StoreURLs, *short)
	return
}

// Set all user's short and original URLs from Short slice.
func (m *Memory) SetAll(ctx context.Context, shotrs []entities.Short) error {
	m.StoreURLs = append(m.StoreURLs, shotrs...)
	return nil
}

// Get original URL from storage.
func (m *Memory) GetOrigin(ctx context.Context, brief string) (origin string, existed bool, isDeleted bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s entities.Short) bool { return s.Brief == brief })
	if id != -1 {
		origin = m.StoreURLs[id].Origin
		existed = true
		isDeleted = m.StoreURLs[id].IsDeleted
	}
	return
}

// Get short URL from storage.
func (m *Memory) GetBrief(ctx context.Context, origin string) (brief string, existed bool, isDeleted bool) {
	id := slices.IndexFunc(m.StoreURLs, func(s entities.Short) bool { return s.Origin == origin })
	if id != -1 {
		brief = m.StoreURLs[id].Brief
		existed = true
		isDeleted = m.StoreURLs[id].IsDeleted
	}
	return

}

// Get all short and original URLs in Short slice.
func (m *Memory) GetAll(ctx context.Context) []entities.Short {
	return m.StoreURLs
}

// Get all user's short and original URLs in Short slice.
func (m *Memory) GetUserAll(ctx context.Context, userID string) []entities.Short {
	shorts := make([]entities.Short, 0)
	for _, short := range m.StoreURLs {
		id, err := short.UserID.Value()
		if err == nil {
			if id == userID {
				shorts = append(shorts, short)
			}
		}
	}
	return shorts
}

// Mark all user's URLs by short URL in briefs slice.
func (m *Memory) DeleteBatch(ctx context.Context, userID string, briefs []string) error {
	for _, brief := range briefs {
		id := slices.IndexFunc(m.StoreURLs, func(s entities.Short) bool { return s.Brief == brief && s.UserID.String == userID })
		if id != -1 {
			m.StoreURLs[id].IsDeleted = true
		}
	}
	return nil
}

// Get totoal number of shorts.
func (m *Memory) GetNumShorts(ctx context.Context) (num int, err error) {
	return len(m.StoreURLs), nil
}

// Get totoal number of Users.
func (m *Memory) GetNumUsers(ctx context.Context) (num int, err error) {
	// Number of users.
	users := make(map[string]struct{})
	for _, short := range m.StoreURLs {
		user := short.UserID.String
		if len(user) != 0 {
			users[user] = struct{}{}
		}
	}
	return len(users), nil
}
