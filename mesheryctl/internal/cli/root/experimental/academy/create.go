package academy

import (
	"os"
	"path/filepath"
	"strings"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	academyModel "github.com/meshery/schemas/models/v1beta3/academy"
	"github.com/spf13/cobra"
)

type cmdAcademyCreateFlags struct {
	Type        string `json:"type" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Into        string `json:"into"`
	OrgID       string `json:"org"`
	Level       string `json:"level"`
	Category    string `json:"category"`
	Tags        string `json:"tags"`
	Force       bool   `json:"force"`
	ID          string `json:"id"`
}

var createAcademyFlags cmdAcademyCreateFlags

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Scaffold Layer5 Academy content",
	Long: `Create scaffolding for Layer5 Academy content types such as learning paths, courses, modules, and pages.
For 'learning-path', it creates a full starter tree.
For others, it adds a single node into an existing tree at the path specified by '--into'.`,
	Example: `
// Scaffold a full learning path tree
mesheryctl exp academy create --type learning-path --title "My Path" --description "Desc" --level beginner --org my-org

// Scaffold a single course into an existing tree
mesheryctl exp academy create --type course --title "New Course" --description "Desc" --into ./my-path --org my-org

// Scaffold a challenge
mesheryctl exp academy create --type challenge --title "My Challenge" --description "Desc" --org my-org
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &createAcademyFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if createAcademyFlags.Type == "" {
			return errTaxonomyType("")
		}
		cType := createAcademyFlags.Type
		if cType != string(academyModel.LearningPath) && cType != string(academyModel.Certification) && cType != string(academyModel.Challenge) {
			return errTaxonomyType(cType)
		}
		return executeCreate()
	},
}

func executeCreate() error {
	cType := createAcademyFlags.Type
	if !IsValidNodeType(cType) {
		return errTaxonomyType(createAcademyFlags.Type)
	}

	targetDir := createAcademyFlags.Into

	if targetDir == "" {
		if cType != string(academyModel.LearningPath) && cType != string(academyModel.Certification) && cType != string(academyModel.Challenge) {
			return errMissingInto()
		}
		var err error
		targetDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	orgID := createAcademyFlags.OrgID

	if cType == string(academyModel.LearningPath) || cType == string(academyModel.Certification) || cType == string(academyModel.Challenge) {
		if orgID == "" {
			return errMissingOrgID()
		}
		if err := validatePathSegment(orgID); err != nil {
			return err
		}
		targetDir = filepath.Join(targetDir, "content", contentDirSegment(cType), orgID)
	}

	var tagsList []string
	if createAcademyFlags.Tags != "" {
		for _, t := range strings.Split(createAcademyFlags.Tags, ",") {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				tagsList = append(tagsList, trimmed)
			}
		}
	}

	if createAcademyFlags.Level == "" {
		if cType == string(academyModel.LearningPath) || cType == string(academyModel.Certification) || cType == string(academyModel.Challenge) {
			createAcademyFlags.Level = string(academyModel.Beginner)
		}
	} else {
		if err := validateLevel(createAcademyFlags.Level); err != nil {
			return err
		}
	}

	opts := ScaffoldOptions{
		Type:        cType,
		Title:       createAcademyFlags.Title,
		Description: createAcademyFlags.Description,
		Level:       createAcademyFlags.Level,
		OrgID:       orgID,
		Category:    createAcademyFlags.Category,
		Tags:        tagsList,
		TargetDir:   targetDir,
		Force:       createAcademyFlags.Force,
		ID:          createAcademyFlags.ID,
	}

	// course/module routes through scaffoldTree, which auto-generates a module-1/page-1 stub tree beneath them
	if cType == string(academyModel.LearningPath) || cType == string(Course) || cType == string(Module) || cType == string(academyModel.Certification) {
		return scaffoldTree(opts)
	} else if cType == string(academyModel.Challenge) {
		return scaffoldChallenge(opts)
	}

	return scaffoldNode(opts, "")
}

func makeSubCmd(kind string) *cobra.Command {
	subCmd := &cobra.Command{
		Use:   kind + " <title>",
		Short: "Scaffold a " + kind,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			createAcademyFlags.Type = kind
			createAcademyFlags.Title = args[0]
			return mesheryctlflags.ValidateCmdFlags(cmd, &createAcademyFlags)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeCreate()
		},
	}
	return subCmd
}

func init() {
	createCmd.Flags().StringVarP(&createAcademyFlags.Type, "type", "t", "", "Content type (learning-path, certification, challenge)")
	createCmd.Flags().StringVar(&createAcademyFlags.Title, "title", "", "Title of the content")
	createCmd.Flags().StringVar(&createAcademyFlags.Description, "description", "", "Description of the content")
	createCmd.Flags().StringVar(&createAcademyFlags.Into, "into", "", "Target directory path")
	createCmd.Flags().StringVar(&createAcademyFlags.OrgID, "org", "", "Organization ID")
	createCmd.Flags().StringVar(&createAcademyFlags.Level, "level", "", "Content level (e.g., beginner, intermediate, advanced)")
	createCmd.Flags().StringVar(&createAcademyFlags.Category, "category", "", "Category of the content")
	createCmd.Flags().StringVar(&createAcademyFlags.Tags, "tags", "", "Comma-separated list of tags")
	createCmd.Flags().StringVar(&createAcademyFlags.ID, "id", "", "Content ID for Instructor Console")
	createCmd.Flags().BoolVarP(&createAcademyFlags.Force, "force", "f", false, "Overwrite existing files")

	subcommands := []string{string(Course), string(Module), string(Page), string(Lab), string(Test), string(Exam)}
	for _, kind := range subcommands {
		subCmd := makeSubCmd(kind)
		subCmd.Flags().StringVar(&createAcademyFlags.Description, "description", "", "Description of the content")
		subCmd.Flags().StringVar(&createAcademyFlags.Into, "into", "", "Target directory path")
		subCmd.Flags().StringVar(&createAcademyFlags.Level, "level", "", "Content level (e.g., beginner, intermediate, advanced)")
		subCmd.Flags().StringVar(&createAcademyFlags.Category, "category", "", "Category of the content")
		subCmd.Flags().StringVar(&createAcademyFlags.Tags, "tags", "", "Comma-separated list of tags")
		subCmd.Flags().BoolVarP(&createAcademyFlags.Force, "force", "f", false, "Overwrite existing files")
		createCmd.AddCommand(subCmd)
	}
}
