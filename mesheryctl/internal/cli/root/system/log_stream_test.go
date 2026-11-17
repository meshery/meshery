package system

import (
	"io"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestPrintLogStreamPreservesLinesAcrossReads(t *testing.T) {
	b := utils.SetupMeshkitLoggerTesting(t, false)
	reader := io.MultiReader(strings.NewReader("hello "), strings.NewReader("world\nlast"))
	if err := printLogStream(reader, "pod"); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "pod\t|\thello world\npod\t|\tlast\n"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
