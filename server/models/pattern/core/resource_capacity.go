package core

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/meshery/schemas/models/v1beta2/component"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Status constants for ResourceCapacity
const (
	StatusWithinCapacity = "within_capacity"
	StatusOverCapacity   = "over_capacity"
	StatusUnknown        = "unknown"

	EstimateLikelyToFit       = "likely_to_fit"
	EstimateReview            = "review"
	EstimateOverCapacity      = "over_capacity"
	EstimateUnknown           = "unknown"
)

// ResourceCapacityPreview summarizes resource capacity vs requirements for a cluster
type ResourceCapacityPreview struct {
	ClusterID     string                       `json:"clusterId"`
	ClusterName   string                       `json:"clusterName"`
	WorkloadCount int                          `json:"workloadCount"`
	EstimatedPods int                          `json:"estimatedPods"`
	CPU           ResourceMetricComparison     `json:"cpu"`
	Memory        ResourceMetricComparison     `json:"memory"`
	Workloads     []WorkloadSchedulingEstimate `json:"workloads"`
	Warnings      []string                     `json:"warnings"`
	Status        string                       `json:"status"` // within_capacity, over_capacity, unknown
}

// ResourceMetricComparison compares required vs available for CPU or Memory
type ResourceMetricComparison struct {
	RequiredFormatted   string `json:"requiredFormatted"`   // e.g. "7.5 cores" or "12 GiB"
	AvailableFormatted  string `json:"availableFormatted"`  // e.g. "6 cores" or "16 GiB"
	DifferenceFormatted string `json:"differenceFormatted"` // e.g. "1.5 cores above estimated availability" or "4 GiB estimated headroom"
	RequiredValue       int64  `json:"requiredValue"`       // millicores or bytes
	AvailableValue      int64  `json:"availableValue"`      // millicores or bytes
	Percentage          int    `json:"percentage"`          // e.g. 125, 75
	Status              string `json:"status"`              // within_capacity, over_capacity, unknown
	IsUnknown           bool   `json:"isUnknown"`
}

// WorkloadSchedulingEstimate is a breakdown per workload in the design
type WorkloadSchedulingEstimate struct {
	Name           string `json:"name"`
	Kind           string `json:"kind"`
	Replicas       int    `json:"replicas"`
	CPURequest     string `json:"cpuRequest"`     // "500m", "1.5 cores", or "Unspecified"
	MemoryRequest  string `json:"memoryRequest"`  // "512Mi", "2Gi", or "Unspecified"
	CPUMillicores  int64  `json:"cpuMillicores"`  // total for all replicas
	MemoryBytes    int64  `json:"memoryBytes"`    // total for all replicas
	Estimate       string `json:"estimate"`       // likely_to_fit, review, over_capacity, unknown
	HasUnspecified bool   `json:"hasUnspecified"`
}

// ClusterNodeCapacity holds allocatable and allocated capacity aggregated across nodes in a cluster
type ClusterNodeCapacity struct {
	NodeCount           int
	TotalAllocatableCPU int64 // millicores
	TotalAllocatableMem int64 // bytes
	TotalAllocatedCPU   int64 // millicores (from active pods)
	TotalAllocatedMem   int64 // bytes (from active pods)
	HasData             bool
}

// AvailableCPU returns available CPU millicores (Allocatable - Allocated, min 0)
func (c ClusterNodeCapacity) AvailableCPU() int64 {
	avail := c.TotalAllocatableCPU - c.TotalAllocatedCPU
	if avail < 0 {
		return 0
	}
	return avail
}

// AvailableMemory returns available Memory bytes (Allocatable - Allocated, min 0)
func (c ClusterNodeCapacity) AvailableMemory() int64 {
	avail := c.TotalAllocatableMem - c.TotalAllocatedMem
	if avail < 0 {
		return 0
	}
	return avail
}

