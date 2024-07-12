package registration

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
)

// this smells bad.. refactor it
func getEntity(byt []byte, filetype string) (entity.Entity, error) {
    var versionMeta v1beta1.VersionMeta
    switch filetype {
    case "yaml":
		err := yaml.Unmarshal(byt, &versionMeta)
		if(err != nil){
			return nil, err
		}
    case "json":
		err := json.Unmarshal(byt, &versionMeta)
		if(err != nil){
			return nil, err
		}
    }
	switch (versionMeta.SchemaVersion) {
		case v1beta1.ComponentSchemaVersion:
			var compDef v1beta1.ComponentDefinition
			switch filetype {
				case "yaml":
					err := yaml.Unmarshal(byt, &compDef)
					if(err == nil){
						return &compDef, nil
					}
				case "json":
					err := json.Unmarshal(byt, &compDef)
					if(err == nil){
						return &compDef, nil
					}
			}
		case v1beta1.ModelSchemaVersion:
			var model v1beta1.Model
			switch filetype {
				case "yaml":
					err := yaml.Unmarshal(byt, &model)
					if(err == nil){
						return &model, nil
					}
				case "json":
					err := json.Unmarshal(byt, &model)
					if(err == nil){
						return &model, nil
					}
			}

		case v1alpha2.RelationshipSchemaVersion:
			var rel v1alpha2.RelationshipDefinition
			switch filetype {
				case "yaml":
					err := yaml.Unmarshal(byt, &rel)
					if(err == nil){
						return &rel, nil
					}
				case "json":
					err := json.Unmarshal(byt, &rel)
					if(err == nil){
						return &rel, nil
					}
			}
	}
	return nil, fmt.Errorf("")
}