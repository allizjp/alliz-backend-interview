package graph

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

func newResolver(t *testing.T) *Resolver {
	return &Resolver{client: newClient(t)}
}

func eventInput() ent.CreateEventInput {
	return ent.CreateEventInput{
		Title:    "Go meetup",
		StartsAt: startsAt,
		EndsAt:   startsAt.Add(2 * time.Hour),
	}
}

func requestInput(eventID int) ent.CreateEventRequestInput {
	return ent.CreateEventRequestInput{
		EventID: eventID,
		Name:    "Ana",
		Email:   "ana@example.com",
	}
}
