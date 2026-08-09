package academy

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/experimental/academy/templates"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	academyModel "github.com/meshery/schemas/models/v1beta3/academy"
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

func getTemplateString(cType string) string {
	switch cType {
	case string(academyModel.LearningPath):
		return templates.LearningPathTemplate
	case string(Course):
		return templates.CourseTemplate
	case string(Module):
		return templates.ModuleTemplate
	case string(Page):
		return templates.PageTemplate
	case string(academyModel.Certification):
		return templates.CertificationTemplate
	case string(Exam):
		return templates.ExamTemplate
	case string(Lab):
		return templates.LabTemplate
	case string(Test):
		return templates.TestTemplate
	case string(academyModel.Challenge):
		return templates.ChallengeTemplate
	default:
		return ""
	}
}

func contentDirSegment(cType string) string {
	return cType + "s"
}

type ParentFrontmatter struct {
	Type     string   `yaml:"type"`
	Level    string   `yaml:"level"`
	Category string   `yaml:"categories"`
	Tags     []string `yaml:"tags"`
}

func checkNesting(cType string, parentDir string) (ParentFrontmatter, error) {
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

type ScaffoldOptions struct {
	Type        string
	Title       string
	Description string
	Level       string
	OrgID       string
	Category    string
	Tags        []string
	TargetDir   string
	Force       bool
	ID          string
}

func scaffoldNode(opts ScaffoldOptions, explicitFolderName string) error {
	tmplStr := getTemplateString(opts.Type)
	if tmplStr == "" {
		return errTaxonomyType(string(opts.Type))
	}

	pf, err := checkNesting(opts.Type, opts.TargetDir)
	if err != nil {
		return err
	}
	parentType := pf.Type

	if opts.Level == "" && pf.Level != "" {
		opts.Level = pf.Level
	}
	if opts.Category == "" && pf.Category != "" {
		opts.Category = pf.Category
	}
	if len(opts.Tags) == 0 && len(pf.Tags) > 0 {
		opts.Tags = pf.Tags
	}

	if opts.ID == "" && (opts.Type == string(academyModel.LearningPath) || opts.Type == string(academyModel.Certification) || opts.Type == string(academyModel.Challenge)) {
		opts.ID = "REPLACE_WITH_INSTRUCTOR_CONSOLE_ID"
	}

	var indexPath string
	folderName := explicitFolderName
	if folderName == "" {
		var err error
		folderName, err = makeSlug(opts.Title)
		if err != nil {
			return err
		}
	}

	if opts.Type == string(Test) && (parentType == string(Course) || parentType == string(Module)) {
		indexPath = filepath.Join(opts.TargetDir, "test.md")
	} else if opts.Type == string(Exam) && parentType == string(Course) {
		indexPath = filepath.Join(opts.TargetDir, "course-exam.md")
	} else {
		if opts.Type == string(Test) && parentType == string(academyModel.Certification) {
			const maxTests = 1000
			found := false
			for testNum := 1; testNum <= maxTests; testNum++ {
				testFolderName := fmt.Sprintf("test-%d", testNum)
				_, statErr := os.Stat(filepath.Join(opts.TargetDir, testFolderName))
				if statErr == nil {
					continue
				}
				if !os.IsNotExist(statErr) {
					return statErr
				}
				folderName = testFolderName
				found = true
				break
			}
			if !found {
				return errScaffoldExists(filepath.Join(opts.TargetDir, "test-*"))
			}
		}

		nodeDir := filepath.Join(opts.TargetDir, folderName)
		if err := os.MkdirAll(nodeDir, 0755); err != nil {
			return err
		}
		indexPath = filepath.Join(nodeDir, "_index.md")
	}
	if _, err := os.Stat(indexPath); err == nil && !opts.Force {
		return errScaffoldExists(indexPath)
	}

	weight := inferWeight(opts.TargetDir, folderName)

	tmpl, err := template.New(string(opts.Type)).Funcs(template.FuncMap{
		"yamlQuote": strconv.Quote,
	}).Parse(tmplStr)
	if err != nil {
		return err
	}

	data := TemplateData{
		Title:       opts.Title,
		Description: opts.Description,
		Level:       opts.Level,
		Weight:      weight,
		OrgID:       opts.OrgID,
		Category:    opts.Category,
		Tags:        opts.Tags,
		ID:          opts.ID,
	}

	f, err := os.Create(indexPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			utils.Log.Errorf("failed to close file %s: %v", indexPath, cerr)
		}
	}()

	if err := tmpl.Execute(f, data); err != nil {
		return err
	}

	utils.Log.Infof("Scaffolded %s '%s' at %s", opts.Type, opts.Title, indexPath)
	return nil
}

