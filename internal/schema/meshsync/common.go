package meshsync

import "time"

var (
	// ServiTypes
	Deployment  ServiceType = "deployment"
	StatefulSet ServiceType = "statefulset"
	DaemonSet   ServiceType = "daemonset"

	// Addons
	Prometheus Addon = "prometheus"
	Grafana    Addon = "grafana"
	Jaeger     Addon = "jaeger"
)

type ServiceType string

type Addon string

type Count float64

type TimeStamp time.Time

type ID string

type ProgressMeter struct {
	Current int64
	Total   int64
}
