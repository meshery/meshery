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
	MockOutput       string
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

// Path to the current file
func GetBasePath(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("problems recovering caller information")
	}

	return filepath.Dir(filename)
}

// Load a Golden file
func (tf *GoldenFile) Load() string {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return string(content)
}

// write a Golden file
func (tf *GoldenFile) Write(content string) {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}
