# alliz backend interview

Events and event requests over GraphQL, on a simplified alliz stack: Go, gqlgen, ent + entgql, SQLite.

## Stack

| Tool                                                                                       | Role in this repo                                                         |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------- |
| [Go](https://go.dev)                                                                       | Language                                                                  |
| [gqlgen](https://gqlgen.com)                                                               | GraphQL server, generated from `internal/graph/schemas/`                  |
| [ent](https://entgo.io)                                                                    | ORM, generated from `internal/ent/schema/`; hooks hold the business rules |
| [entgql](https://entgo.io/docs/graphql)                                                    | Generates the GraphQL types and create inputs from the ent schemas        |
| [SQLite](https://sqlite.org) ([modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)) | Database; pure-Go driver                                                  |
| [testify](https://github.com/stretchr/testify)                                             | Test assertions                                                           |

## Run

Requires Go 1.26.

```sh
go run .
```

Open http://localhost:8080 for the GraphQL playground. Data is kept in `alliz.db`.

```graphql
mutation {
  createEvent(
    input: {
      title: "Go meetup"
      startsAt: "2030-01-01T18:00:00Z"
      endsAt: "2030-01-01T20:00:00Z"
      capacity: 20
    }
  ) {
    id
    title
  }
}
```

```graphql
query {
  events {
    id
    title
    startsAt
    requests {
      name
      email
    }
  }
}
```

## Commands

| Command             | Description                                                    |
| ------------------- | -------------------------------------------------------------- |
| `go run .`          | Starts the server                                              |
| `go test ./...`     | Runs all tests                                                 |
| `go generate ./...` | Regenerates ent, then gqlgen code; generated code is committed |
