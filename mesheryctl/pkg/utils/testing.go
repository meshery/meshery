package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/kr/pretty"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/spf13/viper"
)

type TestHelper struct {
	Version string
	BaseURL string
}

func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		Version: "v0.5.10",
		BaseURL: "http://localhost:9081",
	}
}

type CmdTestInput struct {
	Name             string
	Args             []string
	ExpectedResponse string
}

type GoldenFile struct {
	t    *testing.T
	name string
	dir  string
}

func NewGoldenFile(t *testing.T, name string, directory string) *GoldenFile {
	return &GoldenFile{t: t, name: name, dir: directory}
}

// equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
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

// Load a Golden file
func (tf *GoldenFile) LoadByte() []byte {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return content
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

// write a Golden file
func (tf *GoldenFile) WriteInByte(content []byte) {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}

func SetupContextEnv(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("unable to locate meshery directory")
	}
	viper.Reset()
	viper.SetConfigFile(path + "/../../../../pkg/utils/TestConfig.yaml")
	//fmt.Println(viper.ConfigFileUsed())
	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	_, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}
}

func StartMockery(t *testing.T) {
	// activate http mocking
	httpmock.Activate()
}

func StopMockery(t *testing.T) {
	httpmock.DeactivateAndReset()
}
