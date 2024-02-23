package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {

	t.Run("Check duplicated error", func(t *testing.T) {
		origin := "yandex.ru"
		brief := "qwertu"
		duplicated := errors.New("By duplicated error")
		error := NewErrDuplicatedURL(brief, origin, duplicated)
		require.True(t, errors.Is(error, duplicated))
		require.NotEmpty(t, error.Error())
	})

	t.Run("Check duplicated short error", func(t *testing.T) {
		origin := "yandex.ru"
		brief := "qwertu"
		session := "125453"
		duplicated := errors.New("By duplicated error")
		error := NewErrDuplicatedShort(session, brief, origin, duplicated)
		require.True(t, errors.Is(error, duplicated))
		require.NotEmpty(t, error.Error())
	})
}
