package schema

import (
	"context"
	"strings"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/allizjp/alliz-backend-interview/internal/apperr"
	gen "github.com/allizjp/alliz-backend-interview/internal/ent"
	"github.com/allizjp/alliz-backend-interview/internal/ent/hook"
)

var (
	ErrEventRequestNameBlank    = apperr.New("name can't be blank")
	ErrEventRequestEmailInvalid = apperr.New("email is invalid")
	ErrEventRequestDuplicate    = apperr.New("email already requested")
)

// eventRequestEmailIndexColumns is how SQLite names the unique email index
// in the error it returns when the index is violated.
const eventRequestEmailIndexColumns = "event_requests.email, event_requests.event_id"

type EventRequest struct {
	ent.Schema
}

func (EventRequest) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Mutations(entgql.MutationCreate()),
	}
}

func (EventRequest) Fields() []ent.Field {
	return []ent.Field{
		field.Int("event_id").Immutable(),
		field.String("name").Validate(notBlank(ErrEventRequestNameBlank)),
		field.String("email").Validate(validateEmail),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Annotations(entgql.Skip(entgql.SkipMutationCreateInput)),
	}
}

func (EventRequest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("requests").
			Field("event_id").
			Unique().
			Required().
			Immutable(),
	}
}

func (EventRequest) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email", "event_id").Unique(),
	}
}

func (EventRequest) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(normalizeEmail, ent.OpCreate),
		hook.On(checkEventTakesRequests, ent.OpCreate),
	}
}

func normalizeEmail(next ent.Mutator) ent.Mutator {
	return hook.EventRequestFunc(func(ctx context.Context, m *gen.EventRequestMutation) (ent.Value, error) {
		if email, ok := m.Email(); ok {
			m.SetEmail(strings.ToLower(strings.TrimSpace(email)))
		}
		return next.Mutate(ctx, m)
	})
}

func checkEventTakesRequests(next ent.Mutator) ent.Mutator {
	return hook.EventRequestFunc(func(ctx context.Context, m *gen.EventRequestMutation) (ent.Value, error) {
		eventID, _ := m.EventID()
		requested, err := m.Client().Event.Get(ctx, eventID)
		if gen.IsNotFound(err) {
			return nil, ErrEventNotFound
		}
		if err != nil {
			return nil, err
		}

		if requested.Capacity != nil {
			taken, err := requested.QueryRequests().Count(ctx)
			if err != nil {
				return nil, err
			}
			if taken >= *requested.Capacity {
				return nil, ErrEventFull
			}
		}

		created, err := next.Mutate(ctx, m)
		if gen.IsConstraintError(err) && strings.Contains(err.Error(), eventRequestEmailIndexColumns) {
			return nil, ErrEventRequestDuplicate
		}
		return created, err
	})
}