// CalculateResourceCapacity computes the ResourceCapacityPreview for a given set of components against cluster capacity
func CalculateResourceCapacity(
	clusterID string,
	clusterName string,
	comps []*component.ComponentDefinition,
	clusterCapacity ClusterNodeCapacity,
) *ResourceCapacityPreview {
	preview := &ResourceCapacityPreview{
		ClusterID:   clusterID,
		ClusterName: clusterName,
		Warnings:    make([]string, 0),
		Workloads:   make([]WorkloadSchedulingEstimate, 0),
		Status:      StatusWithinCapacity,
	}

	if clusterName == "" {
		preview.ClusterName = clusterID
	}

	workloads, totalCPUReq, totalMemReq, totalPods := ExtractWorkloads(comps, clusterCapacity.NodeCount)
	preview.Workloads = workloads
	preview.WorkloadCount = len(workloads)
	preview.EstimatedPods = totalPods

	// CPU Metric Comparison
	cpuComp := ResourceMetricComparison{
		RequiredValue:     totalCPUReq,
		RequiredFormatted: FormatCPU(totalCPUReq),
	}

	// Memory Metric Comparison
	memComp := ResourceMetricComparison{
		RequiredValue:     totalMemReq,
		RequiredFormatted: FormatMemory(totalMemReq),
	}

	if !clusterCapacity.HasData {
		cpuComp.IsUnknown = true
		cpuComp.Status = StatusUnknown
		cpuComp.AvailableFormatted = "Unknown"
		cpuComp.DifferenceFormatted = "Cluster node data unavailable"

		memComp.IsUnknown = true
		memComp.Status = StatusUnknown
		memComp.AvailableFormatted = "Unknown"
		memComp.DifferenceFormatted = "Cluster node data unavailable"

		preview.Status = StatusUnknown
		preview.Warnings = append(preview.Warnings, "Cluster capacity data is unavailable from MeshSync. Scheduling estimate cannot be computed.")
	} else {
		availCPU := clusterCapacity.AvailableCPU()
		availMem := clusterCapacity.AvailableMemory()

		cpuComp.AvailableValue = availCPU
		cpuComp.AvailableFormatted = FormatCPU(availCPU)

		memComp.AvailableValue = availMem
		memComp.AvailableFormatted = FormatMemory(availMem)

		// CPU percentage & diff
		if availCPU > 0 {
			cpuComp.Percentage = int(math.Round(float64(totalCPUReq) * 100.0 / float64(availCPU)))
		} else if totalCPUReq > 0 {
			cpuComp.Percentage = 100
		} else {
			cpuComp.Percentage = 0
		}

		if totalCPUReq > availCPU {
			cpuComp.Status = StatusOverCapacity
			diff := totalCPUReq - availCPU
			cpuComp.DifferenceFormatted = fmt.Sprintf("%s above estimated availability", FormatCPU(diff))
			preview.Status = StatusOverCapacity
			preview.Warnings = append(preview.Warnings, fmt.Sprintf("CPU requests exceed available cluster capacity by %s.", FormatCPU(diff)))
		} else {
			cpuComp.Status = StatusWithinCapacity
			headroom := availCPU - totalCPUReq
			cpuComp.DifferenceFormatted = fmt.Sprintf("%s estimated headroom", FormatCPU(headroom))
		}

		// Memory percentage & diff
		if availMem > 0 {
			memComp.Percentage = int(math.Round(float64(totalMemReq) * 100.0 / float64(availMem)))
		} else if totalMemReq > 0 {
			memComp.Percentage = 100
		} else {
			memComp.Percentage = 0
		}

		if totalMemReq > availMem {
			memComp.Status = StatusOverCapacity
			diff := totalMemReq - availMem
			memComp.DifferenceFormatted = fmt.Sprintf("%s above estimated availability", FormatMemory(diff))
			preview.Status = StatusOverCapacity
			preview.Warnings = append(preview.Warnings, fmt.Sprintf("Memory requests exceed available cluster capacity by %s.", FormatMemory(diff)))
		} else {
			memComp.Status = StatusWithinCapacity
			headroom := availMem - totalMemReq
			memComp.DifferenceFormatted = fmt.Sprintf("%s estimated headroom", FormatMemory(headroom))
		}
	}

	// Refine workload estimates
	for i := range preview.Workloads {
		w := &preview.Workloads[i]
		if w.HasUnspecified {
			preview.Warnings = append(preview.Warnings, fmt.Sprintf("Workload '%s' has missing CPU or memory resource requests.", w.Name))
		}
		if !clusterCapacity.HasData {
			w.Estimate = EstimateUnknown
		} else if w.HasUnspecified {
			w.Estimate = EstimateReview
		} else if w.CPUMillicores > clusterCapacity.AvailableCPU() || w.MemoryBytes > clusterCapacity.AvailableMemory() {
			w.Estimate = EstimateReview
		} else {
			w.Estimate = EstimateLikelyToFit
		}
	}

	preview.CPU = cpuComp
	preview.Memory = memComp
	return preview
}

