package academy

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"
	academyModel "github.com/meshery/schemas/models/v1beta3/academy"
)

// Structural content types for Layer5 Academy content.
const (
	Course academyModel.ContentType = "course"
	Module academyModel.ContentType = "module"
	Page   academyModel.ContentType = "page"
	Lab    academyModel.ContentType = "lab"
	Test   academyModel.ContentType = "test"
	Exam   academyModel.ContentType = "exam"
)

// AllowedChildren maps a content type to the types of children it can contain.
var AllowedChildren = map[string][]string{
	string(academyModel.LearningPath):  {string(Course)},
	string(Course):                     {string(Module), string(Test), string(Exam)},
	string(Module):                     {string(Page), string(Test), string(Lab)},
	string(Page):                       {},
	string(academyModel.Certification): {string(Test), string(Exam)},
	string(Exam):                       {},
	string(Lab):                        {},
	string(Test):                       {},
	string(academyModel.Challenge):     {string(Lab), string(Exam)},
}

// IsValidNodeType checks if the provided string is a valid content type or node kind.
func IsValidNodeType(c string) bool {
	switch c {
	case string(academyModel.LearningPath), string(Course), string(Module), string(Page), string(academyModel.Certification), string(Lab), string(Test), string(Exam), string(academyModel.Challenge):
		return true
	}
	return false
}

// validateLevel checks if the provided level is valid according to the schema enum.
func validateLevel(level string) error {
	switch level {
	case string(academyModel.Beginner), string(academyModel.Intermediate), string(academyModel.Advanced):
		return nil
	}
	return errInvalidLevel(level)
}

func validatePathSegment(segment string) error {
	if segment == "" || segment == "." || segment == ".." {
		return fmt.Errorf("invalid path segment: %q", segment)
	}
	if strings.ContainsAny(segment, "/\\") {
		return fmt.Errorf("path segment contains directory separators: %q", segment)
	}
	return nil
}

// validateOrgID checks that the provided org ID is a valid UUID, matching the
// schema's core.Uuid type for academy content ownership.
func validateOrgID(orgID string) error {
	if _, err := uuid.FromString(orgID); err != nil {
		return errInvalidOrgID(orgID)
	}
	return nil
}

var windowsReservedNames = map[string]bool{
	"con": true, "prn": true, "aux": true, "nul": true,
	"com1": true, "com2": true, "com3": true, "com4": true, "com5": true,
	"com6": true, "com7": true, "com8": true, "com9": true,
	"lpt1": true, "lpt2": true, "lpt3": true, "lpt4": true, "lpt5": true,
	"lpt6": true, "lpt7": true, "lpt8": true, "lpt9": true,
}

var nonAlphaNumRegex = regexp.MustCompile(`[^a-z0-9]+`)

func makeSlug(input string) (string, error) {
	slug := strings.ToLower(input)
	result := nonAlphaNumRegex.ReplaceAllString(slug, "-")
	result = strings.Trim(result, "-")

	if result == "" || result == "." || result == ".." {
		return "", fmt.Errorf("invalid or unsafe path segment generated from title: %q", input)
	}
	if windowsReservedNames[result] {
		return "", fmt.Errorf("title %q produces a reserved Windows device name %q; choose a different title", input, result)
	}
	return result, nil
}
