// Package store provides methods for interacting
// with a thread safe global store.

package store

import (
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	// Reset global store
	globalStore = newThreadSafeStore()

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

func TestGet(t *testing.T) {
	// Reset global store
	globalStore = newThreadSafeStore()

	// Initialize the store
	Initialize()

	type args struct {
		key string
	}

	// Add data to the store
	Set("key1", 1234)
	Set("key1", 1234)
	Set("key2", struct{ name string }{name: "val1"})

	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{
			name: "When more than one value against a key exists",
			args: args{
				key: "key1",
			},
			want:  1234,
			want1: true,
		},
		{
			name: "When only one value against a key exists",
			args: args{
				key: "key2",
			},
			want:  struct{ name string }{name: "val1"},
			want1: true,
		},
		{
			name: "When key does not exits",
			args: args{
				key: "key3",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	// Reset global store
	globalStore = newThreadSafeStore()

	// Initialize the store
	Initialize()

	type args struct {
		key string
	}

	// Add data to the store
	Set("key1", 1234)
	Set("key1", 1235)
	Set("key1", 1235)
	Set("key2", struct{ Name string }{Name: "val1"})
	Set("key2", struct{ Name string }{Name: "val2"})

	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "When key exists with multiple values including duplicates",
			args: args{
				key: "key1",
			},
			want: []interface{}{1234, 1235},
		},
		{
			name: "When key exists with multiple values without duplicates",
			args: args{
				key: "key2",
			},
			want: []interface{}{
				struct{ Name string }{Name: "val1"},
				struct{ Name string }{Name: "val2"},
			},
		},
		{
			name: "When key does not exists",
			args: args{
				key: "key3",
			},
			want: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAll(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrefixMatch(t *testing.T) {
	// Reset global store
	globalStore = newThreadSafeStore()

	// Initialize the store
	Initialize()

	type args struct {
		key string
	}

	// Add data to the store
	Set("key1", 1234)
	Set("key1", 1235)
	Set("key1", 1235)
	Set("key2", struct{ Name string }{Name: "val1"})
	Set("key2", struct{ Name string }{Name: "val2"})

	tests := []struct {
		name    string
		args    args
		wantRes []interface{}
	}{
		{
			name: "",
			args: args{
				key: "key",
			},
			wantRes: []interface{}{
				1234,
				1235,
				struct{ Name string }{Name: "val1"},
				struct{ Name string }{Name: "val2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := PrefixMatch(tt.args.key); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("PrefixMatch() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func includes(s map[string]interface{}, v interface{}) bool {
	for _, si := range s {
		if reflect.DeepEqual(si, v) {
			return true
		}
	}

	return false
}
