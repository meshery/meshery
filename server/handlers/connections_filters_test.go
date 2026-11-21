package handlers

import (
	"slices"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta3/environment"
)

// TestGetConnectionsRepeatedFilterParams drives the connections list through
// the real handler and Local provider. An empty filter value (`status=`,
// `kind=`, `type=`) must be ignored rather than matched literally: the
// Dashboard resource view sends `status=` and was getting an empty list back.
func TestGetConnectionsRepeatedFilterParams(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(
		connections.Connection{},
		environment.Environment{},
		environment.EnvironmentConnectionMapping{},
	); err != nil {
		t.Fatalf("migrate connection tables: %v", err)
	}

	cp := &models.ConnectionPersister{DB: &db}
	for _, c := range []struct {
		name, kind, connType string
		status               connections.ConnectionStatus
	}{
		{"k8s-connected", "kubernetes", "platform", connections.CONNECTED},
		{"k8s-disconnected", "kubernetes", "platform", connections.DISCONNECTED},
		{"prom-connected", "prometheus", "telemetry", connections.CONNECTED},
	} {
		if _, err := cp.SaveConnection(&connections.Connection{
			ID:             core.Uuid(uuid.Must(uuid.NewV4())),
			Name:           c.name,
			Kind:           c.kind,
			ConnectionType: c.connType,
			Status:         c.status,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}); err != nil {
			t.Fatalf("save %s: %v", c.name, err)
		}
	}

	h := &Handler{config: &models.HandlerConfig{}, log: newTestLogger(t)}
	provider := &models.DefaultLocalProvider{ConnectionPersister: cp}

	all := []string{"k8s-connected", "k8s-disconnected", "prom-connected"}
	k8s := []string{"k8s-connected", "k8s-disconnected"}
	connected := []string{"k8s-connected", "prom-connected"}

	for _, tc := range []struct {
		query string
		want  []string
	}{
		{"", all},
		{"status=", all},
		{"kind=", all},
		{"type=", all},
		{"status=%20", all},
		{"kind=%20%20&type=%09", all},
		{"status=connected", connected},
		{"kind=kubernetes", k8s},
		{"type=telemetry", []string{"prom-connected"}},
		{"status=&status=connected", connected},
		{"kind=kubernetes&kind=", k8s},
		{"type=&type=telemetry", []string{"prom-connected"}},
		{"status=connected&status=disconnected", all},
		{"kind=kubernetes&kind=prometheus", all},
		{"kind=kubernetes&status=connected", []string{"k8s-connected"}},
		// The exact query the Dashboard resource detail view sends.
		{"page=0&search=&order=&kind=kubernetes&status=", k8s},
	} {
		t.Run("?"+tc.query, func(t *testing.T) {
			page, body := getConnectionsAPI(t, h, provider, tc.query)

			got := make([]string, 0, len(page.Connections))
			for _, c := range page.Connections {
				got = append(got, c.Name)
			}
			slices.Sort(got)

			if !slices.Equal(got, tc.want) || page.TotalCount != len(tc.want) {
				t.Errorf("got %v (totalCount=%d), want %v; body: %s", got, page.TotalCount, tc.want, body)
			}
		})
	}
}
