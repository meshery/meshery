package utils

// import (
// 	"bufio"
// 	"net/http"
// 	"reflect"
// 	"testing"
// )

// func TestConvertRespToSSE(t *testing.T) {
// 	type args struct {
// 		resp *http.Response
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    chan Event
// 		wantErr bool
// 	}{
// 		{}
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := ConvertRespToSSE(tt.args.resp)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ConvertRespToSSE() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ConvertRespToSSE() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_loop(t *testing.T) {
// 	type args struct {
// 		reader *bufio.Reader
// 		events chan Event
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			loop(tt.args.reader, tt.args.events)
// 		})
// 	}
// }
