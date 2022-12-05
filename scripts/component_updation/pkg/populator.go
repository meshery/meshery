package pkg

import (
	"io/fs"
	"path/filepath"
)

type jsonpopulator func(filepath string, changeFields map[string]string) error

// PopulateEntries goes through all entries and for the key=entryAsFileName,it will find all files recursively with this name and populate it
func PopulateEntries(dir string, entries []map[string]string, entryAsFileName string, populator jsonpopulator) error {
	//Step 1: Create a map of fileNames->[filePaths]
	dirNames := make(map[string][]string)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			name := d.Name()
			if dirNames[name] == nil {
				dirNames[name] = make([]string, 0)
			}
			dirNames[name] = append(dirNames[name], path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, entry := range entries {
		dirPaths := dirNames[entry[entryAsFileName]]
		for _, fp := range dirPaths {
			err = populator(fp, entry)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
