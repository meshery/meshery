package main

import (
	"os"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/viper"
)

func newInstanceIDTestViper(t *testing.T) (*viper.Viper, logger.Handler) {
	t.Helper()
	log, err := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})
	if err != nil {
		t.Fatalf("logger.New() failed with error: %v", err)
	}
	v := viper.New()
	v.AutomaticEnv() // as main does
	return v, log
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

// TestConfigureInstanceIDFromEnvironment sets INSTANCE_ID in the environment,
// as the environment variable docs show (INSTANCE_ID=$(uuidgen) make server).
// The server must use that ID, and every reader must get it as the *core.Uuid
// that handlers type-assert, not as the raw string.
func TestConfigureInstanceIDFromEnvironment(t *testing.T) {
	want := uuid.Must(uuid.NewV4())
	t.Setenv("INSTANCE_ID", want.String())
	v, log := newInstanceIDTestViper(t)

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
}

// TestConfigureInstanceIDGenerated covers an unset INSTANCE_ID and values that
// are not a usable UUID, which each get a newly generated ID, and a valid UUID
// with surrounding whitespace, which is trimmed and used. In every case the ID
// is stored for every reader.
func TestConfigureInstanceIDGenerated(t *testing.T) {
	padded := uuid.Must(uuid.NewV4())

	cases := []struct {
		name string
		env  string
		set  bool
		want *uuid.UUID // nil: any newly generated, non-nil ID
	}{
		{name: "unset"},
		{name: "empty", env: "", set: true},
		{name: "not a UUID", env: "my-meshery", set: true},
		{name: "nil UUID", env: uuid.Nil.String(), set: true},
		{name: "valid UUID with surrounding spaces", env: "  " + padded.String() + "\n", set: true, want: &padded},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setenv restores the variable after the test; unset it for
			// cases that need it absent, whatever the outer environment holds.
			t.Setenv("INSTANCE_ID", tc.env)
			if !tc.set {
				if err := os.Unsetenv("INSTANCE_ID"); err != nil {
					t.Fatalf("os.Unsetenv(INSTANCE_ID) failed with error: %v", err)
				}
			}
			v, log := newInstanceIDTestViper(t)

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
		})
	}
}
