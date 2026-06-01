package performance

const (
	ServiceMesh_INVALID_MESH = "INVALID_MESH"
)

// Protocol defines the protocol type
type Protocol int32

const (
	Protocol_UNSET Protocol = 0
	Protocol_HTTP  Protocol = 1
	Protocol_GRPC  Protocol = 2
	Protocol_TCP   Protocol = 3
)

func (p Protocol) String() string {
	switch p {
	case Protocol_HTTP:
		return "http"
	case Protocol_GRPC:
		return "grpc"
	case Protocol_TCP:
		return "tcp"
	default:
		return ""
	}
}

// LoadGenerator holds load generator configuration
type LoadGenerator int32

const (
	LoadGenerator_UNSET LoadGenerator = 0
	LoadGenerator_FORTIO LoadGenerator = 1
)

// Client holds client configuration for load test
type Client struct {
	Protocol      Protocol  `json:"protocol,omitempty"`
	Connections   int32     `json:"connections,omitempty"`
	Rps           int32     `json:"rps,omitempty"`
	LoadGenerator string    `json:"loadGenerator,omitempty"`
	EndpointUrls  []string  `json:"endpointUrls,omitempty"`
}

// PerformanceTestConfig_Client is a nested type for backward compatibility
type PerformanceTestConfig_Client struct {
	Protocol      int32     `json:"protocol,omitempty"`
	Connections   int32     `json:"connections,omitempty"`
	Rps           int32     `json:"rps,omitempty"`
	LoadGenerator string    `json:"loadGenerator,omitempty"`
	EndpointUrls  []string  `json:"endpointUrls,omitempty"`
}

// Protocol constants for nested type
const (
	PerformanceTestConfig_Client_PROTOCOL_UNSET PerformanceTestConfig_Client_Protocol = 0
	PerformanceTestConfig_Client_PROTOCOL_HTTP  PerformanceTestConfig_Client_Protocol = 1
	PerformanceTestConfig_Client_PROTOCOL_GRPC  PerformanceTestConfig_Client_Protocol = 2
	PerformanceTestConfig_Client_PROTOCOL_TCP   PerformanceTestConfig_Client_Protocol = 3
)

type PerformanceTestConfig_Client_Protocol int32

// PerformanceTestConfig holds the performance test configuration
type PerformanceTestConfig struct {
	Name              string                            `json:"name,omitempty"`
	Id                string                            `json:"id,omitempty"`
	Duration          string                            `json:"duration,omitempty"`
	QueriesPerSecond  int32                             `json:"queriesPerSecond,omitempty"`
	ConcurrentRequest int32                             `json:"concurrentRequest,omitempty"`
	TestURL           string                            `json:"testUrl,omitempty"`
	RequestHeaders    map[string]string                 `json:"requestHeaders,omitempty"`
	Clients           []Client                          `json:"clients,omitempty"`
	ClientsNested     []*PerformanceTestConfig_Client   `json:"clientsNested,omitempty"`
}

// ServiceMesh holds service mesh information
type ServiceMesh struct {
	Type int32  `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

// ServiceMesh_Type_name provides string names for ServiceMesh types
var ServiceMesh_Type_name = map[int32]string{
	0: "UNSET",
	1: "ISTIO",
	2: "LINKERD",
	3: "CONSUL",
	4: "OSMI",
	5: "APPMESH",
	6: "KUMA",
	7: "TRAEFIK",
	8: "OCTARINE",
}

// ServiceMesh_Type_value provides integer values for ServiceMesh type names
var ServiceMesh_Type_value = map[string]int32{
	"UNSET":     0,
	"ISTIO":     1,
	"LINKERD":   2,
	"CONSUL":    3,
	"OSMI":      4,
	"APPMESH":   5,
	"KUMA":      6,
	"TRAEFIK":   7,
	"OCTARINE":  8,
}

