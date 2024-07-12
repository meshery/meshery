package registration

import (

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"

)

type PackagingUnit struct {
	model v1beta1.Model
	components []v1beta1.ComponentDefinition
	relationships []v1alpha2.RelationshipDefinition
	policies []v1beta1.PolicyDefinition
}

// Register will accept a dir, oci or tar. Register assumes that whatever is given in its argument is a single unit of
// packaging (model) and interpret it like that. inputs that have models inside models are not considered valid. This is
// enforce the fact that registration should happen through the unit of packaing (model)
// Given that, it does not care about the folder structure as long as we have one file defining the `model` and others
// `components` etc.
// For bulk registration, use BulkRegister
// Using any because of lack of type strength in go
func Register(entity RegisterableEntity) error {
	// get the packaging units
	pu, err := entity.PkgUnit()
	if(err != nil){
		// given input is not a valid model, or it should not be packaged like this
		return err
	}
	return register(pu)
}

// This function will take meshery's unit of packaging, a model in any format (tar.gz, oci, or a directory) and register
// it
func register(pkg PackagingUnit) error {
	// 1. Register the model

	// 2. Register other units
	return nil
}
























