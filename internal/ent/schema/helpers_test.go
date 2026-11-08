package schema_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/allizjp/alliz-backend-interview/internal/db"
	"github.com/allizjp/alliz-backend-interview/internal/ent"
)

var startsAt = time.Date(2030, 1, 1, 10, 0, 0, 0, time.UTC)

func newClient(t *testing.T) *ent.Client {
	t.Helper()
	conn, err := db.SQLite(db.InMemory)
	require.NoError(t, err)
	client, err := db.Open(t.Context(), conn)
	require.NoError(t, err)
	t.Cleanup(func() { client.Close() })
	return client
}

func eventCreate(client *ent.Client) *ent.EventCreate {
	return client.Event.Create().
		SetTitle("Go meetup").
		SetStartsAt(startsAt).
		SetEndsAt(startsAt.Add(2 * time.Hour))
}

func requestCreate(client *ent.Client, event *ent.Event) *ent.EventRequestCreate {
	return client.EventRequest.Create().
		SetEvent(event).
		SetName("Ana").
		SetEmail("ana@example.com")
}
