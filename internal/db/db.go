package db

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"

	"github.com/allizjp/alliz-backend-interview/internal/ent"
	_ "github.com/allizjp/alliz-backend-interview/internal/ent/runtime"
)

const InMemory = ":memory:"

// Times are stored as text and compared as text, so they are all written in
// UTC to keep sorting and range filters correct across time zones.
const sqliteOptions = "?_pragma=foreign_keys(1)&_timezone=UTC&_time_format=sqlite"

func SQLite(path string) (dialect.Driver, error) {
	sqlDB, err := sql.Open("sqlite", path+sqliteOptions)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	// Each connection to :memory: gets its own empty database.
	sqlDB.SetMaxOpenConns(1)
	return entsql.OpenDB(dialect.SQLite, sqlDB), nil
}

func Open(ctx context.Context, conn dialect.Driver) (*ent.Client, error) {
	client := ent.NewClient(ent.Driver(conn))
	if err := client.Schema.Create(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("creating tables: %w", err)
	}
	return client, nil
}
