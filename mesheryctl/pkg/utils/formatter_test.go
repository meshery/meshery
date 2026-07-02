package utils

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTerminalFormatter_Format(t *testing.T) {
	formatter := &TerminalFormatter{}
	entry := &log.Entry{Message: "test message"}

	b, err := formatter.Format(entry)

	assert.NoError(t, err)
	assert.Equal(t, "test message\n", string(b))
}

func TestSetupLogrusFormatter(t *testing.T) {
	assert.NotPanics(t, func() {
		SetupLogrusFormatter()
	})
}