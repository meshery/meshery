package academy

import (
	"fmt"
	"strings"
)

// ContentType represents the types of Layer5 Academy content.
type ContentType string

const (
	LearningPath  ContentType = "learning-path"
	Course        ContentType = "course"
	Module        ContentType = "module"
	Page          ContentType = "page"
	Certification ContentType = "certification"
	Lab           ContentType = "lab"
	Test          ContentType = "test"
	Exam          ContentType = "exam"
)

// AllowedChildren maps a content type to the types of children it can contain.
var AllowedChildren = map[ContentType][]ContentType{
	LearningPath:  {Course},
	Course:        {Module, Test, Exam},
	Module:        {Page, Test, Lab},
	Page:          {},
	Certification: {Test, Exam},
	Exam:          {},
	Lab:           {},
	Test:          {},
}

// IsValid checks if the provided string is a valid ContentType.
func (c ContentType) IsValid() bool {
	switch c {
	case LearningPath, Course, Module, Page, Certification, Lab, Test, Exam:
		return true
	}
	return false
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

var windowsReservedNames = map[string]bool{
	"con": true, "prn": true, "aux": true, "nul": true,
	"com1": true, "com2": true, "com3": true, "com4": true, "com5": true,
	"com6": true, "com7": true, "com8": true, "com9": true,
	"lpt1": true, "lpt2": true, "lpt3": true, "lpt4": true, "lpt5": true,
	"lpt6": true, "lpt7": true, "lpt8": true, "lpt9": true,
}

func makeSlug(input string) (string, error) {
	// First convert to lower case and replace spaces with hyphens
	slug := strings.ToLower(strings.ReplaceAll(input, " ", "-"))

	// Replace non-allowed characters with -
	var b strings.Builder
	for _, ch := range slug {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' {
			b.WriteRune(ch)
		} else {
			b.WriteRune('-')
		}
	}

	// Collapse repeated hyphens
	var collapsed strings.Builder
	lastDash := false
	for _, ch := range b.String() {
		if ch == '-' {
			if !lastDash {
				collapsed.WriteRune(ch)
				lastDash = true
			}
		} else {
			collapsed.WriteRune(ch)
			lastDash = false
		}
	}

	// Trim leading and trailing hyphens
	result := strings.Trim(collapsed.String(), "-")

	if result == "" || result == "." || result == ".." {
		return "", fmt.Errorf("invalid or unsafe path segment generated from title: %q", input)
	}
	if windowsReservedNames[result] {
		return "", fmt.Errorf("title %q produces a reserved Windows device name %q; choose a different title", input, result)
	}
	return result, nil
}
