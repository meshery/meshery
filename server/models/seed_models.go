package models

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshkit/logger"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/registration"
)

var ModelsPath = "../meshmodel"

// Define the maximum number of worker goroutines
const maxWorkers = 20

func getModelDirectoryPaths() ([]string, error) {

	dirEntries := make([]string, 0)
	modelsDirs, err := os.ReadDir(ModelsPath)
	if err != nil {
		return dirEntries, err
	}

	// Channel to collect directory paths
	dirPathsChan := make(chan string, 100)
	// Channel to collect errors
	errorChan := make(chan error, 10)
	// Wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Function to process a single directory
	processDirectory := func(modelDir os.DirEntry) {
		defer wg.Done()
		if !modelDir.IsDir() {
			return
		}
		modelVersionsDirPath := filepath.Join(ModelsPath, modelDir.Name())
		modelVersionsDir, err := os.ReadDir(modelVersionsDirPath)
		if err != nil {
			errorChan <- err
			return
		}
		for _, version := range modelVersionsDir {
			if !version.IsDir() {
				continue
			}
			modelDefVersionsDirPath := filepath.Join(modelVersionsDirPath, version.Name())
			modelDefVersionsDir, err := os.ReadDir(modelDefVersionsDirPath)
			if err != nil {
				errorChan <- err
				return
			}
			for _, defVersion := range modelDefVersionsDir {
				if !defVersion.IsDir() {
					continue
				}
				defPath := filepath.Join(modelDefVersionsDirPath, defVersion.Name())
				// Send the directory path to the channel
				dirPathsChan <- defPath
			}
		}
	}

	// Start a goroutine for each top-level directory
	for _, modelDir := range modelsDirs {
		wg.Add(1)
		go processDirectory(modelDir)
	}

	// Close the directory paths channel once all directories are processed
	go func() {
		wg.Wait()
		close(dirPathsChan)
		close(errorChan)
	}()

	// Collect all directory paths and errors from the channels
	for dirPathsChan != nil || errorChan != nil {
		select {
		case dirPath, ok := <-dirPathsChan:
			if !ok {
				dirPathsChan = nil
			} else {
				dirEntries = append(dirEntries, dirPath)
			}
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
			} else {
				// Return error immediately
				return dirEntries, err
			}
		}
	}
	return dirEntries, nil
}

func SeedComponents(log logger.Handler, hc *HandlerConfig, regm *meshmodel.RegistryManager) {

	regErrorStore := NewRegistrationFailureLogHandler()
	regHelper := registration.NewRegistrationHelper(utils.UI, regm, regErrorStore)
	modelDirPaths, err := getModelDirectoryPaths()
	if err != nil {
		log.Error(ErrSeedingComponents(err))
		return
	}

	// Create a buffered channel to control the number of concurrent tasks
	tasks := make(chan string, len(modelDirPaths))
	// Create a wait group to wait for all worker goroutines to complete
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for dirPath := range tasks {
				dir := registration.NewDir(dirPath)
				regHelper.Register(dir) // Errors are handled inside Register
			}
		}()
	}

	// Distribute tasks to the worker pool
	for _, dirPath := range modelDirPaths {
		tasks <- dirPath
	}

	// Close the tasks channel and wait for all workers to complete
	close(tasks)
	wg.Wait()

	// Call RegistryLog after all registrations are complete
	RegistryLog(log, hc, regm, regErrorStore)
}
