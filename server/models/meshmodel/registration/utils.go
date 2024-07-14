package registration

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"gopkg.in/yaml.v2"
)

func unmarshal(byt []byte, filetype string, out interface{}, strict bool) error {
	switch filetype {
    case "yaml":
		decoder := yaml.NewDecoder(bytes.NewReader(byt))
		// this makes sure it errors out if the schema does not match the expected struct
		// decoder.SetStrict(true)
		err := decoder.Decode(out)
		if(err != nil){
			return err
		}
    case "json":
		decoder := json.NewDecoder(bytes.NewReader(byt))
		// this makes sure it errors out if the schema does not match the expected struct
		if(strict){
			decoder.DisallowUnknownFields()
		}
		err := decoder.Decode(out)
		if(err != nil){
			return err
		}
    }
	return nil
}

// TODO: refactor this and use CUE
func getEntity(byt []byte, filetype string) (et entity.Entity, _ error) {
    var versionMeta v1beta1.VersionMeta
	err  := unmarshal(byt, filetype, &versionMeta, false)
	if err != nil {
		return nil, ErrGetEntity(fmt.Errorf("Does not contain versionmeta: %s", err.Error()))
	}
	switch (versionMeta.SchemaVersion) {
		case v1beta1.ComponentSchemaVersion:
			var compDef v1beta1.ComponentDefinition
			err  := unmarshal(byt, filetype, &compDef, false)
			if err != nil {
				return nil, ErrGetEntity(fmt.Errorf("Invalid component definition: %s", err.Error()))
			}
			et = &compDef
		case v1beta1.ModelSchemaVersion:
			var model v1beta1.Model
			err  := unmarshal(byt, filetype, &model, false)
			if err != nil {
				return nil, ErrGetEntity(fmt.Errorf("Invalid model definition: %s", err.Error()))
			}
			et = &model
		case v1alpha2.RelationshipSchemaVersion:
			var rel v1alpha2.RelationshipDefinition
			err  := unmarshal(byt, filetype, &rel, false)
			if err != nil {
				return nil, ErrGetEntity(fmt.Errorf("Invalid relationship definition: %s", err.Error()))
			}
			et = &rel
		default:
			return nil, ErrGetEntity(fmt.Errorf("Not a valid component definition, model definition, or relationship definition"))
	}
	return et, nil
}