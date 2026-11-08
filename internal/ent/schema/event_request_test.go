package schema_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/allizjp/alliz-backend-interview/internal/ent/schema"
)

func TestEventRequestFields(t *testing.T) {
	t.Run("normalizes email", func(t *testing.T) {
		client := newClient(t)
		event := eventCreate(client).SaveX(t.Context())

		request, err := requestCreate(client, event).SetEmail("  Ana@Example.COM ").Save(t.Context())

		require.NoError(t, err)
		assert.Equal(t, "ana@example.com", request.Email)
	})

	t.Run("rejects blank name", func(t *testing.T) {
		client := newClient(t)
		event := eventCreate(client).SaveX(t.Context())

		_, err := requestCreate(client, event).SetName(" ").Save(t.Context())

		assert.ErrorIs(t, err, schema.ErrEventRequestNameBlank)
	})

	t.Run("rejects invalid email", func(t *testing.T) {
		client := newClient(t)
		event := eventCreate(client).SaveX(t.Context())

		_, err := requestCreate(client, event).SetEmail("ana.example.com").Save(t.Context())

		assert.ErrorIs(t, err, schema.ErrEventRequestEmailInvalid)
	})
}

func TestCreateEventRequest(t *testing.T) {
	t.Run("rejects unknown event", func(t *testing.T) {
		client := newClient(t)

		_, err := client.EventRequest.Create().
			SetEventID(999).
			SetName("Ana").
			SetEmail("ana@example.com").
			Save(t.Context())

		assert.ErrorIs(t, err, schema.ErrEventNotFound)
	})

	t.Run("allows started event", func(t *testing.T) {
		client := newClient(t)
		event := eventCreate(client).
			SetStartsAt(time.Now().Add(-time.Minute)).
			SetEndsAt(time.Now().Add(time.Hour)).
			SaveX(t.Context())

		_, err := requestCreate(client, event).Save(t.Context())

		assert.NoError(t, err)
	})

	t.Run("rejects full event", func(t *testing.T) {
		client := newClient(t)
		event := eventCreate(client).SetCapacity(1).SaveX(t.Context())
		requestCreate(client, event).SaveX(t.Context())

		_, err := requestCreate(client, event).SetEmail("bo@example.com").Save(t.Context())

		assert.ErrorIs(t, err, schema.ErrEventFull)
	})

	t.Run("is unlimited without capacity", func(t *testing.T) {
		client := newClient(t)
		event := eventCreate(client).SaveX(t.Context())
		requestCreate(client, event).SetEmail("ana@example.com").SaveX(t.Context())
		requestCreate(client, event).SetEmail("bo@example.com").SaveX(t.Context())

		_, err := requestCreate(client, event).SetEmail("cy@example.com").Save(t.Context())

		assert.NoError(t, err)
	})

	t.Run("rejects duplicate email", func(t *testing.T) {
		client := newClient(t)
		event := eventCreate(client).SaveX(t.Context())
		requestCreate(client, event).SetEmail("ana@example.com").SaveX(t.Context())

		_, err := requestCreate(client, event).SetEmail("ANA@example.com").Save(t.Context())

		assert.ErrorIs(t, err, schema.ErrEventRequestDuplicate)
	})

	t.Run("allows same email on another event", func(t *testing.T) {
		client := newClient(t)
		first := eventCreate(client).SaveX(t.Context())
		second := eventCreate(client).SaveX(t.Context())
		requestCreate(client, first).SaveX(t.Context())

		_, err := requestCreate(client, second).Save(t.Context())

		assert.NoError(t, err)
	})
}
