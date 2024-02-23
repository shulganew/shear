package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {

	t.Run("Check parse host", func(t *testing.T) {
		server := "localhost:8080"
		host, port := CheckURL(server)
		assert.Equal(t, "localhost", host)
		assert.Equal(t, "8080", port)
	})
}
