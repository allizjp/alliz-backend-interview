package graph

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/allizjp/alliz-backend-interview/internal/apperr"
)

func presentError(ctx context.Context, err error) *gqlerror.Error {
	gqlErr := graphql.DefaultErrorPresenter(ctx, err)
	if public, ok := errors.AsType[*apperr.Public](err); ok {
		gqlErr.Message = public.Error()
		return gqlErr
	}

	slog.ErrorContext(ctx, "request failed", "path", gqlErr.Path.String(), "error", err)
	gqlErr.Message = "internal error"
	return gqlErr
}

func recoverPanic(_ context.Context, p any) error {
	return fmt.Errorf("panic: %v\n%s", p, debug.Stack())
}
