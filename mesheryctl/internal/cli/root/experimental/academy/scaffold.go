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

type ScaffoldOptions struct {
	Type        ContentType
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

func scaffoldNode(opts ScaffoldOptions) error {
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

	if opts.ID == "" && (opts.Type == LearningPath || opts.Type == Certification) {
		opts.ID = "REPLACE_WITH_INSTRUCTOR_CONSOLE_ID"
	}

	var indexPath string
	folderName, err := makeSlug(opts.Title)
	if err != nil {
		return err
	}

	if opts.Type == Test && (parentType == Course || parentType == Module) {
		indexPath = filepath.Join(opts.TargetDir, "test.md")
	} else if opts.Type == Exam && parentType == Course {
		indexPath = filepath.Join(opts.TargetDir, "course-exam.md")
	} else {
		if opts.Type == Test && parentType == Certification {
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

	err = scaffoldNode(opts)
	if err != nil {
		return err
	}

	currentDir := baseDir

	if opts.Type == LearningPath {
		courseTitle := "Course 1"
		courseOpts := opts
		courseOpts.Type = Course
		courseOpts.Title = courseTitle
		courseOpts.Description = ""
		courseOpts.Level = ""
		courseOpts.Category = ""
		courseOpts.Tags = nil
		courseOpts.ID = ""
		courseOpts.TargetDir = currentDir
		err = scaffoldNode(courseOpts)
		if err != nil {
			return err
		}
		slug, err := makeSlug(courseTitle)
		if err != nil {
			return err
		}
		currentDir = filepath.Join(currentDir, slug)
	}

	if opts.Type == LearningPath || opts.Type == Course {
		moduleTitle := "Module 1"
		moduleOpts := opts
		moduleOpts.Type = Module
		moduleOpts.Title = moduleTitle
		moduleOpts.Description = ""
		moduleOpts.Level = ""
		moduleOpts.Category = ""
		moduleOpts.Tags = nil
		moduleOpts.ID = ""
		moduleOpts.TargetDir = currentDir
		err = scaffoldNode(moduleOpts)
		if err != nil {
			return err
		}
		slug, err := makeSlug(moduleTitle)
		if err != nil {
			return err
		}
		currentDir = filepath.Join(currentDir, slug)
	}

	if opts.Type == LearningPath || opts.Type == Course || opts.Type == Module {
		pageTitle := "Page 1"
		pageOpts := opts
		pageOpts.Type = Page
		pageOpts.Title = pageTitle
		pageOpts.Description = ""
		pageOpts.Level = ""
		pageOpts.Category = ""
		pageOpts.Tags = nil
		pageOpts.ID = ""
		pageOpts.TargetDir = currentDir
		err = scaffoldNode(pageOpts)
		if err != nil {
			return err
		}
	}

	if opts.Type == Certification {
		examTitle := "Exam 1"
		examOpts := opts
		examOpts.Type = Exam
		examOpts.Title = examTitle
		examOpts.Description = ""
		examOpts.Level = ""
		examOpts.Category = ""
		examOpts.Tags = nil
		examOpts.ID = ""
		examOpts.TargetDir = currentDir
		err = scaffoldNode(examOpts)
		if err != nil {
			return err
		}
	}

	return nil
}
