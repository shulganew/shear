package storage

import (
	"context"
	"net/url"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/shulganew/shear.git/internal/service"
	"github.com/stretchr/testify/require"
)

func TestMem(t *testing.T) {

	mem := NewMemory()
	ctx := context.Background()

	var usersID []uuid.UUID
	for i := 0; i < 3; i++ {
		// create new user uuid
		userID, err := uuid.NewV7()
		require.NoError(t, err)
		usersID = append(usersID, userID)
		for j := 0; j < 10; j++ {
			URLstr, err := url.JoinPath("http://", "yandex"+strconv.Itoa(i*10+j), ".ru")
			require.NoError(t, err)
			brief := service.GenerateShortLink()
			mem.Set(ctx, userID.String(), brief, URLstr)
		}
	}

	t.Run("Memory", func(t *testing.T) {
		shorts := mem.GetAll(ctx)
		require.Equal(t, 30, len(shorts))

		for _, userID := range usersID {
			shorts = mem.GetUserAll(ctx, userID.String())
			require.Equal(t, len(shorts), 10)
		}

		for _, userID := range usersID {
			shorts := mem.GetUserAll(ctx, userID.String())
			briefs := []string{shorts[0].Brief, shorts[2].Brief, shorts[4].Brief}
			mem.DeleteBatch(ctx, userID.String(), briefs)
		}

		shorts = mem.GetAll(ctx)
		require.Equal(t, 30, len(shorts))
		total := 0
		for _, short := range shorts {
			if !short.IsDeleted {
				total++
			}
		}
		require.Equal(t, 21, total)

	})
}
