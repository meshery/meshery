package pkg

import (
	"log"
	"fmt"
	"io"
	"os"
	"strings"
)

type SystemType int

const (
	Meshery SystemType = iota
	Docs
	RemoteProvider
	rowIndex = 1
	shouldRegisterColIndex = -1
)

func (dt SystemType) String() string {
	switch dt {
	case Meshery:
		return "meshery"

	case Docs:
		return "docs"

	case RemoteProvider:
		return "remote-provider"
	}
	return ""
}

func GetIndexForRegisterCol(cols []string, shouldRegister string) int {
	if shouldRegisterColIndex != -1 {
		return shouldRegisterColIndex
	}

	for index, col := range cols {
		if col == shouldRegister {
			return index
		}
	}
	return shouldRegisterColIndex
}

func formatName(input string) string {
	formattedName := strings.ReplaceAll(input, " ", "-")
	formattedName = strings.ToLower(formattedName)
	return formattedName
}

func CreateFiles(path, filetype, name, content string) {
	formattedName := formatName(name)

	fullPath := path + "/" + formattedName + filetype

	file, err := os.Create(fullPath)
	fmt.Println("Creating file:", fullPath)
	if err != nil {
		log.Println("Error creating filetype file:", err)
		return
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		log.Println("Error writing to filetype file:", err)
	}
}

