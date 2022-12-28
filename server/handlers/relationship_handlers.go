package handlers

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/models/meshmodel/core/types"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

func (h *Handler) GetAllMeshmodelRelationships(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	res := h.registryManager.GetEntities(&v1alpha1.RelationshipFilter{})
	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterMeshmodelRelationships(rw http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var cc meshmodel.MeshModelRegistrantData
	err := dec.Decode(&cc)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	switch cc.EntityType {
	case types.RelationshipDefinition:
		var r v1alpha1.RelationshipDefinition
		err = json.Unmarshal(cc.Entity, &r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.registryManager.RegisterEntity(cc.Host, r)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

// while parsing, if an error is encountered, it will return the list of relationships that have already been parsed along with the error
func parseStaticRelationships(sourceDirPath string) (rs []v1alpha1.RelationshipDefinition, err error) {
	err = filepath.Walk(sourceDirPath, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return fmt.Errorf("invalid/nil fileinfo while walking %s", path)
		}
		if !info.IsDir() {
			var rel v1alpha1.RelationshipDefinition
			byt, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			err = json.Unmarshal(byt, &rel)
			if err != nil {
				return err
			}
			rs = append(rs, rel)
		}
		return nil
	})
	return
}

func RegisterStaticMeshmodelRelationships(rm meshmodel.RegistryManager, sourceDirPath string) (err error) {
	host := meshmodel.Host{Hostname: "meshery"}
	rs, err := parseStaticRelationships(path.Clean(sourceDirPath))
	if err != nil && len(rs) == 0 {
		return
	}
	for _, r := range rs {
		err = rm.RegisterEntity(host, r)
		if err != nil {
			return
		}
	}
	return
}
