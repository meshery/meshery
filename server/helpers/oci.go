package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/opencontainers/image-spec/specs-go"
)

func main() {
	// Create a sample JSON file 1
	file1Content := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	file1Data, err := json.MarshalIndent(file1Content, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON for file1:", err)
		os.Exit(1)
	}
	err = ioutil.WriteFile("file1.json", file1Data, 0644)
	if err != nil {
		fmt.Println("Error writing file1:", err)
		os.Exit(1)
	}

	// Create a sample JSON file 2
	file2Content := map[string]interface{}{
		"key3": "value3",
		"key4": true,
	}
	file2Data, err := json.MarshalIndent(file2Content, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON for file2:", err)
		os.Exit(1)
	}
	err = ioutil.WriteFile("file2.json", file2Data, 0644)
	if err != nil {
		fmt.Println("Error writing file2:", err)
		os.Exit(1)
	}

	// Create the OCI manifest
	manifest := &specs.Manifest{
		Config: specs.Descriptor{
			MediaType: "application/vnd.opencontainers.image.config.v1+json",
			Digest:    "sha256:" + "<SHA256_SUM>", // Replace with the actual digest of your config file
			Size:      int64(len(file1Data)),
		},
		Layers: []specs.Descriptor{
			{
				MediaType: "application/vnd.opencontainers.image.layer.v1.tar+gzip",
				Digest:    "sha256:" + "<SHA256_SUM>", // Replace with the actual digest of your file1 layer
				Size:      int64(len(file1Data)),
			},
			{
				MediaType: "application/vnd.opencontainers.image.layer.v1.tar+gzip",
				Digest:    "sha256:" + "<SHA256_SUM>", // Replace with the actual digest of your file2 layer
				Size:      int64(len(file2Data)),
			},
		},
	}

	// Encode the manifest as JSON
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Println("Error encoding OCI manifest:", err)
		os.Exit(1)
	}

	// Write the manifest to a file
	err = ioutil.WriteFile("oci-manifest.json", manifestData, 0644)
	if err != nil {
		fmt.Println("Error writing OCI manifest:", err)
		os.Exit(1)
	}

	fmt.Println("OCI package created successfully.")
}