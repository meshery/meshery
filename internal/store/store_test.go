// Package store provides methods for interacting
// with a thread safe global store.

package store

import (
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	// Initialize the store
	Initialize()

	// Dummy struct
	type dummy struct{}

	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "When value is a string", args: args{key: "k1", value: "val1"}},
		{name: "When value is a number", args: args{key: "k2", value: 1234}},
		{name: "When value is custom struct", args: args{key: "k3", value: dummy{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.args.key, tt.args.value)

			if !includes(globalStore.store[tt.args.key], tt.args.value) {
				t.Errorf("Set() = %v, want %v", globalStore.store[tt.args.key], tt.args.value)
			}
		})
	}
}

func includes(s []interface{}, v interface{}) bool {
	for _, si := range s {
		if reflect.DeepEqual(si, v) {
			return true
		}
	}

	return false
}
