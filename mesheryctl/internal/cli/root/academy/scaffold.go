package academy

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/academy/templates"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	academyModel "github.com/meshery/schemas/models/v1beta3/academy"
	"gopkg.in/yaml.v3"
)

// TemplateData carries the data rendered into the frontmatter template.
// The local ID field intentionally shadows ChildNode.ID, which is a typed uuid,
// because the scaffolded frontmatter exposes a replaceable placeholder (e.g."REPLACE_WITH_INSTRUCTOR_CONSOLE_ID")
// rather than a uuid.
type TemplateData struct {
	academyModel.ChildNode
	Level    academyModel.Level
	OrgID    string
	Category string
	Tags     []string
	ID       string
}

// TypeString returns the content type for the frontmatter template.
func (t TemplateData) TypeString() string {
	if t.Type != nil {
		return string(*t.Type)
	}
	return ""
}

// LevelString returns the level for the frontmatter template.
func (t TemplateData) LevelString() string {
	return string(t.Level)
}

// WeightInt returns the weight for the frontmatter template.
func (t TemplateData) WeightInt() int {
	if t.Weight != nil {
		return int(*t.Weight)
	}
	return 0
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
	if IsValidNodeType(cType) {
		return templates.NodeTemplate
	}
	return ""
}

func contentDirSegment(cType string) string {
	return cType + "s"
}

// ParentFrontmatter is the slice of a parent _index.md frontmatter needed for
// nesting validation and metadata inheritance. Categories and tags have no schema
// equivalent and stay as plain values.
type ParentFrontmatter struct {
	Type     academyModel.ContentType `yaml:"type"`
	Level    academyModel.Level       `yaml:"level"`
	Category string                   `yaml:"categories"`
	Tags     []string                 `yaml:"tags"`
}

func isRootType(cType string) bool {
	return cType == string(academyModel.LearningPath) || cType == string(academyModel.Certification) || cType == string(academyModel.Challenge)
}

func checkNesting(cType string, parentDir string) (ParentFrontmatter, error) {
	var pf ParentFrontmatter
	indexPath := filepath.Join(parentDir, "_index.md")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		if isRootType(cType) {
			return pf, nil
		}
		return pf, errInvalidParentMetadata(indexPath, cType, "no parent _index.md found")
	}

	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		if isRootType(cType) {
			return pf, nil
		}
		return pf, errInvalidParentMetadata(indexPath, cType, "malformed frontmatter")
	}

	if err := yaml.Unmarshal(parts[1], &pf); err != nil {
		if isRootType(cType) {
			return pf, nil
		}
		return pf, errInvalidParentMetadata(indexPath, cType, fmt.Sprintf("invalid YAML: %v", err))
	}

	if pf.Type == "" {
		if isRootType(cType) {
			return pf, nil
		}
		return pf, errInvalidParentMetadata(indexPath, cType, "has no type field")
	}

	allowed, exists := AllowedChildren[string(pf.Type)]
	if !exists {
		return pf, errInvalidNesting(string(pf.Type), cType)
	}
	for _, child := range allowed {
		if child == cType {
			return pf, nil
		}
	}
	return pf, errInvalidNesting(string(pf.Type), cType)
}

type ScaffoldOptions struct {
	Type        academyModel.ContentType
	Title       string
	Description string
	Level       academyModel.Level
	OrgID       string
	Category    string
	Tags        []string
	TargetDir   string
	Force       bool
	ID          string
	SkipNesting bool
}

func scaffoldNode(opts ScaffoldOptions, explicitFolderName string) error {
	tmplStr := getTemplateString(string(opts.Type))
	if tmplStr == "" {
		return errTaxonomyType(string(opts.Type))
	}

	var pf ParentFrontmatter
	if !opts.SkipNesting {
		var err error
		pf, err = checkNesting(string(opts.Type), opts.TargetDir)
		if err != nil {
			return err
		}
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

	if opts.ID == "" && isRootType(string(opts.Type)) {
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

	if opts.Type == Test && (parentType == Course || parentType == Module) {
		indexPath = filepath.Join(opts.TargetDir, "test.md")
	} else if opts.Type == Exam && parentType == Course {
		indexPath = filepath.Join(opts.TargetDir, "course-exam.md")
	} else {
		if opts.Type == Test && parentType == academyModel.Certification {
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

	weightF := float32(weight)
	data := TemplateData{
		ChildNode: academyModel.ChildNode{
			Title:       opts.Title,
			Description: opts.Description,
			Type:        &opts.Type,
			Weight:      &weightF,
		},
		Level:    opts.Level,
		OrgID:    opts.OrgID,
		Category: opts.Category,
		Tags:     opts.Tags,
		ID:       opts.ID,
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

func scaffoldChild(opts ScaffoldOptions, cType academyModel.ContentType, title, into string) (string, error) {
	child := opts
	child.Type = cType
	child.Title = title
	child.Description = ""
	child.Level = ""
	child.Category = ""
	child.Tags = nil
	child.ID = ""
	child.TargetDir = into

	if err := scaffoldNode(child, ""); err != nil {
		return "", err
	}

	slug, err := makeSlug(title)
	if err != nil {
		return "", err
	}

	return filepath.Join(into, slug), nil
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

	// Only root types reach scaffoldTree: learning-path builds a course/module/page
	// starter tree; certification builds an exam.
	if opts.Type == academyModel.LearningPath {
		currentDir, err = scaffoldChild(opts, Course, "Course 1", currentDir)
		if err != nil {
			return err
		}
		currentDir, err = scaffoldChild(opts, Module, "Module 1", currentDir)
		if err != nil {
			return err
		}
		_, err = scaffoldChild(opts, Page, "Page 1", currentDir)
		if err != nil {
			return err
		}
	}

	if opts.Type == academyModel.Certification {
		_, err = scaffoldChild(opts, Exam, "Exam 1", currentDir)
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

	_, err = scaffoldChild(opts, Lab, "Lab", currentDir)
	if err != nil {
		return err
	}

	_, err = scaffoldChild(opts, Exam, "Exam", currentDir)
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
		pageOpts.Type = Page
		pageOpts.Title = dir.name
		pageOpts.Description = ""
		pageOpts.Category = ""
		pageOpts.Tags = nil
		pageOpts.TargetDir = filepath.Join(currentDir, "content")
		pageOpts.ID = ""
		pageOpts.SkipNesting = true
		err = scaffoldNode(pageOpts, dir.name)
		if err != nil {
			return err
		}
	}

	return nil
}