func scaffoldTree(opts ScaffoldOptions) error {
	folderName, err := makeSlug(opts.Title)
	if err != nil {
		return err
	}
	baseDir := filepath.Join(opts.TargetDir, folderName)

	err = scaffoldNode(opts, folderName)
	if err != nil {
		return err
	}

	currentDir := baseDir

	if opts.Type == string(academyModel.LearningPath) {
		courseTitle := "Course 1"
		courseOpts := opts
		courseOpts.Type = string(Course)
		courseOpts.Title = courseTitle
		courseOpts.Description = ""
		courseOpts.Level = ""
		courseOpts.Category = ""
		courseOpts.Tags = nil
		courseOpts.ID = ""
		courseOpts.TargetDir = currentDir
		err = scaffoldNode(courseOpts, "")
		if err != nil {
			return err
		}
		slug, err := makeSlug(courseTitle)
		if err != nil {
			return err
		}
		currentDir = filepath.Join(currentDir, slug)
	}

	if opts.Type == string(academyModel.LearningPath) || opts.Type == string(Course) {
		moduleTitle := "Module 1"
		moduleOpts := opts
		moduleOpts.Type = string(Module)
		moduleOpts.Title = moduleTitle
		moduleOpts.Description = ""
		moduleOpts.Level = ""
		moduleOpts.Category = ""
		moduleOpts.Tags = nil
		moduleOpts.ID = ""
		moduleOpts.TargetDir = currentDir
		err = scaffoldNode(moduleOpts, "")
		if err != nil {
			return err
		}
		slug, err := makeSlug(moduleTitle)
		if err != nil {
			return err
		}
		currentDir = filepath.Join(currentDir, slug)
	}

	if opts.Type == string(academyModel.LearningPath) || opts.Type == string(Course) || opts.Type == string(Module) {
		pageTitle := "Page 1"
		pageOpts := opts
		pageOpts.Type = string(Page)
		pageOpts.Title = pageTitle
		pageOpts.Description = ""
		pageOpts.Level = ""
		pageOpts.Category = ""
		pageOpts.Tags = nil
		pageOpts.ID = ""
		pageOpts.TargetDir = currentDir
		err = scaffoldNode(pageOpts, "")
		if err != nil {
			return err
		}
	}

	if opts.Type == string(academyModel.Certification) {
		examTitle := "Exam 1"
		examOpts := opts
		examOpts.Type = string(Exam)
		examOpts.Title = examTitle
		examOpts.Description = ""
		examOpts.Level = ""
		examOpts.Category = ""
		examOpts.Tags = nil
		examOpts.ID = ""
		examOpts.TargetDir = currentDir
		err = scaffoldNode(examOpts, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func scaffoldChallenge(opts ScaffoldOptions) error {
	folderName, err := makeSlug(opts.Title)
	if err != nil {
		return err
	}
	baseDir := filepath.Join(opts.TargetDir, folderName)

	err = scaffoldNode(opts, folderName)
	if err != nil {
		return err
	}

	currentDir := baseDir

	labOpts := opts
	labOpts.Type = string(Lab)
	labOpts.Title = "Lab"
	labOpts.Description = ""
	labOpts.Category = ""
	labOpts.Tags = nil
	labOpts.TargetDir = currentDir
	labOpts.ID = ""
	err = scaffoldNode(labOpts, "lab")
	if err != nil {
		return err
	}

	examOpts := opts
	examOpts.Type = string(Exam)
	examOpts.Title = "Exam"
	examOpts.Description = ""
	examOpts.Category = ""
	examOpts.Tags = nil
	examOpts.TargetDir = currentDir
	examOpts.ID = ""
	err = scaffoldNode(examOpts, "exam")
	if err != nil {
		return err
	}

	contentDirs := []struct {
		name string
	}{
		{"description"},
		{"getting-started"},
		{"faq"},
	}

	for _, dir := range contentDirs {
		pageOpts := opts
		pageOpts.Type = string(Page)
		pageOpts.Title = dir.name // makeSlug will handle generating the right folder name
		pageOpts.Description = ""
		pageOpts.Category = ""
		pageOpts.Tags = nil
		pageOpts.TargetDir = filepath.Join(currentDir, "content")
		pageOpts.ID = ""
		err = scaffoldNode(pageOpts, dir.name)
		if err != nil {
			return err
		}
	}

	return nil
}
