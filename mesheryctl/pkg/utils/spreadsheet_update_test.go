package utils

import (
	"reflect"
	"testing"
)

func TestProcessModelToComponentsMap(t *testing.T) {
	tests := []struct {
		name               string
		existingComponents map[string]map[string][]ComponentCSV
		expectedMap        map[string]map[string]map[string]bool
	}{
		{
			name:               "Empty input map",
			existingComponents: map[string]map[string][]ComponentCSV{},
			expectedMap:        map[string]map[string]map[string]bool{},
		},
		{
			name: "Single registrant and model",
			existingComponents: map[string]map[string][]ComponentCSV{
				"registrant1": {
					"model1": {
						{Component: "comp1"},
					},
				},
			},
			expectedMap: map[string]map[string]map[string]bool{
				"registrant1": {
					"model1": {
						"comp1": true,
					},
				},
			},
		},
		{
			name: "Multiple registrants and models",
			existingComponents: map[string]map[string][]ComponentCSV{
				"registrant1": {
					"model1": {
						{Component: "comp1"},
					},
				},
				"registrant2": {
					"model2": {
						{Component: "comp2"},
					},
				},
			},
			expectedMap: map[string]map[string]map[string]bool{
				"registrant1": {
					"model1": {
						"comp1": true,
					},
				},
				"registrant2": {
					"model2": {
						"comp2": true,
					},
				},
			},
		},
		{
			name: "Multiple components for a single model",
			existingComponents: map[string]map[string][]ComponentCSV{
				"registrant1": {
					"model1": {
						{Component: "comp1"},
						{Component: "comp2"},
					},
				},
			},
			expectedMap: map[string]map[string]map[string]bool{
				"registrant1": {
					"model1": {
						"comp1": true,
						"comp2": true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ProcessModelToComponentsMap(tt.existingComponents)
			if !reflect.DeepEqual(RegistrantToModelsToComponentsMap, tt.expectedMap) {
				t.Errorf("ProcessModelToComponentsMap() = %v, want %v", RegistrantToModelsToComponentsMap, tt.expectedMap)
			}
		})
	}
}

func TestMarshalStructToCSValues(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
	}

	tests := []struct {
		name    string
		input   []*TestStruct
		want    [][]interface{}
		wantErr bool
	}{
		{
			name:    "Empty input",
			input:   []*TestStruct{},
			want:    [][]interface{}{},
			wantErr: false,
		},
		{
			name: "Single element input",
			input: []*TestStruct{
				{Field1: "test1", Field2: 1},
			},
			want: [][]interface{}{
				{"test1", "1"},
			},
			wantErr: false,
		},
		{
			name: "Multiple elements input",
			input: []*TestStruct{
				{Field1: "test1", Field2: 1},
				{Field1: "test2", Field2: 2},
			},
			want: [][]interface{}{
				{"test1", "1"},
				{"test2", "2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshalStructToCSValues(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshalStructToCSValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("marshalStructToCSValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
