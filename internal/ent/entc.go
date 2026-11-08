//go:build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("internal/graph/schemas/ent.graphqls"),
		entgql.WithConfigPath("gqlgen.yml"),
		entgql.WithRelaySpec(false),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}

	if err := entc.Generate(
		"./internal/ent/schema",
		&gen.Config{
			Target:  "internal/ent",
			Package: "github.com/allizjp/alliz-backend-interview/internal/ent",
		},
		entc.Extensions(ex),
	); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
