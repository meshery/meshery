// Command orderbylint runs the orderby analyzer over the packages named on the
// command line and exits non-zero when any gorm ORDER BY clause is built from a
// value that is neither constant nor sanitized.
//
// It runs as a step of the golangci-lint-server CI job and of `make golangci`,
// so a violation fails the same gate that already blocks the build. golangci-lint
// cannot host this check without building a custom binary from its module plugin
// system, which is more machinery than one analyzer is worth; see
// docs/content/en/project/contributing/contributing-lint.md.
//
//	go run ./server/internal/lint/orderby/cmd/orderbylint ./...
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/meshery/meshery/server/internal/lint/orderby"
)

// main runs the analyzer over the packages named on the command line, exiting
// non-zero when any of them reports.
func main() {
	singlechecker.Main(orderby.Analyzer)
}
