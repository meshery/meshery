package pkg

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadCSV(url string) (string, error) {
	file, err := os.CreateTemp("./", "*.csv")
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	byt, _ := io.ReadAll(resp.Body)
	file.WriteString(string(byt))
	byt, _ = io.ReadAll(file)
	path, _ := filepath.Abs(file.Name())
	return path, nil
}
