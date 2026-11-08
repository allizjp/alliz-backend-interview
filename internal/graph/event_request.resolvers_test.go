package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEventRequest(t *testing.T) {
	t.Run("creates the request", func(t *testing.T) {
		r := newResolver(t)
		event, err := r.Mutation().CreateEvent(t.Context(), eventInput())
		require.NoError(t, err)

		request, err := r.Mutation().CreateEventRequest(t.Context(), requestInput(event.ID))

		require.NoError(t, err)
		assert.Equal(t, event.ID, request.EventID)
		assert.Equal(t, "ana@example.com", request.Email)
	})
}
