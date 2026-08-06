package academy

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/experimental/academy/templates"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"gopkg.in/yaml.v3"
)

type TemplateData struct {
	Title       string
	Description string
	Level       string
	Weight      int
	OrgID       string
	Category    string
	Tags        []string
	ID          string
}

// Frontmatter representation to extract weight
type Frontmatter struct {
	Weight int `yaml:"weight"`
}

func extractWeight(content []byte) int {
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) >= 3 {
		var fm Frontmatter
		err := yaml.Unmarshal(parts[1], &fm)
		if err == nil {
			return fm.Weight
		}
	}
	return 0
}

func inferWeight(parentDir, excludeDir string) int {
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return 1
	}

	maxWeight := 0
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != excludeDir {
			indexPath := filepath.Join(parentDir, entry.Name(), "_index.md")
			content, err := os.ReadFile(indexPath)
			if err == nil {
				weight := extractWeight(content)
				if weight > maxWeight {
					maxWeight = weight
				}
			}
		}
	}
	if maxWeight == 0 {
		return 1
	}
	return maxWeight + 1
}

func getTemplateString(contentType ContentType) string {
	switch contentType {
	case LearningPath:
		return templates.LearningPathTemplate
	case Course:
		return templates.CourseTemplate
	case Module:
		return templates.ModuleTemplate
	case Page:
		return templates.PageTemplate
	case Certification:
		return templates.CertificationTemplate
	case Exam:
		return templates.ExamTemplate
	case Lab:
		return templates.LabTemplate
	case Test:
		return templates.TestTemplate
	default:
		return ""
	}
}

func contentDirSegment(cType ContentType) string {
	return string(cType) + "s"
}

type ParentFrontmatter struct {
	Type     ContentType `yaml:"type"`
	Level    string      `yaml:"level"`
	Category string      `yaml:"categories"`
	Tags     []string    `yaml:"tags"`
}

func checkNesting(cType ContentType, parentDir string) (ParentFrontmatter, error) {
	var pf ParentFrontmatter
	indexPath := filepath.Join(parentDir, "_index.md")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		return pf, nil
	}

	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) >= 3 {
		if err := yaml.Unmarshal(parts[1], &pf); err == nil {
			allowed, exists := AllowedChildren[pf.Type]
			if exists {
				for _, child := range allowed {
					if child == cType {
						return pf, nil
					}
				}
				return pf, errInvalidNesting(pf.Type, cType)
			}
			return pf, nil
		}
	}
	return pf, nil
}

func scaffoldNode(cType ContentType, title, description, level, orgId, category string, tags []string, targetDir string, force bool, id string) error {
	tmplStr := getTemplateString(cType)
	if tmplStr == "" {
		return errTaxonomyType(string(cType))
	}

	pf, err := checkNesting(cType, targetDir)
	if err != nil {
		return err
	}
	parentType := pf.Type

	if level == "" && pf.Level != "" {
		level = pf.Level
	}
	if category == "" && pf.Category != "" {
		category = pf.Category
	}
	if len(tags) == 0 && len(pf.Tags) > 0 {
		tags = pf.Tags
	}

	if id == "" && (cType == LearningPath || cType == Certification) {
		id = "REPLACE_WITH_INSTRUCTOR_CONSOLE_ID"
	}

	var indexPath string
	folderName := strings.ToLower(strings.ReplaceAll(title, " ", "-"))

	if cType == Test && (parentType == Course || parentType == Module) {
		indexPath = filepath.Join(targetDir, "test.md")
	} else if cType == Exam && parentType == Course {
		indexPath = filepath.Join(targetDir, "course-exam.md")
	} else {
		if cType == Test && parentType == Certification {
			testNum := 1
			for {
				testFolderName := fmt.Sprintf("test-%d", testNum)
				if _, err := os.Stat(filepath.Join(targetDir, testFolderName)); os.IsNotExist(err) {
					folderName = testFolderName
					break
				}
				testNum++
			}
		}

		nodeDir := filepath.Join(targetDir, folderName)
		if err := os.MkdirAll(nodeDir, 0755); err != nil {
			return err
		}
		indexPath = filepath.Join(nodeDir, "_index.md")
	}
	if _, err := os.Stat(indexPath); err == nil && !force {
		return errScaffoldExists(indexPath)
	}

	weight := inferWeight(targetDir, folderName)

	tmpl, err := template.New(string(cType)).Parse(tmplStr)
	if err != nil {
		return err
	}

	data := TemplateData{
		Title:       title,
		Description: description,
		Level:       level,
		Weight:      weight,
		OrgID:       orgId,
		Category:    category,
		Tags:        tags,
		ID:          id,
	}

	f, err := os.Create(indexPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return err
	}

	utils.Log.Infof("Scaffolded %s '%s' at %s", cType, title, indexPath)
	return nil
}

func scaffoldTree(cType ContentType, title, description, level, orgId, category string, tags []string, targetDir string, force bool, id string) error {
	folderName := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	baseDir := filepath.Join(targetDir, folderName)

	err := scaffoldNode(cType, title, description, level, orgId, category, tags, targetDir, force, id)
	if err != nil {
		return err
	}

	currentDir := baseDir

	if cType == LearningPath {
		courseTitle := "Course 1"
		err = scaffoldNode(Course, courseTitle, "", "", orgId, "", nil, currentDir, force, "")
		if err != nil {
			return err
		}
		currentDir = filepath.Join(currentDir, strings.ToLower(strings.ReplaceAll(courseTitle, " ", "-")))
	}

	if cType == LearningPath || cType == Course {
		moduleTitle := "Module 1"
		err = scaffoldNode(Module, moduleTitle, "", "", orgId, "", nil, currentDir, force, "")
		if err != nil {
			return err
		}
		currentDir = filepath.Join(currentDir, strings.ToLower(strings.ReplaceAll(moduleTitle, " ", "-")))
	}

	if cType == LearningPath || cType == Course || cType == Module {
		pageTitle := "Page 1"
		err = scaffoldNode(Page, pageTitle, "", "", orgId, "", nil, currentDir, force, "")
		if err != nil {
			return err
		}
	}

	if cType == Certification {
		examTitle := "Exam 1"
		err = scaffoldNode(Exam, examTitle, "", "", orgId, "", nil, currentDir, force, "")
		if err != nil {
			return err
		}
	}

	return nil
}
