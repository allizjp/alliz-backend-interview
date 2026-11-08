package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/allizjp/alliz-backend-interview/internal/db"
	"github.com/allizjp/alliz-backend-interview/internal/graph"
)

const (
	addr   = "localhost:8080"
	dbPath = "alliz.db"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	conn, err := db.SQLite(dbPath)
	if err != nil {
		return err
	}
	client, err := db.Open(context.Background(), conn)
	if err != nil {
		return err
	}
	defer client.Close()

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("alliz", "/query"))
	mux.Handle("/query", graph.NewHandler(client))

	slog.Info("playground ready", "url", "http://"+addr)
	return http.ListenAndServe(addr, mux)
}
