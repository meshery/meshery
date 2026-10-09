package stages

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/meshery/meshery/server/models/pattern/planner"
	"github.com/meshery/meshery/server/models/pattern/utils"
	mutils "github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
	"github.com/pkg/errors"
)

const FillerPattern = `\$\(#ref\..+\)`

var FillerRegex = regexp.MustCompile(FillerPattern)

// Filler - filler stage processes the pattern to subsitute Pattern
func Filler(skipPrintLogs bool) ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			next(data, err)
			return
		}
		// Flatten the service map to perform queries
		flattenedComponent := map[string]interface{}{}
		utils.FlattenMap("", utils.ToMapStringInterface(data.Pattern), flattenedComponent)
		if !skipPrintLogs {
			fmt.Printf("%+#v\n", flattenedComponent)
		}
		err = fill(data.Pattern, flattenedComponent)
		if next != nil {
			next(data, err)
		}
	}
}

func fill(p *pattern.PatternFile, flattenedComponent map[string]interface{}) error {
	var errs []error
	for _, component := range p.Components {
		if err := fillDependsOn(p.Name, component, flattenedComponent); err != nil {
			errs = append(errs, resolveReferenceError(err))
		}
		if err := fillNamespace(component, flattenedComponent); err != nil {
			errs = append(errs, resolveReferenceError(err))
		}
		if err := fillVersion(component, flattenedComponent); err != nil {
			errs = append(errs, resolveReferenceError(err))
		}
		if err := fillConfiguration(component, flattenedComponent); err != nil {
			errs = append(errs, resolveReferenceError(err))
		}

		if err := fillType(component, flattenedComponent); err != nil {
			errs = append(errs, resolveReferenceError(err))
		}
	}

	return mergeErrors(errs)
}

// fillDependsOn resolves the references a component's "dependsOn" entries carry
// and leaves the resolved names behind as a []string.
//
// The design arrives here as it was decoded from JSON, where "dependsOn" is not
// a first-class field and so survives as an untyped []interface{}. Normalizing
// it back onto the component is what makes the resolution performed here
// durable: every later reader - the execution plan above all - then sees one
// shape carrying the resolved names rather than the references they started as.
func fillDependsOn(designName string, component *component.ComponentDefinition, flattenedPattern map[string]interface{}) error {
	dependsOn, err := planner.DeclaredDependencies(designName, component)
	if err != nil {
		return err
	}

	if dependsOn == nil {
		return nil
	}

	for i, d := range dependsOn {
		k, ok := matchPattern(d)
		if !ok {
			continue
		}

		val, found := flattenedPattern[k]
		if !found || mutils.IsInterfaceNil(val) {
			return fmt.Errorf("failed to resolve reference \"%s\": %s", "dependsOn", k)
		}

		cval, err := mutils.Cast[string](val)
		if err != nil {
			return errors.Wrapf(err, "failed to resolve reference \"%s\": %s", "dependsOn", val)
		}

		dependsOn[i] = cval
	}

	component.Metadata.AdditionalProperties["dependsOn"] = dependsOn

	return nil
}

func fillVersion(component *component.ComponentDefinition, flattenedPattern map[string]interface{}) error {
	// Refernce to a version?
	// So that if user chooses a comp def of other comps but all other comps in design refer to some other version,
	// So instead of choosing new comp of correct version, user can quickly point to particular version?
	// But it can be that the version are not compatible??
	versionKey, ok := matchPattern(component.Model.Model.Version)
	if !ok {
		return nil
	}

	val, found := flattenedPattern[versionKey]
	if !found || mutils.IsInterfaceNil(val) {
		return fmt.Errorf("failed to resolve reference \"%s\": %s", "version", versionKey)
	}

	vVal, err := mutils.Cast[string](val)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve reference \"%s\": %s", "version", versionKey)
	}

	component.Model.Model.Version = vVal
	return nil
}

