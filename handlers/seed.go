package handlers

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/models/pattern/core"
	"github.com/layer5io/meshery/models/pattern/stages"
	"github.com/layer5io/meshkit/utils"
)

func SeedApplications(l models.Provider) {
	names, contents, err := getSeededApplication()
	if err != nil {
		return
	}
	for i, name := range names {
		content := []byte(contents[i])

		p, err := core.NewPatternFile(content)
		if err != nil {
			return
		}
		chain := stages.CreateChain()
		sip := &serviceInfoProvider{
			provider: l,
		}
		// we only need the error field as we dont further along after the third stage
		sap := &serviceActionProvider{
			err: nil,
		}
		chain.
			Add(stages.ServiceIdentifier(sip, sap)).
			Add(stages.Filler).
			Add(stages.Validator(sip, sap)).Process(&stages.Data{
			Pattern: &p,
			Other:   map[string]interface{}{},
		})
		if sap.err == nil {
			app := models.MesheryApplication{
				Name:            name,
				ApplicationFile: string(content),
			}
			_, err := l.SaveMesheryApplication("", &app) //This function is only used for local provider, which doesn't make use of token hence it is passed as an empty string
			if err != nil {
				return
			}
		}
	}
	return

}

// getSeededComponents reads the directory recursively looking for seed content
func getSeededApplication() ([]string, []string, error) {
	wd := utils.GetHome()
	wd = filepath.Join(wd, ".meshery", "seed_content", "applications")
	var names []string
	var contents []string
	err := filepath.WalkDir(wd,
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				file, err := os.OpenFile(path, os.O_RDONLY, 0444)
				if err != nil {
					return err
				}
				content, err := ioutil.ReadAll(file)
				if err != nil {
					return err
				}
				names = append(names, d.Name())
				contents = append(contents, string(content))
			}
			return nil
		})
	if err != nil {
		return nil, nil, err
	}
	return names, contents, nil
}
