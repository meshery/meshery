package registration

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
)

type Dir struct {
	dirpath string
}

func NewDir(path string) Dir {
	return Dir{dirpath: path}
}

func (d Dir) PkgUnit() (PackagingUnit, error){
	pkg := PackagingUnit{}
	filepath.Walk(d.dirpath, func (path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if(f.IsDir()){
				return nil
		}
		byt, _ := os.ReadFile(path)
		var e entity.Entity
		switch filepath.Ext(path) {
		case ".yaml":
			e, err = getEntity(byt, "yaml")
			// TODO: properly handle error
			if(err != nil){
				return nil
			}
		case ".json":
			e, err = getEntity(byt, "json")
			// TODO: properly handle error
			if(err != nil){
				return nil
			}
		case ".cue":
		}
		// set it to pkgunit
		switch e.Type() {
			case entity.Model:
				if !reflect.ValueOf(pkg.model).IsZero() {
					// currently models inside models are not handled
					return nil
				}
				model, _ := e.(*v1beta1.Model)
				pkg.model = *model
			case entity.ComponentDefinition:
				comp, _ := e.(*v1beta1.ComponentDefinition)
				pkg.components = append(pkg.components, *comp)
			case entity.RelationshipDefinition:
				rel, _ := e.(*v1alpha2.RelationshipDefinition)
				pkg.relationships = append(pkg.relationships, *rel)
		}
		return nil
	})
	return pkg, nil
}
