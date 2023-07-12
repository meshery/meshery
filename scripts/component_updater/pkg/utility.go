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
