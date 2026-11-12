package orderby_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/meshery/meshery/server/internal/lint/orderby"
)

// TestAnalyzer pins the rule against the ORDER BY shapes the server actually
// builds. The `// want` annotations in testdata are the assertion: every
// rejected shape must report, and - just as important - every accepted shape
// must stay silent, because a rule that fires on the existing call sites would
// be turned off rather than obeyed.
func TestAnalyzer(t *testing.T) {
	analysistest.Run(
		t,
		analysistest.TestData(),
		orderby.Analyzer,
		"persister",
		"github.com/meshery/meshery/server/models",
	)
}
