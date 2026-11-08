package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/allizjp/alliz-backend-interview/internal/ent/schema"
)

func TestEventFields(t *testing.T) {
	t.Run("rejects blank title", func(t *testing.T) {
		client := newClient(t)

		_, err := eventCreate(client).SetTitle("  ").Save(t.Context())

		assert.ErrorIs(t, err, schema.ErrEventTitleBlank)
	})

	t.Run("rejects capacity below one", func(t *testing.T) {
		client := newClient(t)

		_, err := eventCreate(client).SetCapacity(0).Save(t.Context())

		assert.ErrorIs(t, err, schema.ErrEventCapacityTooLow)
	})

	t.Run("allows no capacity", func(t *testing.T) {
		client := newClient(t)

		event, err := eventCreate(client).Save(t.Context())

		require.NoError(t, err)
		assert.Nil(t, event.Capacity)
	})
}

func TestCreateEvent(t *testing.T) {
	t.Run("rejects end at start", func(t *testing.T) {
		client := newClient(t)

		_, err := eventCreate(client).SetEndsAt(startsAt).Save(t.Context())

		assert.ErrorIs(t, err, schema.ErrEventEndsBeforeStart)
	})
}
