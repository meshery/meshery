package helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
	models "github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

type contentandName struct {
	name    string
	content string
}
type Component string

const (
	Filter  Component = "FILTER"
	Pattern Component = "PATTERN"
)

var seededUUIDs []uuid.UUID
var localProvider *models.DefaultLocalProvider

// SeedContent seeds the components from file system into the database
func SeedContent(l *models.DefaultLocalProvider) {
	localProvider = l
	logrus.Info("Starting to seed patterns")
	components, err := getSeededComponents(Pattern)
	if err != nil {
		logrus.Error("[SEEDING] ", err.Error())
	}

	for _, p := range components {
		id, _ := uuid.NewV4()
		var pattern *models.MesheryPattern = &models.MesheryPattern{
			PatternFile: p.content,
			Name:        p.name,
			ID:          &id,
		}
		logrus.Info("[SEEDING] ", "Saving pattern- ", p.name)
		_, err := l.MesheryPatternPersister.SaveMesheryPattern(pattern)
		if err != nil {
			logrus.Error("[SEEDING] ", err.Error())
		}
		seededUUIDs = append(seededUUIDs, id)
	}
	logrus.Info("Starting to seed filters")
	components, err = getSeededComponents(Filter)
	if err != nil {
		logrus.Error("[SEEDING] ", err.Error())
	}

	for _, f := range components {
		id, _ := uuid.NewV4()
		var filter *models.MesheryFilter = &models.MesheryFilter{
			FilterFile: f.content,
			Name:       f.name,
			ID:         &id,
		}
		logrus.Info("[SEEDING] ", "Saving filter- ", f.name)
		_, err := l.MesheryFilterPersister.SaveMesheryFilter(filter)
		if err != nil {
			logrus.Error("[SEEDING] ", err.Error())
		}
		seededUUIDs = append(seededUUIDs, id)
	}

}

func getSeededComponents(comp Component) ([]contentandName, error) {
	wd, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	switch comp {
	case Pattern:
		wd = filepath.Join(wd, ".meshery", "seed_content", "patterns")
	case Filter:
		wd = filepath.Join(wd, ".meshery", "seed_content", "filters", "binaries")
	}
	logrus.Info("[SEEDING] ", "Extracting "+comp+"S from ", wd)
	var contentandNames []contentandName
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		fPath := filepath.Join(wd, f.Name())
		file, err := os.OpenFile(fPath, os.O_RDONLY, 444)
		if err != nil {
			return nil, err
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		contentandName := contentandName{
			name:    f.Name(),
			content: string(content),
		}
		contentandNames = append(contentandNames, contentandName)
	}
	return contentandNames, nil
}

// DeleteSeededContent deletes the seeded components, so that copy of same components can be avoided at restarts.
func DeleteSeededContent() {
	for _, u := range seededUUIDs {
		localProvider.MesheryFilterPersister.DeleteMesheryFilter(u)
		localProvider.MesheryPatternPersister.DeleteMesheryPattern(u)
	}
}