func fillNamespace(component *component.ComponentDefinition, flattenedPattern map[string]interface{}) error {
	_metadata, ok := component.Configuration["metadata"]
	if !ok || mutils.IsInterfaceNil(_metadata) {
		return nil
	}

	configurationMetadata, err := mutils.Cast[map[string]interface{}](_metadata)
	fmt.Println("configurationMetadata: ", configurationMetadata, err)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve namespace reference for \"%s: %s\"", component.DisplayName, component.Component.Kind)
	}

	_namespace, ok := configurationMetadata["namespace"]
	if !ok || mutils.IsInterfaceNil(_namespace) {
		return nil
	}

	namespaceKey, err := mutils.Cast[string](_namespace)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve namespace reference for \"%s: %s\"", component.DisplayName, component.Component.Kind)
	}

	nsKey, ok := matchPattern(namespaceKey)
	if !ok {
		return nil
	}

	val, found := flattenedPattern[nsKey]
	if !found || mutils.IsInterfaceNil(val) {
		return fmt.Errorf("invalid reference query: %s", nsKey)
	}

	nsVal, err := mutils.Cast[string](val)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve reference \"%s\": %s", "namespace", nsKey)
	}

	configurationMetadata["namespace"] = nsVal
	component.Configuration["metadata"] = configurationMetadata
	return nil
}

func fillType(component *component.ComponentDefinition, flattenedPattern map[string]interface{}) error {
	kindKey, ok := matchPattern(component.Component.Kind)
	if !ok {
		return nil
	}

	val, found := flattenedPattern[kindKey]
	if !found || mutils.IsInterfaceNil(val) {
		return errors.Wrapf(fmt.Errorf("failed to resolve reference"), "failed to resolve \"kind\" reference for \"%s: %s\"", kindKey, component.Component.Kind)
	}

	kindVal, err := mutils.Cast[string](val)
	if err != nil {
		return errors.Wrapf(err, "failed to resolved \"kind\" reference: %s", kindKey)
	}

	component.Component.Kind = kindVal
	return nil
}

func fillConfiguration(component *component.ComponentDefinition, flattenedPattern map[string]interface{}) (err error) {
	component.Configuration, err = fillMap(component.Configuration, flattenedPattern)
	return
}

func fillMap(mp map[string]interface{}, flattenedPattern map[string]interface{}) (map[string]interface{}, error) {
	var _fillMap func(mp map[string]interface{}) (map[string]interface{}, error)

	_fillMap = func(mp map[string]interface{}) (map[string]interface{}, error) {
		for k, v := range mp {
			switch cNode := v.(type) {
			case string:
				val, ok, err := fillMapString(cNode, flattenedPattern)
				if err != nil {
					return mp, err
				}

				if !ok {
					continue
				}

				mp[k] = val
			case []interface{}:
				for i, el := range cNode {
					switch ccNode := el.(type) {
					case string:
						val, ok, err := fillMapString(ccNode, flattenedPattern)
						if err != nil {
							return mp, err
						}

						if !ok {
							continue
						}

						mp[k].([]interface{})[i] = val
					case map[string]interface{}:
						val, err := _fillMap(ccNode)
						if err != nil {
							return mp, err
						}

						mp[k].([]interface{})[i] = val
					}
				}
			case map[string]interface{}:
				var err error
				mp[k], err = _fillMap(cNode)
				if err != nil {
					return mp, err
				}
			}
		}

		return mp, nil
	}

	return _fillMap(mp)
}

func fillMapString(str string, flattenedPattern map[string]interface{}) (string, bool, error) {
	res, ok := matchPattern(str)
	if !ok {
		return "", false, nil
	}

	val, found := flattenedPattern[res]
	if !found || mutils.IsInterfaceNil(val) {
		return "", false, fmt.Errorf("invalid reference query: %s", res)
	}

	cval, err := mutils.Cast[string](val)
	if err != nil {
		return "", false, fmt.Errorf("resolved reference query [%s] does not return string", res)
	}

	return cval, true, nil
}

// matchPattern takes in a string and tests it against the pattern
// if match is successful it returns the string after removing the pattern
// specific details from it
func matchPattern(str string) (string, bool) {
	res := FillerRegex.FindString(str)
	if res == "" {
		return res, false
	}

	return strings.TrimSuffix(strings.TrimPrefix(res, "$(#ref."), ")"), true
}
