package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// newTestSeedLog redirects the seed-log root to a temp directory for the
// duration of the test and creates the dedicated log for one content type.
func newTestSeedLog(t *testing.T, process logger.Handler, stage SeedStage) *SeedLog {
	t.Helper()

	oldSeedLogsRoot := seedLogsRoot
	seedLogsRoot = func() (string, error) { return t.TempDir(), nil }
	t.Cleanup(func() { seedLogsRoot = oldSeedLogsRoot })

	seedLog, err := NewSeedLog(process, stage)
	if err != nil {
		t.Fatalf("NewSeedLog(%q): %v", stage, err)
	}
	return seedLog
}

// readSeedLog returns the full content of the dedicated seed log file.
func readSeedLog(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read seed log %s: %v", path, err)
	}
	return string(content)
}

// TestSeedLogUsesADedicatedFilePerContentType pins the per-content-type log
// contract: each seed stage owns its own file under the seed-log directory, so
// keys, models, policies and designs each have a self-contained review surface.
func TestSeedLogUsesADedicatedFilePerContentType(t *testing.T) {
	stages := []SeedStage{SeedStageKeys, SeedStageModels, SeedStagePolicies, SeedStageDesigns}
	seen := make(map[string]bool)
	for _, stage := range stages {
		t.Run(string(stage), func(t *testing.T) {
			seedLog := newTestSeedLog(t, nil, stage)
			if seedLog.Stage() != stage {
				t.Fatalf("stage mismatch: got %q want %q", seedLog.Stage(), stage)
			}
			want := filepath.Base(seedLog.Path())
			if want != stage.fileName() {
				t.Fatalf("dedicated log file for %q: got %q want %q", stage, want, stage.fileName())
			}
			if !filepath.IsAbs(seedLog.Path()) {
				t.Fatalf("seed log path %q is not absolute", seedLog.Path())
			}
			seen[want] = true
		})
	}
	if len(seen) != len(stages) {
		t.Fatalf("expected %d distinct seed log files, got %d", len(stages), len(seen))
	}
}

// TestSeedLogRouting pins where each line type goes: summary lines reach both
// the process log and the file, detail lines reach the file only, and error
// lines reach both.
func TestSeedLogRouting(t *testing.T) {
	log, sink := capturingTestLogger(t)
	seedLog := newTestSeedLog(t, log, SeedStageModels)

	seedLog.Reportf("Seeded %d models.", 3)
	seedLog.Detailf("detail that stays in the file only")
	seedLog.Errorf("registration failed for model %s", "jaegar")

	seedLog.Close()
	content := readSeedLog(t, seedLog.Path())

	if !strings.Contains(content, "Seeded 3 models.") {
		t.Fatalf("Reportf line missing from the seed log:\n%s", content)
	}
	if !strings.Contains(content, "detail that stays in the file only") {
		t.Fatalf("Detailf line missing from the seed log:\n%s", content)
	}
	if !strings.Contains(content, "registration failed for model jaegar") {
		t.Fatalf("Errorf line missing from the seed log:\n%s", content)
	}
	if got := len(sink.records(t)); got != 1 {
		t.Fatalf("expected the error line only in the captured process log, got %d records", got)
	}
}

// TestRunSeedStageWritesOutcomeAndPublishesEvent pins the UI-review contract: a
// finished stage records its header and outcome in its dedicated log file and
// publishes a system event whose metadata carries the file path as
// DownloadLink/ViewLink, which the notification center renders as "Download
// File" and "Get Logs".
func TestRunSeedStageWritesOutcomeAndPublishesEvent(t *testing.T) {
	log := registryLogTestLogger(t)
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open in-memory database: %v", err)
	}
	if err := db.AutoMigrate(events.Event{}); err != nil {
		t.Fatalf("migrate events: %v", err)
	}

	seedLog := newTestSeedLog(t, log, SeedStagePolicies)
	seedLog.WithEventSink(&EventsPersister{DB: &db}, core.Uuid{})
	RunSeedStage(log, seedLog, func(sl *SeedLog) {
		sl.Reportf("Loaded 4 policies from ./policies.")
	})

	content := readSeedLog(t, seedLog.Path())
	if !strings.Contains(content, "Seeding \"policies\" started") || !strings.Contains(content, "completed") {
		t.Fatalf("dedicated seed log lacks header or outcome:\n%s", content)
	}
	if !strings.Contains(content, "Loaded 4 policies from ./policies.") {
		t.Fatalf("stage summary missing from the dedicated seed log:\n%s", content)
	}

	persisted := persistedSystemEvents(t, &db)
	if len(persisted) != 1 {
		t.Fatalf("expected 1 persisted seeding event, got %d", len(persisted))
	}
	meta, ok := persisted[0].Metadata["DownloadLink"].(string)
	if !ok || meta != seedLog.Path() {
		t.Fatalf("DownloadLink metadata does not point at the dedicated log: %#v", persisted[0].Metadata)
	}
	if meta, ok := persisted[0].Metadata["ViewLink"].(string); !ok || meta != seedLog.Path() {
		t.Fatalf("ViewLink metadata does not point at the dedicated log: %#v", persisted[0].Metadata)
	}
	if stage, ok := persisted[0].Metadata["stage"].(string); !ok || stage != string(SeedStagePolicies) {
		t.Fatalf("stage metadata mismatch: %#v", persisted[0].Metadata)
	}
}