// ExtractWorkloads extracts workload components, calculates requests and replica counts
func ExtractWorkloads(comps []*component.ComponentDefinition, nodeCount int) ([]WorkloadSchedulingEstimate, int64, int64, int) {
	workloads := make([]WorkloadSchedulingEstimate, 0)
	var totalCPUMillicores int64
	var totalMemoryBytes int64
	var totalPods int

	for _, cmp := range comps {
		if cmp == nil || cmp.Component.Kind == "" {
			continue
		}
		kind := cmp.Component.Kind
		if !isWorkloadKind(kind) {
			continue
		}

		name := cmp.DisplayName
		if name == "" {
			name = cmp.Component.Kind
		}

		config := cmp.Configuration
		if config == nil {
			config = make(map[string]interface{})
		}

		replicas := extractReplicas(kind, config, nodeCount)
		totalPods += replicas

		podCPU, podMem, hasUnspecified := extractContainerRequests(kind, config)

		totalWorkloadCPU := podCPU * int64(replicas)
		totalWorkloadMem := podMem * int64(replicas)

		totalCPUMillicores += totalWorkloadCPU
		totalMemoryBytes += totalWorkloadMem

		cpuReqFormatted := "Unspecified"
		if podCPU > 0 || !hasUnspecified {
			cpuReqFormatted = FormatCPU(podCPU)
		}

		memReqFormatted := "Unspecified"
		if podMem > 0 || !hasUnspecified {
			memReqFormatted = FormatMemory(podMem)
		}

		workload := WorkloadSchedulingEstimate{
			Name:           name,
			Kind:           kind,
			Replicas:       replicas,
			CPURequest:     cpuReqFormatted,
			MemoryRequest:  memReqFormatted,
			CPUMillicores:  totalWorkloadCPU,
			MemoryBytes:    totalWorkloadMem,
			HasUnspecified: hasUnspecified,
			Estimate:       EstimateLikelyToFit,
		}
		workloads = append(workloads, workload)
	}

	return workloads, totalCPUMillicores, totalMemoryBytes, totalPods
}

func isWorkloadKind(kind string) bool {
	switch strings.ToLower(kind) {
	case "deployment", "statefulset", "daemonset", "pod", "job", "cronjob", "replicaset":
		return true
	default:
		return false
	}
}

func extractReplicas(kind string, config map[string]interface{}, nodeCount int) int {
	k := strings.ToLower(kind)
	if k == "pod" {
		return 1
	}
	if k == "cronjob" {
		if spec, ok := getMap(config, "spec"); ok {
			if susp, exists := spec["suspend"]; exists {
				if b, ok := susp.(bool); ok && b {
					return 0
				}
			}
			if jt, ok := getMap(spec, "jobTemplate"); ok {
				if jtSpec, ok := getMap(jt, "spec"); ok {
					if susp, exists := jtSpec["suspend"]; exists {
						if b, ok := susp.(bool); ok && b {
							return 0
						}
					}
					if pVal, exists := jtSpec["parallelism"]; exists && pVal != nil {
						if p := toInt(pVal); p >= 0 {
							return p
						}
					}
				}
			}
		}
		return 1
	}
	if k == "daemonset" {
		if nodeCount > 0 {
			return nodeCount
		}
		return 1
	}

	// For Deployment, StatefulSet, ReplicaSet, Job
	spec, ok := getMap(config, "spec")
	if !ok {
		return 1
	}

	if k == "job" {
		if susp, exists := spec["suspend"]; exists {
			if b, ok := susp.(bool); ok && b {
				return 0
			}
		}
		if pVal, exists := spec["parallelism"]; exists && pVal != nil {
			if p := toInt(pVal); p >= 0 {
				return p
			}
		}
		return 1
	}

	if rVal, exists := spec["replicas"]; exists && rVal != nil {
		if r := toInt(rVal); r >= 0 {
			return r
		}
	}
	return 1
}

