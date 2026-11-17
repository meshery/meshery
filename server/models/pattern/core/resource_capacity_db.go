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

// PodSpecAttribute holds the container resource requests from MeshSync Pod.Spec.Attribute
type PodSpecAttribute struct {
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
	capResult.HasData = true

	for _, node := range nodes {
		if node.Status == nil || node.Status.Attribute == "" {
			continue
		}

		var nodeStatus NodeStatusAttribute
		if err := json.Unmarshal([]byte(node.Status.Attribute), &nodeStatus); err != nil {
			continue
		}

		// Prefer Allocatable over Capacity
		resourceMap := nodeStatus.Allocatable
		if len(resourceMap) == 0 {
			resourceMap = nodeStatus.Capacity
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

	// 2. Query active Pods for cluster to determine claimed resources
	var pods []meshsyncmodel.KubernetesResource
	err = db.Preload("Spec").Preload("Status").
		Where("cluster_id IN (?) AND kind = ?", clusterIDs, "Pod").
		Find(&pods).Error

	if err == nil {
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

			for _, c := range podSpec.Containers {
				if cpuVal, ok := c.Resources.Requests["cpu"]; ok {
					if cpuMilli, specified := ParseCPU(cpuVal); specified {
						capResult.TotalAllocatedCPU += cpuMilli
					}
				}
				if memVal, ok := c.Resources.Requests["memory"]; ok {
					if memBytes, specified := ParseMemory(memVal); specified {
						capResult.TotalAllocatedMem += memBytes
					}
				}
			}
		}
	}

	return capResult
}
