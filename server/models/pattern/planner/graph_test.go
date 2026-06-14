package planner

import (
	"testing"
)

func TestGraph_DetectCycle(t *testing.T) {
	type fields struct {
		Nodes map[string]*Node
		Edges map[string][]string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Cycle does not exists when the graph is disconnected",
			fields: fields{
				Nodes: map[string]*Node{
					"1": {},
					"2": {},
					"3": {},
					"4": {},
				},
				Edges: map[string][]string{},
			},
			want: false,
		},
		{
			name: "Cycle does not exists when the graph is connected",
			fields: fields{
				Nodes: map[string]*Node{
					"1": {},
					"2": {},
					"3": {},
					"4": {},
				},
				Edges: map[string][]string{
					"1": {"2", "3"},
					"3": {"4"},
				},
			},
			want: false,
		},
		{
			name: "Cycle exists when the graph is connected",
			fields: fields{
				Nodes: map[string]*Node{
					"1": {},
					"2": {},
					"3": {},
					"4": {},
				},
				Edges: map[string][]string{
					"1": {"2", "3"},
					"3": {"1"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGraph()

			// Add nodes
			for node, val := range tt.fields.Nodes {
				g.AddNode(node, val.Data)
			}

			// Add edges
			for src, dests := range tt.fields.Edges {
				for _, dest := range dests {
					g.AddEdge(src, dest)
				}
			}

			if got := g.DetectCycle(); got != tt.want {
				t.Errorf("Graph.DetectCycle() = %v, want %v", got, tt.want)
			}
		})
	}
}
