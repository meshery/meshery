package academy

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
}

// IsValid checks if the provided string is a valid ContentType.
func (c ContentType) IsValid() bool {
	switch c {
	case LearningPath, Course, Module, Page, Certification, Lab, Test, Exam:
		return true
	}
	return false
}
