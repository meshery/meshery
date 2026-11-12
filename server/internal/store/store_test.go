// Package store provides methods for interacting
// with a thread safe global store.

package store

import (
	"reflect"
	"testing"
)

// Dummy struct
type dummyValue struct {
	ID    string
	Value interface{}
}

func (d *dummyValue) SetID(id string) { d.ID = id }
func (d dummyValue) GetID() string    { return d.ID }

func TestSet(t *testing.T) {
	// Reset global store
	globalStore = newThreadSafeStore()

	// Initialize the store
	Initialize()

	type args struct {
		key   string
		value Value
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "When value is an int", args: args{key: "k2", value: &dummyValue{Value: 123}}},
		{name: "When value is a string", args: args{key: "k3", value: &dummyValue{Value: "123"}}},
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

func TestGetAll(t *testing.T) {
	// Reset global store
	globalStore = newThreadSafeStore()

	// Initialize the store
	Initialize()

	type args struct {
		key string
	}

	// Add data to the store
	Set("key1", &dummyValue{Value: 12345})
	Set("key1", &dummyValue{Value: 12345})
	Set("key2", &dummyValue{Value: struct{ Name string }{Name: "val1"}})
	Set("key2", &dummyValue{Value: struct{ Name string }{Name: "val2"}})

	tests := []struct {
		name string
		args args
		want []Value
	}{
		{
			name: "When key exists with multiple values including duplicates",
			args: args{
				key: "key1",
			},
			want: []Value{&dummyValue{Value: 12345}},
		},
		{
			name: "When key exists with multiple values without duplicates",
			args: args{
				key: "key2",
			},
			want: []Value{
				&dummyValue{Value: struct{ Name string }{Name: "val1"}},
				&dummyValue{Value: struct{ Name string }{Name: "val2"}},
			},
		},
		{
			name: "When key does not exists",
			args: args{
				key: "key3",
			},
			want: []Value{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAll(tt.args.key); !matchSlice(got, tt.want) {
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
	Set("key1", &dummyValue{Value: 1234})
	Set("key1", &dummyValue{Value: 1235})
	Set("key1", &dummyValue{Value: 1235})
	Set("key2", &dummyValue{Value: struct{ Name string }{Name: "val1"}})
	Set("key2", &dummyValue{Value: struct{ Name string }{Name: "val2"}})

	tests := []struct {
		name    string
		args    args
		wantRes []Value
	}{
		{
			name: "",
			args: args{
				key: "key",
			},
			wantRes: []Value{
				&dummyValue{Value: 1234},
				&dummyValue{Value: 1235},
				&dummyValue{Value: struct{ Name string }{Name: "val1"}},
				&dummyValue{Value: struct{ Name string }{Name: "val2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := PrefixMatch(tt.args.key); !matchSlice(gotRes, tt.wantRes) {
				t.Errorf("PrefixMatch() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func includes(s map[string]Value, v Value) bool {
	for _, si := range s {
		if reflect.DeepEqual(si, v) {
			return true
		}
	}

	return false
}

func matchSlice(arr1 []Value, arr2 []Value) bool {
	mp := map[string]bool{}
	for _, el := range arr1 {
		mp[el.GetID()] = true
	}

	for _, el := range arr2 {
		id := md5Hash(el)
		if !mp[id] {
			return false
		}
	}

	return true
}
