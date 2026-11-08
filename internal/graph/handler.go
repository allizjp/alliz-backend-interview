package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/allizjp/alliz-backend-interview/internal/ent"
	"github.com/allizjp/alliz-backend-interview/internal/graph/generated"
)

type Resolver struct {
	client *ent.Client
}

func NewHandler(client *ent.Client) *handler.Server {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{client: client},
	}))
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})
	srv.SetErrorPresenter(presentError)
	srv.SetRecoverFunc(recoverPanic)
	return srv
}
