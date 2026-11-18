package core

import (
	"encoding/json"

	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"gorm.io/gorm"
)

// NodeStatusAttribute holds the relevant fields from MeshSync Node.Status.Attribute
type NodeStatusAttribute struct {
	Allocatable map[string]string `json:"allocatable"`
	Capacity    map[string]string `json:"capacity"`
}

// PodStatusAttribute holds the relevant fields from MeshSync Pod.Status.Attribute
type PodStatusAttribute struct {
	Phase string `json:"phase"`
}

// PodSpecAttribute holds the container and pod-level resource requests from MeshSync Pod.Spec.Attribute
type PodSpecAttribute struct {
	NodeName   string `json:"nodeName,omitempty"`
	Containers []struct {
		Resources struct {
			Requests map[string]interface{} `json:"requests"`
		} `json:"resources"`
	} `json:"containers"`
	InitContainers []struct {
		Resources struct {
			Requests map[string]interface{} `json:"requests"`
		} `json:"resources"`
	} `json:"initContainers"`
	Resources *struct {
		Requests map[string]interface{} `json:"requests"`
	} `json:"resources,omitempty"`
}

// GetClusterNodeCapacityFromPersister queries MeshSync database records to determine cluster allocatable and allocated capacity
func GetClusterNodeCapacityFromPersister(db *gorm.DB, clusterIDs []string) ClusterNodeCapacity {
	capResult := ClusterNodeCapacity{
		HasData: false,
	}

	if db == nil || len(clusterIDs) == 0 {
		return capResult
	}

	// 1. Query Nodes for cluster
	var nodes []meshsyncmodel.KubernetesResource
	err := db.Preload("Status").
		Where("cluster_id IN (?) AND kind = ?", clusterIDs, "Node").
		Find(&nodes).Error

	if err != nil || len(nodes) == 0 {
		return capResult
	}

	capResult.NodeCount = len(nodes)

	for _, node := range nodes {
		if node.Status == nil || node.Status.Attribute == "" {
			continue
		}

		var nodeStatus NodeStatusAttribute
		if err := json.Unmarshal([]byte(node.Status.Attribute), &nodeStatus); err != nil {
			continue
		}

		// In Kubernetes scheduling, only Allocatable resources can be scheduled
		resourceMap := nodeStatus.Allocatable
		if len(resourceMap) == 0 {
			continue
		}

		if cpuStr, ok := resourceMap["cpu"]; ok {
			if cpuMilli, specified := ParseCPU(cpuStr); specified {
				capResult.TotalAllocatableCPU += cpuMilli
			}
		}

		if memStr, ok := resourceMap["memory"]; ok {
			if memBytes, specified := ParseMemory(memStr); specified {
				capResult.TotalAllocatableMem += memBytes
			}
		}
	}

	// If node resource data is missing or zero, cluster capacity is unknown
	if capResult.TotalAllocatableCPU == 0 && capResult.TotalAllocatableMem == 0 {
		return ClusterNodeCapacity{HasData: false}
	}

	// 2. Query active Pods for cluster to determine claimed resources
	var pods []meshsyncmodel.KubernetesResource
	err = db.Preload("Spec").Preload("Status").
		Where("cluster_id IN (?) AND kind = ?", clusterIDs, "Pod").
		Find(&pods).Error

	if err != nil {
		// Return unknown if pod query fails so we never under-report allocated claims
		return ClusterNodeCapacity{HasData: false}
	}

	for _, pod := range pods {
		if pod.Status != nil && pod.Status.Attribute != "" {
			var podStatus PodStatusAttribute
			if err := json.Unmarshal([]byte(pod.Status.Attribute), &podStatus); err == nil {
				// Succeeded and Failed pods no longer consume cluster resources
				if podStatus.Phase == "Succeeded" || podStatus.Phase == "Failed" {
					continue
				}
			}
		}

		if pod.Spec == nil || pod.Spec.Attribute == "" {
			continue
		}

		var podSpec PodSpecAttribute
		if err := json.Unmarshal([]byte(pod.Spec.Attribute), &podSpec); err != nil {
			continue
		}

		// Only pods assigned to a node consume node allocatable capacity
		if podSpec.NodeName == "" {
			continue
		}

		var containersCPU int64
		var containersMem int64
		for _, c := range podSpec.Containers {
			if cpuVal, ok := c.Resources.Requests["cpu"]; ok {
				if cpuMilli, specified := ParseCPU(cpuVal); specified {
					containersCPU += cpuMilli
				}
			}
			if memVal, ok := c.Resources.Requests["memory"]; ok {
				if memBytes, specified := ParseMemory(memVal); specified {
					containersMem += memBytes
				}
			}
		}

		// In Kubernetes scheduling, init containers run sequentially before app containers;
		// effective pod request is max(sum(containers), max(initContainers))
		var maxInitCPU int64
		var maxInitMem int64
		for _, ic := range podSpec.InitContainers {
			if cpuVal, ok := ic.Resources.Requests["cpu"]; ok {
				if cpuMilli, specified := ParseCPU(cpuVal); specified && cpuMilli > maxInitCPU {
					maxInitCPU = cpuMilli
				}
			}
			if memVal, ok := ic.Resources.Requests["memory"]; ok {
				if memBytes, specified := ParseMemory(memVal); specified && memBytes > maxInitMem {
					maxInitMem = memBytes
				}
			}
		}

		effectiveCPU := containersCPU
		if maxInitCPU > effectiveCPU {
			effectiveCPU = maxInitCPU
		}
		effectiveMem := containersMem
		if maxInitMem > effectiveMem {
			effectiveMem = maxInitMem
		}

		// Check if explicit pod-level requests are set (Kubernetes v1.35)
		if podSpec.Resources != nil && podSpec.Resources.Requests != nil {
			if cpuVal, ok := podSpec.Resources.Requests["cpu"]; ok {
				if cpuMilli, specified := ParseCPU(cpuVal); specified {
					effectiveCPU = cpuMilli
				}
			}
			if memVal, ok := podSpec.Resources.Requests["memory"]; ok {
				if memBytes, specified := ParseMemory(memVal); specified {
					effectiveMem = memBytes
				}
			}
		}

		capResult.TotalAllocatedCPU += effectiveCPU
		capResult.TotalAllocatedMem += effectiveMem
	}

	capResult.HasData = true
	return capResult
}
