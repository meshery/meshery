package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// invalidInstanceIDWarning is the part of configureInstanceID's warning that
// operators see in the server log.
const invalidInstanceIDWarning = "is not a valid UUID"

// newInstanceIDTestViper returns a viper set up as main sets it up, and a
// logger whose output is captured in the returned buffer.
func newInstanceIDTestViper(t *testing.T) (*viper.Viper, logger.Handler, *bytes.Buffer) {
	t.Helper()
	var logs bytes.Buffer
	log, err := logger.New("test", logger.Options{
		Format:      logger.SyslogLogFormat,
		LogLevel:    int(logrus.InfoLevel), // 0 is panic-only, which drops warnings
		Output:      &logs,
		ErrorOutput: &logs,
	})
	if err != nil {
		t.Fatalf("logger.New() failed with error: %v", err)
	}
	v := viper.New()
	v.AutomaticEnv() // as main does
	return v, log, &logs
}

// setInstanceIDEnv sets INSTANCE_ID for the test, or unsets it when set is
// false, whatever the outer environment holds. t.Setenv restores it afterwards.
func setInstanceIDEnv(t *testing.T, value string, set bool) {
	t.Helper()
	t.Setenv("INSTANCE_ID", value)
	if !set {
		if err := os.Unsetenv("INSTANCE_ID"); err != nil {
			t.Fatalf("os.Unsetenv(INSTANCE_ID) failed with error: %v", err)
		}
	}
}

// storedInstanceID makes the type assertion handlers.NewHandlerInstance makes
// on INSTANCE_ID at boot. It panics there when the assertion fails.
func storedInstanceID(t *testing.T, v *viper.Viper) *core.Uuid {
	t.Helper()
	stored, ok := v.Get("INSTANCE_ID").(*core.Uuid)
	if !ok || stored == nil {
		t.Fatalf("viper.Get(\"INSTANCE_ID\") = %#v (%T), want a non-nil *core.Uuid: the server panics at boot on this assertion", v.Get("INSTANCE_ID"), v.Get("INSTANCE_ID"))
	}
	return stored
}

// checkWarning fails the test unless the captured log holds the invalid-value
// warning exactly when want is true.
func checkWarning(t *testing.T, logs *bytes.Buffer, want bool) {
	t.Helper()
	if got := strings.Contains(logs.String(), invalidInstanceIDWarning); got != want {
		t.Errorf("logged an invalid INSTANCE_ID warning = %v, want %v; log:\n%s", got, want, logs.String())
	}
}

// TestConfigureInstanceIDFromEnvironment sets INSTANCE_ID in the environment,
// as the environment variable docs show (INSTANCE_ID=$(uuidgen) make server).
// The server must use that ID, and every reader must get it as the *core.Uuid
// that handlers type-assert, not as the raw string.
func TestConfigureInstanceIDFromEnvironment(t *testing.T) {
	want := uuid.Must(uuid.NewV4())
	setInstanceIDEnv(t, want.String(), true)
	v, log, logs := newInstanceIDTestViper(t)

	got, err := configureInstanceID(v, log)
	if err != nil {
		t.Fatalf("configureInstanceID() failed with error: %v", err)
	}
	stored := storedInstanceID(t, v)
	if got == nil || *got != want {
		t.Fatalf("configureInstanceID() = %v, want the INSTANCE_ID from the environment, %v", got, want)
	}
	if stored != got {
		t.Errorf("viper.Get(\"INSTANCE_ID\") = %v, want the ID configureInstanceID returned, %v", stored, got)
	}
	if s := v.GetString("INSTANCE_ID"); s != want.String() {
		t.Errorf("viper.GetString(\"INSTANCE_ID\") = %q, want %q", s, want.String())
	}
	checkWarning(t, logs, false)
}

