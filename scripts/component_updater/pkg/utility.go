package pkg

type SystemType int

const (
	Meshery SystemType = iota
	Docs
	RemoteProvider
)

func (dt SystemType) String() string {
	switch dt {
	case Meshery:
		return "meshery"

	case Docs:
		return "docs"

	case RemoteProvider:
		return "remote-provider"
	}
	return ""
}

func GetIndexForRegisterCol(cols []string, shouldRegister string) int {
	if shouldRegisterColIndex != -1 {
		return shouldRegisterColIndex
	}

	for index, col := range cols {
		if col == shouldRegister {
			return index
		}
	}
	return shouldRegisterColIndex
}
