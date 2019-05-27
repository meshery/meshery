package cmd

import (
	"io"
	"net/http"
	"os"
)

const (
	url                = "http://localhost:9081"
	fileUrl            = "https://raw.githubusercontent.com/layer5io/meshery/master/docker-compose.yaml"
	mesheryLocalFolder = ".meshery"
	dockerComposeFile  = mesheryLocalFolder + "/meshery.yaml"
)

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
