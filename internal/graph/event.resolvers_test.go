package graph

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	t.Run("creates the event", func(t *testing.T) {
		r := newResolver(t)

		event, err := r.Mutation().CreateEvent(t.Context(), eventInput())

		require.NoError(t, err)
		assert.Equal(t, "Go meetup", event.Title)
		assert.Nil(t, event.Capacity)
	})
}

func TestEvent(t *testing.T) {
	t.Run("returns the event", func(t *testing.T) {
		r := newResolver(t)
		created, err := r.Mutation().CreateEvent(t.Context(), eventInput())
		require.NoError(t, err)

		found, err := r.Query().Event(t.Context(), created.ID)

		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
	})

	t.Run("returns nil for unknown ID", func(t *testing.T) {
		r := newResolver(t)

		found, err := r.Query().Event(t.Context(), 999)

		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestEvents(t *testing.T) {
	t.Run("sorts soonest first across time zones", func(t *testing.T) {
		r := newResolver(t)
		later := eventInput()
		later.StartsAt = time.Date(2030, 1, 1, 5, 0, 0, 0, time.UTC)
		later.EndsAt = later.StartsAt.Add(time.Hour)
		sooner := eventInput()
		sooner.StartsAt = time.Date(2030, 1, 1, 10, 0, 0, 0, time.FixedZone("JST", 9*60*60))
		sooner.EndsAt = sooner.StartsAt.Add(time.Hour)
		laterEvent, err := r.Mutation().CreateEvent(t.Context(), later)
		require.NoError(t, err)
		soonerEvent, err := r.Mutation().CreateEvent(t.Context(), sooner)
		require.NoError(t, err)

		events, err := r.Query().Events(t.Context())

		require.NoError(t, err)
		require.Len(t, events, 2)
		assert.Equal(t, []int{soonerEvent.ID, laterEvent.ID}, []int{events[0].ID, events[1].ID})
	})
}