func extractContainerRequests(kind string, config map[string]interface{}) (int64, int64, bool) {
	podSpec, containers, initContainers := findPodSpec(kind, config)
	if len(containers) == 0 && len(initContainers) == 0 {
		hasPodLevel := false
		if podSpec != nil {
			if resMap, ok := getMap(podSpec, "resources"); ok {
				if reqMap, ok := getMap(resMap, "requests"); ok && len(reqMap) > 0 {
					hasPodLevel = true
				}
			}
		}
		if !hasPodLevel {
			return 0, 0, true
		}
	}

	var containersCPU int64
	var containersMem int64
	missingCPU := false
	missingMem := false

	for _, container := range containers {
		cMap, ok := container.(map[string]interface{})
		if !ok {
			continue
		}

		resMap, ok := getMap(cMap, "resources")
		if !ok {
			missingCPU = true
			missingMem = true
			continue
		}

		reqMap, ok := getMap(resMap, "requests")
		if !ok {
			missingCPU = true
			missingMem = true
			continue
		}

		cpuVal := reqMap["cpu"]
		memVal := reqMap["memory"]

		cpuMilli, cpuSpecified := ParseCPU(cpuVal)
		memBytes, memSpecified := ParseMemory(memVal)

		if !cpuSpecified {
			missingCPU = true
		}
		if !memSpecified {
			missingMem = true
		}

		containersCPU += cpuMilli
		containersMem += memBytes
	}

	var maxInitCPU int64
	var maxInitMem int64
	for _, ic := range initContainers {
		icMap, ok := ic.(map[string]interface{})
		if !ok {
			continue
		}

		resMap, ok := getMap(icMap, "resources")
		if !ok {
			continue
		}

		reqMap, ok := getMap(resMap, "requests")
		if !ok {
			continue
		}

		if cpuVal, ok := reqMap["cpu"]; ok {
			if cpuMilli, specified := ParseCPU(cpuVal); specified && cpuMilli > maxInitCPU {
				maxInitCPU = cpuMilli
			}
		}

		if memVal, ok := reqMap["memory"]; ok {
			if memBytes, specified := ParseMemory(memVal); specified && memBytes > maxInitMem {
				maxInitMem = memBytes
			}
		}
	}

	// In Kubernetes scheduling, effective pod request is max(sum(containers), max(initContainers))
	effectiveCPU := containersCPU
	if maxInitCPU > effectiveCPU {
		effectiveCPU = maxInitCPU
	}
	effectiveMem := containersMem
	if maxInitMem > effectiveMem {
		effectiveMem = maxInitMem
	}

	// Check if explicit pod-level requests are set in spec.resources.requests (Kubernetes v1.35)
	if podSpec != nil {
		if resMap, ok := getMap(podSpec, "resources"); ok {
			if reqMap, ok := getMap(resMap, "requests"); ok {
				if cpuVal, ok := reqMap["cpu"]; ok {
					if cpuMilli, specified := ParseCPU(cpuVal); specified {
						effectiveCPU = cpuMilli
						missingCPU = false
					}
				}
				if memVal, ok := reqMap["memory"]; ok {
					if memBytes, specified := ParseMemory(memVal); specified {
						effectiveMem = memBytes
						missingMem = false
					}
				}
			}
		}
	}

	hasUnspecified := missingCPU || missingMem
	return effectiveCPU, effectiveMem, hasUnspecified
}

func findPodSpec(kind string, config map[string]interface{}) (map[string]interface{}, []interface{}, []interface{}) {
	spec, _ := getMap(config, "spec")
	if spec == nil {
		spec = config
	}

	k := strings.ToLower(kind)
	if k == "pod" {
		cList, _ := spec["containers"].([]interface{})
		icList, _ := spec["initContainers"].([]interface{})
		return spec, cList, icList
	}

	if k == "cronjob" {
		if jt, ok := getMap(spec, "jobTemplate"); ok {
			if jtSpec, ok := getMap(jt, "spec"); ok {
				if tpl, ok := getMap(jtSpec, "template"); ok {
					if tplSpec, ok := getMap(tpl, "spec"); ok {
						cList, _ := tplSpec["containers"].([]interface{})
						icList, _ := tplSpec["initContainers"].([]interface{})
						return tplSpec, cList, icList
					}
				}
			}
		}
	}

	// Deployment, StatefulSet, DaemonSet, Job, ReplicaSet
	if tpl, ok := getMap(spec, "template"); ok {
		if tplSpec, ok := getMap(tpl, "spec"); ok {
			cList, _ := tplSpec["containers"].([]interface{})
			icList, _ := tplSpec["initContainers"].([]interface{})
			return tplSpec, cList, icList
		}
	}

	// Fallback check
	cList, _ := spec["containers"].([]interface{})
	icList, _ := spec["initContainers"].([]interface{})
	return spec, cList, icList
}

// ParseCPU parses Kubernetes CPU requests into millicores
func ParseCPU(val interface{}) (int64, bool) {
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case string:
		str := strings.TrimSpace(v)
		if str == "" {
			return 0, false
		}
		// Strip "core" / "cores" suffix so human-formatted values like "2 cores" or "1 core" round-trip
		lower := strings.ToLower(str)
		if strings.HasSuffix(lower, "cores") {
			str = strings.TrimSpace(str[:len(str)-5])
		} else if strings.HasSuffix(lower, "core") {
			str = strings.TrimSpace(str[:len(str)-4])
		}
		q, err := resource.ParseQuantity(str)
		if err == nil {
			return q.MilliValue(), true
		}
		// Fallback manual parse
		if strings.HasSuffix(str, "m") {
			numStr := strings.TrimSuffix(str, "m")
			if num, err := strconv.ParseInt(numStr, 10, 64); err == nil {
				return num, true
			}
		}
		if f, err := strconv.ParseFloat(str, 64); err == nil {
			return int64(math.Round(f * 1000)), true
		}
		return 0, false
	case float64:
		return int64(math.Round(v * 1000)), true
	case int:
		return int64(v * 1000), true
	case int64:
		return v * 1000, true
	default:
		return 0, false
	}
}

