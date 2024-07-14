package registration

import (
	"fmt"
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

/*
	The directory should contain one and only one `model`.
	A directory containing multiple `model` will be invalid.
*/
func NewDir(path string) Dir {
	return Dir{dirpath: path}
}


/* PkgUnit parses all the files inside the directory and finds out if they are any valid meshery definitions. Valid meshery definitions are added to the packagingUnit struct.
*/
func (d Dir) PkgUnit() (_ packagingUnit, err error){
	pkg := packagingUnit{}
	// check if the given is a directory
	_, err = os.ReadDir(d.dirpath)
	if(err != nil){
		return pkg, ErrDirPkgUnitParseFail(d.dirpath, fmt.Errorf("Could not read the directory: %e", err))
	}
	err = filepath.Walk(d.dirpath, func (path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if(f.IsDir()){
				return nil
		}
		byt, _ := os.ReadFile(path)
		var e entity.Entity
		if(filepath.Ext(path) == ".yaml"){
			e, err = getEntity(byt, "yaml")
			if(err != nil){
				// we skip unrecognisable entities
				RegLog.invalidDefinitions[path] = ErrInvalidMeshmodelDefinition(path, err)
				return nil
			}
		} else if(filepath.Ext(path) == ".json"){
			e, err = getEntity(byt, "json")
			if(err != nil){
				RegLog.invalidDefinitions[path] = ErrInvalidMeshmodelDefinition(path, err)
				return nil
			}
		} else if(filepath.Ext(path) == ".cue"){
		} else {
			return nil
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
	if err != nil {
		return pkg, ErrDirPkgUnitParseFail(d.dirpath, fmt.Errorf("Could not completely walk the file tree: %e", err))
	}
	if (reflect.ValueOf(pkg.model).IsZero()){
		err = fmt.Errorf("Cannot find `model` in the directory. Directory without a `model` definition is invalid.")
		RegLog.invalidDefinitions[d.dirpath] = err
		return pkg, err
	}
	return pkg, err
}