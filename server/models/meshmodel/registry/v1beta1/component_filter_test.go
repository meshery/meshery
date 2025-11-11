package v1beta1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentFilter(t *testing.T) {
	t.Run("Test Exclusion Fields", func(t *testing.T) {
		f := ComponentFilter{
			Name:         "test",
			Exclude:      "List",
			ExcludeRegex: ".*List$",
		}

		assert.Equal(t, "List", f.Exclude)
		assert.Equal(t, ".*List$", f.ExcludeRegex)
	})

	t.Run("Test GetById with nil db", func(t *testing.T) {
		f := &ComponentFilter{}
		_, err := f.GetById(nil)
		assert.Error(t, err)
	})

	t.Run("Test Get with nil db", func(t *testing.T) {
		f := &ComponentFilter{}
		_, _, _, err := f.Get(nil)
		assert.Error(t, err)
	})
}
