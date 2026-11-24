package models

import (
	"testing"

	"github.com/meshery/meshery/server/models/connections"
)

func TestConnectionTypeFilterUsesPersistedType(t *testing.T) {
	p := &ConnectionPersister{DB: newMigratedDB(t)}
	for _, c := range []*connections.Connection{
		{Name: "telemetry", Kind: "telemetry", ConnectionType: "telemetry", SubType: "metrics"},
		{Name: "source", Kind: "source", ConnectionType: "source", SubType: "git"},
	} {
		if _, err := p.SaveConnection(c); err != nil {
			t.Fatal(err)
		}
	}
	for _, filter := range []string{"type telemetry", "connection_type telemetry", `{"type":"telemetry"}`, `{"type":"telemetry","sub_type":"metrics"}`} {
		t.Run(filter, func(t *testing.T) {
			page, err := p.GetConnections("", "", 0, 10, filter, nil, nil, nil, "")
			if err != nil {
				t.Fatal(err)
			}
			if len(page.Connections) != 1 || page.Connections[0].Name != "telemetry" {
				t.Fatalf("filter %s: %#v", filter, page.Connections)
			}
		})
	}
}
