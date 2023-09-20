package stages

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/utils"
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
		flatSvc := map[string]interface{}{}
		utils.FlattenMap("", utils.ToMapStringInterface(data.Pattern), flatSvc)
		if !skipPrintLogs {
			fmt.Printf("%+#v\n", flatSvc)
		}
		err = fill(data.Pattern, flatSvc)
		if next != nil {
			next(data, err)
		}
	}
}

func fill(p *core.Pattern, flatSvc map[string]interface{}) error {
	var errs []error
	for _, v := range p.Services {
		if err := fillDependsOn(v, flatSvc); err != nil {
			errs = append(errs, err)
		}
		if err := fillNamespace(v, flatSvc); err != nil {
			errs = append(errs, err)
		}
		if err := fillVersion(v, flatSvc); err != nil {
			errs = append(errs, err)
		}
		if err := fillSettings(v, flatSvc); err != nil {
			errs = append(errs, err)
		}
		if err := fillTraits(v, flatSvc); err != nil {
			errs = append(errs, err)
		}
		if err := fillType(v, flatSvc); err != nil {
			errs = append(errs, err)
		}
	}

	return mergeErrors(errs)
}

func fillDependsOn(svc *core.Service, flatSvc map[string]interface{}) error {
	for i, d := range svc.DependsOn {
		k, ok := matchPattern(d)
		if !ok {
			continue
		}

		val, found := flatSvc[k]
		if !found {
			return fmt.Errorf("invalid reference query: %s", k)
		}

		cval, ok := val.(string)
		if !ok {
			return fmt.Errorf("resolved reference query [%s] does not return string", k)
		}

		svc.DependsOn[i] = cval
	}

	return nil
}
func fillVersion(svc *core.Service, flatSvc map[string]interface{}) error {
	nsKey, ok := matchPattern(svc.Version)
	if !ok {
		return nil
	}

	val, found := flatSvc[nsKey]
	if !found {
		return fmt.Errorf("invalid reference query: %s", nsKey)
	}

	vVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("resolved reference query [%s] does not return string", nsKey)
	}

	svc.Version = vVal
	return nil
}
func fillNamespace(svc *core.Service, flatSvc map[string]interface{}) error {
	nsKey, ok := matchPattern(svc.Namespace)
	if !ok {
		return nil
	}

	val, found := flatSvc[nsKey]
	if !found {
		return fmt.Errorf("invalid reference query: %s", nsKey)
	}

	nsVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("resolved reference query [%s] does not return string", nsKey)
	}

	svc.Namespace = nsVal
	return nil
}

func fillType(svc *core.Service, flatSvc map[string]interface{}) error {
	tKey, ok := matchPattern(svc.Type)
	if !ok {
		return nil
	}

	val, found := flatSvc[tKey]
	if !found {
		return fmt.Errorf("invalid reference query: %s", tKey)
	}

	tVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("resolved reference query [%s] does not return string", tKey)
	}

	svc.Type = tVal
	return nil
}

func fillSettings(svc *core.Service, flatSvc map[string]interface{}) (err error) {
	svc.Settings, err = fillMap(svc.Settings, flatSvc)
	return
}

func fillTraits(svc *core.Service, flatSvc map[string]interface{}) (err error) {
	svc.Traits, err = fillMap(svc.Traits, flatSvc)
	return
}

func fillMap(mp map[string]interface{}, flatSvc map[string]interface{}) (map[string]interface{}, error) {
	var _fillMap func(mp map[string]interface{}) (map[string]interface{}, error)

	_fillMap = func(mp map[string]interface{}) (map[string]interface{}, error) {
		for k, v := range mp {
			switch cNode := v.(type) {
			case string:
				val, ok, err := fillMapString(cNode, flatSvc)
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
						val, ok, err := fillMapString(ccNode, flatSvc)
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

func fillMapString(str string, flatSvc map[string]interface{}) (string, bool, error) {
	res, ok := matchPattern(str)
	if !ok {
		return "", false, nil
	}

	val, found := flatSvc[res]
	if !found {
		return "", false, fmt.Errorf("invalid reference query: %s", res)
	}

	cval, ok := val.(string)
	if !ok {
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
