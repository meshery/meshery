package models

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	gofrs "github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	meshkitUtils "github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/viper"
)

// SeedStage names one content type whose seed content initializes the database.
// Each stage is a bootstrapping concern and gets a dedicated log file, so that
// a boot ("Seeding <stage> started") and its outcome can be reviewed per content
// type - in the aggregate stdout report and, through the per-stage events, in
// the Meshery UI.
type SeedStage string

const (
	SeedStageKeys     SeedStage = "keys"
	SeedStageModels   SeedStage = "models"
	SeedStagePolicies SeedStage = "policies"
	SeedStageDesigns  SeedStage = "designs"
)

func (s SeedStage) fileName() string {
	return fmt.Sprintf("%s.log", s)
}

// Seed status words written to the dedicated log files and the aggregate stdout
// report.
const (
	SeedStatusSuccess = "completed"
	SeedStatusFailed  = "failed"
)

// seedLogsRoot resolves the directory that holds the per-content-type seeding
// logs (<home>/.meshery/logs/seed). It is a var, not a const, so tests can
// redirect it to a temporary directory without mutating HOME.
//
// The logs live under ~/.meshery/logs so the existing /api/system/fileView and
// /api/system/fileDownload endpoints (confined to that directory by
// SafeOpenFile) can serve them for review in the Meshery UI.
var seedLogsRoot = func() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = meshkitUtils.GetHome()
	}
	if home == "" {
		return "", fmt.Errorf("could not resolve the user home directory to place seed logs under")
	}
	return filepath.Join(home, ".meshery", "logs", "seed"), nil
}

// SeedLog is the per-content-type seeding log. It writes each line to a
// dedicated file under ~/.meshery/logs/seed, forwards summary lines to the
// process (stdout) log so the aggregate report carries one line per content
// type, and - once wired with an event sink - publishes a system event whose
// metadata lets the Meshery UI open the file through the existing
// /api/system/fileView and /api/system/fileDownload endpoints.
//
// A SeedLog is always safe to use, even when its backing file could not be
// created: NewSeedLog reports the failure once, and logging then degrades to
// the process log only.
type SeedLog struct {
	stage     SeedStage
	process   logger.Handler
	file      *os.File
	path      string
	startedAt time.Time
	errCount  int

	persister SystemEventPersister
	systemID  core.Uuid
}

// NewSeedLog creates the dedicated seed log for one content type.
func NewSeedLog(process logger.Handler, stage SeedStage) (*SeedLog, error) {
	sl := &SeedLog{
		stage:     stage,
		process:   process,
		startedAt: time.Now(),
	}
	dir, err := seedLogsRoot()
	if err != nil {
		return sl, ErrCreatingSeedLog(stage, err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return sl, ErrCreatingSeedLog(stage, err)
	}
	sl.path = filepath.Join(dir, stage.fileName())
	file, err := os.Create(sl.path)
	if err != nil {
		return sl, ErrCreatingSeedLog(stage, err)
	}
	sl.file = file
	return sl, nil
}

// WithEventSink wires a system-event persister so a finished stage can publish
// an event pointing at its dedicated log file for review in the Meshery UI.
func (sl *SeedLog) WithEventSink(persister SystemEventPersister, systemID core.Uuid) *SeedLog {
	sl.persister = persister
	sl.systemID = systemID
	return sl
}

// NewSeedLogForSystem creates a seed log wired to the system event persister,
// publishing its completed-stage event on behalf of the Meshery instance
// (INSTANCE_ID), which is how every boot-time and reset-path seed stage is run.
func NewSeedLogForSystem(process logger.Handler, persister SystemEventPersister, stage SeedStage) *SeedLog {
	seedLog, err := NewSeedLog(process, stage)
	if err != nil {
		process.Error(err)
	}
	return seedLog.WithEventSink(persister, gofrs.FromStringOrNil(viper.GetString("INSTANCE_ID")))
}

// Stage returns the content type this log covers.
func (sl *SeedLog) Stage() SeedStage { return sl.stage }

// Path returns the absolute path of the dedicated log file, or "" when the file
// could not be created.
func (sl *SeedLog) Path() string { return sl.path }

// Reportf appends a summary line to the dedicated file and forwards it to the
// process (stdout) log, so the aggregate report carries it for this content
// type.
func (sl *SeedLog) Reportf(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	sl.write(line)
	if sl.process != nil {
		sl.process.Info(line)
	}
}

// Detailf appends a line to the dedicated file only, keeping the aggregate
// stdout report light.
func (sl *SeedLog) Detailf(format string, args ...interface{}) {
	sl.write(fmt.Sprintf(format, args...))
}

// Errorf appends an error line to the dedicated file and forwards it to the
// process log. Recording an error marks the stage as failed: the aggregate
// outcome line and the published event use the failure status, so a stage that
// ran into trouble is never reported as completed.
func (sl *SeedLog) Errorf(format string, args ...interface{}) {
	sl.errCount++
	line := fmt.Sprintf(format, args...)
	sl.write(line)
	if sl.process != nil {
		sl.process.Error(fmt.Errorf("%s", line))
	}
}

// HasErrors reports whether the stage recorded any error line. RunSeedStage
// uses it to decide between the completed and failed outcome.
func (sl *SeedLog) HasErrors() bool { return sl.errCount > 0 }

// Header writes the stage header into the dedicated file.
func (sl *SeedLog) Header() {
	sl.write(fmt.Sprintf("Seeding %q started at %s.", sl.stage, sl.startedAt.Format(time.RFC3339)))
}

// Close flushes and closes the dedicated log file.
func (sl *SeedLog) Close() {
	if sl.file == nil {
		return
	}
	_ = sl.file.Sync()
	_ = sl.file.Close()
	sl.file = nil
}

func (sl *SeedLog) write(line string) {
	if sl.file != nil {
		_, _ = sl.file.WriteString(line + "\n")
	}
}

// publishEvent persists a system event describing the finished seed stage. When
// the stage produced a dedicated log file, the metadata carries DownloadLink and
// ViewLink so the notification center renders "Download File" and "Get Logs",
// letting the file be opened and reviewed in the UI. Best-effort: a persistence
// failure is returned, never swallowed silently.
func (sl *SeedLog) publishEvent(status string) error {
	if sl.persister == nil {
		return nil
	}
	severity := events.Informational
	if status != SeedStatusSuccess {
		severity = events.Error
	}
	eventBuilder := events.NewEvent().FromSystem(sl.systemID).FromOwner(sl.systemID).
		WithCategory("seeding").WithAction("get_summary").
		WithSeverity(severity).WithDescription(fmt.Sprintf("Seeding %q %s.", sl.stage, status))
	if sl.path != "" {
		eventBuilder.WithMetadata(map[string]interface{}{
			"stage":        string(sl.stage),
			"logFile":      sl.path,
			"ViewLink":     sl.path,
			"DownloadLink": sl.path,
		})
	}
	return sl.persister.PersistSystemEvent(*eventBuilder.Build())
}
