package models

import (
	"fmt"

	"github.com/layer5io/meshkit/models"
	"github.com/layer5io/meshkit/utils/artifacthub"
)

type ArtifactHubPackageManager struct {
	PackageName string
}

func (ahpm ArtifactHubPackageManager) GetPackage() (models.Package, error) {
	// get relevant packages
	pkgs, err := artifacthub.GetAhPackagesWithName(ahpm.PackageName)
	if err != nil {
		return nil, ErrGetPackage(err)
	}
	// update package information
	for i, ap := range pkgs {
		ap.UpdatePackageData()
		pkgs[i] = ap
	}
	// filter only packages with crds
	pkgs = artifacthub.FilterPackagesWithCrds(pkgs)
	pkgs = artifacthub.SortPackagesWithScore(pkgs)
	if len(pkgs) == 0 {
		return nil, ErrGetPackage(fmt.Errorf("could not find any appropriate artifacthub package"))
	}
	return pkgs[0], nil
}
