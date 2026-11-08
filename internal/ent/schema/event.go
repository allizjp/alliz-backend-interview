package schema

import (
	"context"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/allizjp/alliz-backend-interview/internal/apperr"
	gen "github.com/allizjp/alliz-backend-interview/internal/ent"
	"github.com/allizjp/alliz-backend-interview/internal/ent/hook"
)

var (
	ErrEventTitleBlank      = apperr.New("title can't be blank")
	ErrEventEndsBeforeStart = apperr.New("endsAt must be after startsAt")
	ErrEventCapacityTooLow  = apperr.New("capacity must be at least 1")
	ErrEventNotFound        = apperr.New("event not found")
	ErrEventFull            = apperr.New("event is full")
)

type Event struct {
	ent.Schema
}

func (Event) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Mutations(entgql.MutationCreate()),
	}
}

func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").Validate(notBlank(ErrEventTitleBlank)),
		field.String("description").Optional(),
		field.Time("starts_at"),
		field.Time("ends_at"),
		field.Int("capacity").
			Optional().
			Nillable().
			Comment("Unlimited when null").
			Validate(func(capacity int) error {
				if capacity < 1 {
					return ErrEventCapacityTooLow
				}
				return nil
			}),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Annotations(entgql.Skip(entgql.SkipMutationCreateInput)),
	}
}

func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("requests", EventRequest.Type).
			Annotations(entgql.Skip(entgql.SkipMutationCreateInput)),
	}
}

func (Event) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(checkDates, ent.OpCreate),
	}
}

func checkDates(next ent.Mutator) ent.Mutator {
	return hook.EventFunc(func(ctx context.Context, m *gen.EventMutation) (ent.Value, error) {
		startsAt, _ := m.StartsAt()
		endsAt, _ := m.EndsAt()
		if !endsAt.After(startsAt) {
			return nil, ErrEventEndsBeforeStart
		}
		return next.Mutate(ctx, m)
	})
}
