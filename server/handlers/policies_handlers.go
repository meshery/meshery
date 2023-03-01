package handlers

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/models/meshmodel/core/types"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

func (h *Handler) GetAllMeshmodelPolicy(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]
	name := mux.Vars(r)["name"]

	res := h.registryManager.GetEntities(&v1alpha1.PolicyFilter{
		Kind:      name,
		ModelName: typ,
	})
	unique := make(map[string]bool)
	var pls []v1alpha1.PolicyDefinition
	for _, p := range res {
		pl, ok := p.(v1alpha1.PolicyDefinition)
		if !ok {
			continue
		}
		if _, ok := unique[pl.Model.Name]; !ok {
			pls = append(pls, pl)
			unique[pl.Model.Name] = true
		}
	}
	if err := enc.Encode(pls); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterMeshmodelPolicy(rw http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var cc meshmodel.MeshModelRegistrantData
	err := dec.Decode(&cc)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	switch cc.EntityType {
	case types.PolicyDefinition:
		var p v1alpha1.PolicyDefinition
		err = json.Unmarshal(cc.Entity, &p)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.registryManager.RegisterEntity(cc.Host, p)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	go h.config.MeshModelSummaryChannel.Publish()
}

func parseStaticPolicy(sourceDirPath string) (rs []v1alpha1.PolicyDefinition, err error) {
	err = filepath.Walk(sourceDirPath, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return fmt.Errorf("invalid/nil fileinfo while walking %s", path)
		}
		if !info.IsDir() {
			var p v1alpha1.PolicyDefinition
			byt, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			err = json.Unmarshal(byt, &p)
			if err != nil {
				return err
			}
			rs = append(rs, p)
		}
		return nil
	})
	return
}

func RegisterStaticMeshmodelPolicy(rm meshmodel.RegistryManager, sourceDirPath string) (err error) {
	host := meshmodel.Host{Hostname: "meshery"}
	pl, err := parseStaticPolicy(path.Clean(sourceDirPath))
	if err != nil && len(pl) == 0 {
		return
	}
	for _, p := range pl {
		err = rm.RegisterEntity(host, p)
		if err != nil {
			return
		}
	}
	return
}
