package stages

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/meshery/meshery/server/models/pattern/planner"
	"github.com/meshery/meshery/server/models/pattern/utils"
	mutils "github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
	"github.com/pkg/errors"
)

const FillerPattern = `\$\(#ref\.([^)]+)\)`

var FillerRegex = regexp.MustCompile(FillerPattern)

// formatScalar formats scalar values (strings, integers, floats, booleans) into strings.
func formatScalar(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
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

// fill resolves references across all components within the design pattern file.
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

		dependsOn[i] = formatScalar(val)
	}

	component.Metadata.AdditionalProperties["dependsOn"] = dependsOn

	return nil
}

// fillVersion resolves references declared in a component's model or component version.
func fillVersion(component *component.ComponentDefinition, flattenedPattern map[string]interface{}) error {
	if versionKey, ok := matchPattern(component.Model.Model.Version); ok {
		val, found := flattenedPattern[versionKey]
		if !found || mutils.IsInterfaceNil(val) {
			return fmt.Errorf("failed to resolve reference \"%s\": %s", "version", versionKey)
		}

		component.Model.Model.Version = formatScalar(val)
	}

	if compVerKey, ok := matchPattern(component.Component.Version); ok {
		val, found := flattenedPattern[compVerKey]
		if !found || mutils.IsInterfaceNil(val) {
			return fmt.Errorf("failed to resolve reference \"%s\": %s", "version", compVerKey)
		}

		component.Component.Version = formatScalar(val)
	}

	return nil
}

// fillNamespace resolves namespace references declared within a component's metadata configuration.
func fillNamespace(component *component.ComponentDefinition, flattenedPattern map[string]interface{}) error {
	_metadata, ok := component.Configuration["metadata"]
	if !ok || mutils.IsInterfaceNil(_metadata) {
		return nil
	}

	configurationMetadata, err := mutils.Cast[map[string]interface{}](_metadata)
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

	configurationMetadata["namespace"] = formatScalar(val)
	component.Configuration["metadata"] = configurationMetadata
	return nil
}

// fillType resolves references declared in a component's kind definition.
func fillType(component *component.ComponentDefinition, flattenedPattern map[string]interface{}) error {
	kindKey, ok := matchPattern(component.Component.Kind)
	if !ok {
		return nil
	}

	val, found := flattenedPattern[kindKey]
	if !found || mutils.IsInterfaceNil(val) {
		return errors.Wrapf(fmt.Errorf("failed to resolve reference"), "failed to resolve \"kind\" reference for \"%s: %s\"", kindKey, component.Component.Kind)
	}

	component.Component.Kind = formatScalar(val)
	return nil
}

// fillConfiguration resolves references within a component's configuration dictionary.
func fillConfiguration(component *component.ComponentDefinition, flattenedPattern map[string]interface{}) (err error) {
	component.Configuration, err = fillMap(component.Configuration, flattenedPattern)
	return
}

// fillMap recursively processes maps and slices to interpolate references in strings.
func fillMap(mp map[string]interface{}, flattenedPattern map[string]interface{}) (map[string]interface{}, error) {
	var _fillMap func(mp map[string]interface{}) (map[string]interface{}, error)
	var _fillSlice func(s []interface{}) ([]interface{}, error)

	_fillSlice = func(s []interface{}) ([]interface{}, error) {
		for i, el := range s {
			switch ccNode := el.(type) {
			case string:
				val, ok, err := fillMapString(ccNode, flattenedPattern)
				if err != nil {
					return s, err
				}
				if ok {
					s[i] = val
				}
			case map[string]interface{}:
				val, err := _fillMap(ccNode)
				if err != nil {
					return s, err
				}
				s[i] = val
			case []interface{}:
				val, err := _fillSlice(ccNode)
				if err != nil {
					return s, err
				}
				s[i] = val
			}
		}
		return s, nil
	}

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
				val, err := _fillSlice(cNode)
				if err != nil {
					return mp, err
				}
				mp[k] = val
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

// fillMapString substitutes all references in a string in a single pass without reprocessing substituted values.
func fillMapString(str string, flattenedPattern map[string]interface{}) (string, bool, error) {
	if !FillerRegex.MatchString(str) {
		return "", false, nil
	}

	var resolveErr error
	result := FillerRegex.ReplaceAllStringFunc(str, func(m string) string {
		if resolveErr != nil {
			return m
		}

		sub := FillerRegex.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}

		val, found := flattenedPattern[sub[1]]
		if !found || mutils.IsInterfaceNil(val) {
			resolveErr = fmt.Errorf("invalid reference query: %s", sub[1])
			return m
		}

		return formatScalar(val)
	})

	if resolveErr != nil {
		return "", false, resolveErr
	}

	return result, true, nil
}

// matchPattern takes in a string and tests it against the pattern
// if match is successful it returns the string after removing the pattern
// specific details from it
func matchPattern(str string) (string, bool) {
	sub := FillerRegex.FindStringSubmatch(str)
	if len(sub) < 2 {
		return "", false
	}

	return sub[1], true
}
