package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/artifacthub"
	"github.com/layer5io/meshkit/utils/manifests"
)

type ExternalSystem struct {
	SystemName string `json:"systemName"`
	Version    string `json:"version"`
}

const ArtifactHubApiEndpoint = "https://artifacthub.io/api/v1/"

var (
	AhSearchEndpoint = ArtifactHubApiEndpoint + "packages/search"
	Limit            = "50"
	HELMKind         = "0"
)

func main() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	ahPackages := make([]artifacthub.AhPackage, 0)
	go func(pkgs *[]artifacthub.AhPackage, quitCh chan struct{}) {
		err := getHelmChartsList(&ahPackages)
		if err != nil {
			fmt.Println(err)
		}
		quitCh <- struct{}{}
	}(&ahPackages, quit)

	lastIndex := 0

	for {
		select {
		case <-ticker.C:
			if len(ahPackages) != 0 {
				err := writeExternalSystems(ahPackages)
				if err != nil {
					fmt.Println(err)
				}
				err = writeComponents(ahPackages[lastIndex+1:])
				if err != nil {
					fmt.Println(err)
				}
				lastIndex = len(ahPackages) - 1
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

type ComponentStruct struct {
	PackageName string               `json:"name"`
	Components  []v1alpha1.Component `json:"components"`
}

func writeComponents(pkgs []artifacthub.AhPackage) error {
	comps := make([]ComponentStruct, 0)
	content, err := os.ReadFile("./output/components.json")
	if err == nil {
		err = json.Unmarshal(content, &comps)
		if err != nil {
			return err
		}
	}
	for _, ap := range pkgs {
		cmps_gen, err := ap.GenerateComponents()
		if err != nil {
			continue
		}
		cmp := ComponentStruct{
			PackageName: ap.Name,
			Components:  cmps_gen,
		}
		comps = append(comps, cmp)
	}
	val, err := json.MarshalIndent(comps, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile("./output/components.json", val, 0666)
	if err != nil {
		return err
	}
	return nil
}

func writeExternalSystems(pkgs []artifacthub.AhPackage) error {
	extSytms := make([]ExternalSystem, 0)
	for _, ap := range pkgs {
		extSytms = append(extSytms, ExternalSystem{
			SystemName: ap.Name,
			Version:    ap.Version,
		})
	}
	val, err := json.MarshalIndent(extSytms, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile("./output/external_systems.json", val, 0666)
	if err != nil {
		return err
	}
	return nil
}

func getKindFacet(facetsSlice []map[string]interface{}) map[string]interface{} {
	for _, facet := range facetsSlice {
		if facet["title"] == "Kind" {
			return facet
		}
	}
	return nil
}

func getTotalChartsCount(apiResponse map[string]([]map[string]interface{})) int {
	kindFacet := getKindFacet(apiResponse["facets"])
	if kindFacet == nil {
		return 0
	}
	options := kindFacet["options"].([]interface{})
	for _, option := range options {
		option := option.(map[string]interface{})
		if option["id"].(float64) == 0 {
			return int(option["total"].(float64))
		}
	}
	return 0
}

func getPackagesFromRes(res map[string]([]map[string]interface{})) []artifacthub.AhPackage {
	ahPackages := make([]artifacthub.AhPackage, 0)
	pkgs := res["packages"]
	for _, pkg := range pkgs {
		repo := pkg["repository"].(map[string]interface{})
		ahPkg := artifacthub.AhPackage{Name: pkg["name"].(string), Repository: repo["name"].(string)}
		ahPkg.UpdatePackageData()
		crds, err := manifests.GetCrdsFromHelm(ahPkg.Url)
		if err == nil && len(crds) != 0 {
			fmt.Println("[Debug] Package with name: ", ahPkg.Name, " has been added to the list of HELM packages")
			ahPackages = append(ahPackages, ahPkg)
		}
	}
	return ahPackages
}

func getHelmChartsList(ahPackages *[]artifacthub.AhPackage) error {
	res, err := getApiResponse(0, true)
	if err != nil {
		return err
	}
	pkgs := getPackagesFromRes(res)
	*ahPackages = append(*ahPackages, pkgs...)
	totalChartsCount := getTotalChartsCount(res)
	loopCount := int((totalChartsCount / 50))
	for i := 0; i < loopCount; i++ {
		res, err := getApiResponse(50*(i+1), false)
		if err != nil {
			continue
		}
		pkgs := getPackagesFromRes(res)
		*ahPackages = append(*ahPackages, pkgs...)
	}
	return nil
}

func getApiResponse(offset int, facet bool) (map[string]([]map[string]interface{}), error) {
	reqUrl := fmt.Sprintf("%s?offset=%v&limit=%s&facets=%v&kind=%v&deprecated=false", AhSearchEndpoint, offset, Limit, facet, HELMKind)
	resp, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("status code %d for %s", resp.StatusCode, reqUrl)
		return nil, err
	}
	defer resp.Body.Close()
	var res map[string]([]map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
