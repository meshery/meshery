// // Package store provides methods for interacting
// // with a thread safe global store.

// package store

// import (
// 	"reflect"
// 	"testing"
// )

// // Dummy struct
// type dummy struct{}

// func (d *dummy) SetID(id string) {}

// func TestSet(t *testing.T) {
// 	// Reset global store
// 	globalStore = newThreadSafeStore()

// 	// Initialize the store
// 	Initialize()

// 	type args struct {
// 		key   string
// 		value Value
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{name: "When value is a string", args: args{key: "k1", value: "val1"}},
// 		{name: "When value is a number", args: args{key: "k2", value: 1234}},
// 		{name: "When value is custom struct", args: args{key: "k3", value: &dummy{}}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			Set(tt.args.key, tt.args.value)

// 			if !includes(globalStore.store[tt.args.key], tt.args.value) {
// 				t.Errorf("Set() = %v, want %v", globalStore.store[tt.args.key], tt.args.value)
// 			}
// 		})
// 	}
// }

// func TestGetAll(t *testing.T) {
// 	// Reset global store
// 	globalStore = newThreadSafeStore()

// 	// Initialize the store
// 	Initialize()

// 	type args struct {
// 		key string
// 	}

// 	// Add data to the store
// 	Set("key1", 1234)
// 	Set("key1", 1235)
// 	Set("key1", 1235)
// 	Set("key2", struct{ Name string }{Name: "val1"})
// 	Set("key2", struct{ Name string }{Name: "val2"})

// 	tests := []struct {
// 		name string
// 		args args
// 		want []interface{}
// 	}{
// 		{
// 			name: "When key exists with multiple values including duplicates",
// 			args: args{
// 				key: "key1",
// 			},
// 			want: []interface{}{1234, 1235},
// 		},
// 		{
// 			name: "When key exists with multiple values without duplicates",
// 			args: args{
// 				key: "key2",
// 			},
// 			want: []interface{}{
// 				struct{ Name string }{Name: "val1"},
// 				struct{ Name string }{Name: "val2"},
// 			},
// 		},
// 		{
// 			name: "When key does not exists",
// 			args: args{
// 				key: "key3",
// 			},
// 			want: []interface{}{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := GetAll(tt.args.key); !matchSlice(got, tt.want) {
// 				t.Errorf("GetAll() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestPrefixMatch(t *testing.T) {
// 	// Reset global store
// 	globalStore = newThreadSafeStore()

// 	// Initialize the store
// 	Initialize()

// 	type args struct {
// 		key string
// 	}

// 	// Add data to the store
// 	Set("key1", 1234)
// 	Set("key1", 1235)
// 	Set("key1", 1235)
// 	Set("key2", struct{ Name string }{Name: "val1"})
// 	Set("key2", struct{ Name string }{Name: "val2"})

// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantRes []interface{}
// 	}{
// 		{
// 			name: "",
// 			args: args{
// 				key: "key",
// 			},
// 			wantRes: []interface{}{
// 				1234,
// 				1235,
// 				struct{ Name string }{Name: "val1"},
// 				struct{ Name string }{Name: "val2"},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if gotRes := PrefixMatch(tt.args.key); !matchSlice(gotRes, tt.wantRes) {
// 				t.Errorf("PrefixMatch() = %v, want %v", gotRes, tt.wantRes)
// 			}
// 		})
// 	}
// }

// func includes(s map[string]interface{}, v interface{}) bool {
// 	for _, si := range s {
// 		if reflect.DeepEqual(si, v) {
// 			return true
// 		}
// 	}

// 	return false
// }

// func matchSlice(arr1 []interface{}, arr2 []interface{}) bool {
// 	mp := map[interface{}]bool{}
// 	for _, el := range arr1 {
// 		mp[el] = true
// 	}

// 	for _, el := range arr2 {
// 		if !mp[el] {
// 			return false
// 		}
// 	}

// 	return true
// }
package store