// TestConfigureInstanceIDGenerated covers an unset INSTANCE_ID and values that
// are not a usable UUID, which each get a newly generated ID, and a valid UUID
// with surrounding whitespace, which is trimmed and used. Only a value that
// was given and is not a usable UUID warns; an empty value is treated as
// unset. In every case the ID is stored for every reader.
func TestConfigureInstanceIDGenerated(t *testing.T) {
	padded := uuid.Must(uuid.NewV4())

	cases := []struct {
		name        string
		env         string
		set         bool
		want        *uuid.UUID // nil: any newly generated, non-nil ID
		wantWarning bool
	}{
		{name: "unset"},
		{name: "empty", env: "", set: true},
		{name: "whitespace only", env: "   ", set: true, wantWarning: true},
		{name: "not a UUID", env: "my-meshery", set: true, wantWarning: true},
		{name: "nil UUID", env: uuid.Nil.String(), set: true, wantWarning: true},
		{name: "valid UUID with surrounding spaces", env: "  " + padded.String() + "\n", set: true, want: &padded},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setInstanceIDEnv(t, tc.env, tc.set)
			v, log, logs := newInstanceIDTestViper(t)

			got, err := configureInstanceID(v, log)
			if err != nil {
				t.Fatalf("configureInstanceID() failed with error: %v", err)
			}
			if got == nil || *got == uuid.Nil {
				t.Fatalf("configureInstanceID() = %v, want a non-nil ID", got)
			}
			if tc.want != nil && *got != *tc.want {
				t.Fatalf("configureInstanceID() = %v, want %v", got, *tc.want)
			}
			if stored := storedInstanceID(t, v); stored != got {
				t.Errorf("viper.Get(\"INSTANCE_ID\") = %v, want the ID configureInstanceID returned, %v", stored, got)
			}
			checkWarning(t, logs, tc.wantWarning)
		})
	}
}

// TestConfigureInstanceIDFromConfigFile reads INSTANCE_ID from
// server-config.env, as main does. The file's ID is used when the environment
// leaves INSTANCE_ID unset. An INSTANCE_ID in the environment wins over the
// file, even when it is only whitespace, which is why that case warns instead
// of silently replacing the file's ID.
func TestConfigureInstanceIDFromConfigFile(t *testing.T) {
	fileID := uuid.Must(uuid.NewV4())
	envID := uuid.Must(uuid.NewV4())
	config := filepath.Join(t.TempDir(), "server-config.env")
	if err := os.WriteFile(config, []byte("INSTANCE_ID="+fileID.String()+"\n"), 0o600); err != nil {
		t.Fatalf("writing %s: %v", config, err)
	}

	cases := []struct {
		name        string
		env         string
		set         bool
		want        *uuid.UUID // nil: any newly generated ID other than the file's
		wantWarning bool
	}{
		{name: "environment unset", want: &fileID},
		{name: "environment set", env: envID.String(), set: true, want: &envID},
		{name: "environment only whitespace", env: "   ", set: true, wantWarning: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setInstanceIDEnv(t, tc.env, tc.set)
			v, log, logs := newInstanceIDTestViper(t)
			v.SetConfigFile(config) // as main does with ./server-config.env
			if err := v.ReadInConfig(); err != nil {
				t.Fatalf("ReadInConfig() failed with error: %v", err)
			}

			got, err := configureInstanceID(v, log)
			if err != nil {
				t.Fatalf("configureInstanceID() failed with error: %v", err)
			}
			if tc.want != nil && (got == nil || *got != *tc.want) {
				t.Fatalf("configureInstanceID() = %v, want %v", got, *tc.want)
			}
			if tc.want == nil && (got == nil || *got == uuid.Nil || *got == fileID) {
				t.Fatalf("configureInstanceID() = %v, want a newly generated ID", got)
			}
			if stored := storedInstanceID(t, v); stored != got {
				t.Errorf("viper.Get(\"INSTANCE_ID\") = %v, want the ID configureInstanceID returned, %v", stored, got)
			}
			checkWarning(t, logs, tc.wantWarning)
		})
	}
}
