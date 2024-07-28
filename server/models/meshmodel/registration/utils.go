package registration

import (
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"gopkg.in/yaml.v2"
)

func unmarshal(byt []byte, out interface{}) error {
	err := json.Unmarshal(byt, out)
	if(err != nil){
		err = yaml.Unmarshal(byt, out)
		if(err != nil){
			return fmt.Errorf("Not a valid YAML or JSON")
		}
	}
	return nil
}

// TODO: refactor this and use CUE
func getEntity(byt []byte) (et entity.Entity, _ error) {
    var versionMeta v1beta1.VersionMeta
	err  := unmarshal(byt, &versionMeta)
	if err != nil || versionMeta.SchemaVersion == "" {
		return nil, ErrGetEntity(fmt.Errorf("Does not contain versionmeta"))
	}
	switch (versionMeta.SchemaVersion) {
		case v1beta1.ComponentSchemaVersion:
			var compDef v1beta1.ComponentDefinition
			err  := unmarshal(byt, &compDef)
			if err != nil {
				return nil, ErrGetEntity(fmt.Errorf("Invalid component definition: %s", err.Error()))
			}
			et = &compDef
		case v1beta1.ModelSchemaVersion:
			var model v1beta1.Model
			err  := unmarshal(byt,&model)
			if err != nil {
				return nil, ErrGetEntity(fmt.Errorf("Invalid model definition: %s", err.Error()))
			}
			et = &model
		case v1alpha2.RelationshipSchemaVersion:
			var rel v1alpha2.RelationshipDefinition
			err  := unmarshal(byt,&rel)
			if err != nil {
				return nil, ErrGetEntity(fmt.Errorf("Invalid relationship definition: %s", err.Error()))
			}
			et = &rel
		default:
			return nil, ErrGetEntity(fmt.Errorf("Not a valid component definition, model definition, or relationship definition"))
	}
	return et, nil
}