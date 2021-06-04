package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")
var binaryName = "mesheryctl"
var binaryPath string

func TestMain(m *testing.M) {
	// move one directory up
	err := os.Chdir(".")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}
	// create a mesheryctl binary file
	make := exec.Command("make")
	err = make.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s:%v", binaryName, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
