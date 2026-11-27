package planner

import (
	"testing"

	"github.com/meshery/schemas/models/v1beta2/component"
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

			got, err := g.DetectCycle()
			if err != nil {
				t.Fatalf("Graph.DetectCycle() returned an unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("Graph.DetectCycle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraph_EmptyVertexName(t *testing.T) {
	g := NewGraph()
	g.AddNode("", component.ComponentDefinition{})

	cycle, err := g.DetectCycle()
	if err != nil {
		t.Fatalf("DetectCycle returned unexpected error: %v", err)
	}
	if cycle {
		t.Errorf("DetectCycle returned true for single node with empty name, want false")
	}

	var visited []string
	err = g.Traverse(func(v string, _ component.ComponentDefinition) bool {
		visited = append(visited, v)
		return true
	})
	if err != nil {
		t.Fatalf("Traverse returned unexpected error: %v", err)
	}
	if len(visited) != 1 || visited[0] != "" {
		t.Errorf("expected visited [\"\"], got %v", visited)
	}
}

func TestGraph_DetectCycle_RejectsUnknownTargetNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("nodeA", component.ComponentDefinition{})
	g.AddEdge("nodeA", "nonExistent")

	_, err := g.DetectCycle()
	if err == nil {
		t.Fatal("expected error when edge points to an unknown node, got nil")
	}
}

func BenchmarkGraph_DetectCycle(b *testing.B) {
	g := NewGraph()
	for i := 0; i < 50; i++ {
		g.AddNode(string(rune('A'+i)), component.ComponentDefinition{
			DisplayName: "benchmark-component",
		})
	}
	for i := 0; i < 49; i++ {
		g.AddEdge(string(rune('A'+i)), string(rune('A'+i+1)))
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hasCycle, err := g.DetectCycle()
		if err != nil || hasCycle {
			b.Fatalf("unexpected result: cycle=%v, err=%v", hasCycle, err)
		}
	}
}
