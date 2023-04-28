package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Hello, WebAssembly!")
	logrus.Info("Hello WebAssemly from external package")
}
