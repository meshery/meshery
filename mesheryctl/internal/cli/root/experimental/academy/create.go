package academy

import (
	"os"
	"path/filepath"
	"strings"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
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
mesheryctl exp academy create --type learning-path --title "My Path" --level beginner

// Scaffold a single course into an existing tree
mesheryctl exp academy create --type course --title "New Course" --into ./my-path
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &createAcademyFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cType := ContentType(createAcademyFlags.Type)
		if !cType.IsValid() {
			return errTaxonomyType(createAcademyFlags.Type)
		}

		targetDir := createAcademyFlags.Into

		if targetDir == "" {
			if cType != LearningPath && cType != Certification {
				return errMissingInto()
			}
			var err error
			targetDir, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		orgID := createAcademyFlags.OrgID

		if cType == LearningPath || cType == Certification {
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

		if cType == LearningPath || cType == Course || cType == Module || cType == Certification {
			return scaffoldTree(opts)
		}

		return scaffoldNode(opts)
	},
}

func init() {
	createCmd.Flags().StringVarP(&createAcademyFlags.Type, "type", "t", "", "Content type (learning-path, course, module, page, lab, test, exam, certification)")
	createCmd.Flags().StringVar(&createAcademyFlags.Title, "title", "", "Title of the content")
	createCmd.Flags().StringVar(&createAcademyFlags.Description, "description", "", "Description of the content")
	createCmd.Flags().StringVar(&createAcademyFlags.Into, "into", "", "Target directory path")
	createCmd.Flags().StringVar(&createAcademyFlags.OrgID, "org", "", "Organization ID")
	createCmd.Flags().StringVar(&createAcademyFlags.Level, "level", "", "Content level (e.g., beginner, intermediate, advanced)")
	createCmd.Flags().StringVar(&createAcademyFlags.Category, "category", "", "Category of the content")
	createCmd.Flags().StringVar(&createAcademyFlags.Tags, "tags", "", "Comma-separated list of tags")
	createCmd.Flags().StringVar(&createAcademyFlags.ID, "id", "", "Content ID for Instructor Console")
	createCmd.Flags().BoolVarP(&createAcademyFlags.Force, "force", "f", false, "Overwrite existing files")
}