// ParseMemory parses Kubernetes memory requests into bytes
func ParseMemory(val interface{}) (int64, bool) {
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case string:
		str := strings.TrimSpace(v)
		if str == "" {
			return 0, false
		}
		// Normalize "GiB" -> "Gi", "MiB" -> "Mi", "1000B" -> "1000", etc. so human-formatted values round-trip
		lower := strings.ToLower(str)
		if strings.HasSuffix(lower, "b") {
			trimmed := strings.TrimSpace(str[:len(str)-1])
			if trimmed != "" {
				if strings.HasSuffix(strings.ToLower(trimmed), "k") && !strings.HasSuffix(strings.ToLower(trimmed), "i") {
					trimmed = trimmed[:len(trimmed)-1] + "k"
				}
				str = trimmed
			}
		}
		q, err := resource.ParseQuantity(str)
		if err == nil {
			return q.Value(), true
		}
		// Fallback manual parse
		return parseMemoryManual(str)
	case float64:
		return int64(v), true
	case int:
		return int64(v), true
	case int64:
		return v, true
	default:
		return 0, false
	}
}

func parseMemoryManual(str string) (int64, bool) {
	suffixes := map[string]int64{
		"Ki": 1024,
		"Mi": 1024 * 1024,
		"Gi": 1024 * 1024 * 1024,
		"Ti": 1024 * 1024 * 1024 * 1024,
		"k":  1000,
		"M":  1000 * 1000,
		"G":  1000 * 1000 * 1000,
		"T":  1000 * 1000 * 1000 * 1000,
	}

	for suffix, mult := range suffixes {
		if strings.HasSuffix(str, suffix) {
			numStr := strings.TrimSpace(strings.TrimSuffix(str, suffix))
			if f, err := strconv.ParseFloat(numStr, 64); err == nil {
				return int64(f * float64(mult)), true
			}
		}
	}

	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return int64(f), true
	}
	return 0, false
}

// FormatCPU formats millicores into human-readable cores/millicores string
func FormatCPU(milli int64) string {
	if milli <= 0 {
		return "0 cores"
	}
	if milli < 1000 {
		return fmt.Sprintf("%dm", milli)
	}
	if milli%1000 == 0 {
		cores := milli / 1000
		if cores == 1 {
			return "1 core"
		}
		return fmt.Sprintf("%d cores", cores)
	}
	coresFloat := float64(milli) / 1000.0
	// If one decimal place is enough (e.g. 7.5 cores)
	formatted := fmt.Sprintf("%.1f", coresFloat)
	formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
	return fmt.Sprintf("%s cores", formatted)
}

// FormatMemory formats bytes into human-readable memory string (GiB, MiB, etc.)
func FormatMemory(bytes int64) string {
	if bytes <= 0 {
		return "0 B"
	}
	const (
		KiB = int64(1024)
		MiB = 1024 * KiB
		GiB = 1024 * MiB
		TiB = 1024 * GiB
	)

	if bytes >= GiB {
		if bytes%GiB == 0 {
			return fmt.Sprintf("%d GiB", bytes/GiB)
		}
		gibFloat := float64(bytes) / float64(GiB)
		formatted := fmt.Sprintf("%.1f", gibFloat)
		formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
		return fmt.Sprintf("%s GiB", formatted)
	}

	if bytes >= MiB {
		if bytes%MiB == 0 {
			return fmt.Sprintf("%d MiB", bytes/MiB)
		}
		mibFloat := float64(bytes) / float64(MiB)
		formatted := fmt.Sprintf("%.1f", mibFloat)
		formatted = strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
		return fmt.Sprintf("%s MiB", formatted)
	}

	if bytes >= KiB {
		return fmt.Sprintf("%d KiB", bytes/KiB)
	}

	return fmt.Sprintf("%d B", bytes)
}

func getMap(m map[string]interface{}, key string) (map[string]interface{}, bool) {
	if m == nil {
		return nil, false
	}
	val, ok := m[key]
	if !ok || val == nil {
		return nil, false
	}
	if casted, ok := val.(map[string]interface{}); ok {
		return casted, true
	}
	return nil, false
}

func toInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return int(i)
		}
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}
