package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

var (
	// ManifestsFolder is where the Kubernetes manifests are stored
	ManifestsFolder = "manifests"
	// MesheryFolder is the default relative location of the meshery config
	// related configuration files.
	MesheryFolder = ".meshery"

	ReleaseTag string
)

// SafeClose is a helper function help to close the io
func SafeClose(co io.Closer) {
	if cerr := co.Close(); cerr != nil {
		log.Error(cerr)
	}
}

// CreateManifestsFolder creates a new folder (.meshery/manifests)
func CreateManifestsFolder() error {
	mesheryManifestFolder := filepath.Join(MesheryFolder, ManifestsFolder)

	log.Debug("deleting " + ManifestsFolder + " folder...")
	// delete manifests folder if it already exists
	if err := os.RemoveAll(mesheryManifestFolder); err != nil {
		return err
	}
	log.Debug("creating " + ManifestsFolder + "folder...")
	// create a manifests folder under ~/.meshery to store the manifest files
	if err := os.MkdirAll(mesheryManifestFolder, os.ModePerm); err != nil {
		return errors.Wrapf(err, fmt.Sprintf("failed to make %s directory", ManifestsFolder))
	}
	log.Debug("created manifests folder...")

	return nil
}
