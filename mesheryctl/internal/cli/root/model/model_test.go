package model

import (
	"flag"
	"testing"
)

var update = flag.Bool("update", false, "Update the model")

func TestGenerateModelDataToDisplay_NilResponse(t *testing.T) {
	rows, total := generateModelDataToDisplay(nil)
	if len(rows) != 0 || total != 0 {
		t.Fatalf("expected 0 rows and 0 total on nil response, got %d rows and %d total", len(rows), total)
	}
}
