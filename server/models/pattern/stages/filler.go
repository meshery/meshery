package stages

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/layer5io/meshery/server/models/pattern/utils"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/pkg/errors"
)

const FillerPattern = `\$\(#ref\..+\)`

var FillerRegex *regexp.Regexp

func init() {
	var err error
	FillerRegex, err = regexp.Compile(FillerPattern)
	if err != nil {
		log.Fatal("failed to compile filler design regex")
	}
}

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
		if err := fillDependsOn(component, flattenedComponent); err != nil {
			err = ErrResolveReference(err)
			errs = append(errs, err)
		}
		if err := fillNamespace(component, flattenedComponent); err != nil {
			err = ErrResolveReference(err)
			errs = append(errs, err)
		}
		if err := fillVersion(component, flattenedComponent); err != nil {
			err = ErrResolveReference(err)
			errs = append(errs, err)
		}
		if err := fillConfiguration(component, flattenedComponent); err != nil {
			err = ErrResolveReference(err)
			errs = append(errs, err)
		}

		if err := fillType(component, flattenedComponent); err != nil {
			err = ErrResolveReference(err)
			errs = append(errs, err)
		}
	}

	return mergeErrors(errs)
}

func fillDependsOn(component *component.ComponentDefinition, flattenedPattern map[string]interface{}) error {
	_dependsOn, ok := component.Metadata.AdditionalProperties["dependsOn"]
	if !ok || mutils.IsInterfaceNil(_dependsOn) {
		return nil
	}

	dependsOn, err := mutils.Cast[[]string](_dependsOn)
	if err != nil {
		return err
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
