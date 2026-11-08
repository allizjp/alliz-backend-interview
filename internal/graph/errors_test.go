package graph

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/allizjp/alliz-backend-interview/internal/ent/schema"
)

func TestPresentError(t *testing.T) {
	t.Run("shows public messages", func(t *testing.T) {
		err := fmt.Errorf("saving request: %w", schema.ErrEventFull)

		assert.Equal(t, "event is full", presentError(t.Context(), err).Message)
	})

	t.Run("hides other errors", func(t *testing.T) {
		err := errors.New("database is locked")

		assert.Equal(t, "internal error", presentError(t.Context(), err).Message)
	})
}
