package schema

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {
	errBlank := errors.New("blank")
	validate := notBlank(errBlank)

	t.Run("notBlank rejects blank strings", func(t *testing.T) {
		for _, blank := range []string{"", " ", "\t\n"} {
			assert.ErrorIs(t, validate(blank), errBlank, "%q", blank)
		}
	})

	t.Run("notBlank accepts text", func(t *testing.T) {
		for _, filled := range []string{"a", " Go meetup "} {
			assert.NoError(t, validate(filled), "%q", filled)
		}
	})

	t.Run("validateEmail accepts plain addresses", func(t *testing.T) {
		for _, valid := range []string{"ana@example.com", "ana+events@example.co.jp"} {
			assert.NoError(t, validateEmail(valid), "%q", valid)
		}
	})

	t.Run("validateEmail rejects anything else", func(t *testing.T) {
		for _, invalid := range []string{"", "ana", "ana@", "@example.com", "Ana <ana@example.com>", " ana@example.com"} {
			assert.ErrorIs(t, validateEmail(invalid), ErrEventRequestEmailInvalid, "%q", invalid)
		}
	})
}
