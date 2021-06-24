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

func TestMain(t *testing.T) {
	// create a mesheryctl binary file
	make := exec.Command("make")
	err := make.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s:%v", binaryName, err)
		os.Exit(1)
	}
}
