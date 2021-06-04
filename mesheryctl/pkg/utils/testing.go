package utils

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kr/pretty"
)

type CmdTestInput struct {
	Name             string
	Args             []string
	ExpectedResponse string
	MockResponse     string
}

type GoldenFile struct {
	t    *testing.T
	name string
	dir  string
}

func NewGoldenFile(t *testing.T, name string, directory string) *GoldenFile {
	return &GoldenFile{t: t, name: name, dir: directory}
}

// Print difference
func Diff(expected, actual interface{}) []string {
	return pretty.Diff(expected, actual)
}

// Path to Golden file
func (tf *GoldenFile) path() string {
	tf.t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		tf.t.Fatal("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
}

// Load a testing file
func (tf *GoldenFile) Load() string {
	tf.t.Helper()

	content, err := ioutil.ReadFile(tf.path())
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return string(content)
}
