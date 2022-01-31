package stages

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/layer5io/meshery/models/pattern/core"
	"github.com/layer5io/meshery/models/pattern/utils"
)

const FillerPattern = `\$\(#ref\..+\)`

var FillerRegex *regexp.Regexp

func init() {
	var err error
	FillerRegex, err = regexp.Compile(FillerPattern)
	if err != nil {
		log.Fatal("failed to compile filler pattern regex")
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
		r := data.Other["replacedBy"]
		replacedBy, ok := r.(map[string]string)
		if !ok {
			replacedBy = nil
		}
		if !skipPrintLogs {
			fmt.Printf("%+#v\n", flatSvc)
		}
		err = fill(data.Pattern, flatSvc, replacedBy)

		if next != nil {
			next(data, err)
		}
	}
}

func fill(p *core.Pattern, flatSvc map[string]interface{}, replacedBy map[string]string) error {
	for _, v := range p.Services {
		if err := fillDependsOn(v, flatSvc, replacedBy); err != nil {
			return err
		}
		if err := fillNamespace(v, flatSvc, replacedBy); err != nil {
			return err
		}
		if err := fillSettings(v, flatSvc, replacedBy); err != nil {
			return err
		}
		if err := fillTraits(v, flatSvc, replacedBy); err != nil {
			return err
		}
		if err := fillType(v, flatSvc, replacedBy); err != nil {
			return err
		}
	}

	return nil
}

func fillDependsOn(svc *core.Service, flatSvc map[string]interface{}, replacedBy map[string]string) error {
	for i, d := range svc.DependsOn {
		k, ok := matchFillerPattern(d)
		if !ok {
			continue
		}
		cval, err := getVal(k, flatSvc, replacedBy)
		if err != nil {
			return err
		}
		svc.DependsOn[i] = cval
	}

	return nil
}

func fillNamespace(svc *core.Service, flatSvc map[string]interface{}, replacedBy map[string]string) error {
	nsKey, ok := matchFillerPattern(svc.Namespace)
	if !ok {
		return nil
	}
	nsVal, err := getVal(nsKey, flatSvc, replacedBy)
	if err != nil {
		return err
	}
	svc.Namespace = nsVal
	return nil
}

func fillType(svc *core.Service, flatSvc map[string]interface{}, replacedBy map[string]string) error {
	tKey, ok := matchFillerPattern(svc.Type)
	if !ok {
		return nil
	}
	tVal, err := getVal(tKey, flatSvc, replacedBy)
	if err != nil {
		return err
	}
	svc.Type = tVal
	return nil
}

func fillSettings(svc *core.Service, flatSvc map[string]interface{}, replacedBy map[string]string) (err error) {
	svc.Settings, err = fillMap(svc.Settings, flatSvc, replacedBy)
	return
}

func fillTraits(svc *core.Service, flatSvc map[string]interface{}, replacedBy map[string]string) (err error) {
	svc.Traits, err = fillMap(svc.Traits, flatSvc, replacedBy)
	return
}

func fillMap(mp map[string]interface{}, flatSvc map[string]interface{}, replacedBy map[string]string) (map[string]interface{}, error) {
	var _fillMap func(mp map[string]interface{}) (map[string]interface{}, error)

	_fillMap = func(mp map[string]interface{}) (map[string]interface{}, error) {
		for k, v := range mp {
			switch cNode := v.(type) {
			case string:
				val, ok, err := fillMapString(cNode, flatSvc, replacedBy)
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
						val, ok, err := fillMapString(ccNode, flatSvc, replacedBy)
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

func fillMapString(str string, flatSvc map[string]interface{}, replacedBy map[string]string) (string, bool, error) {
	res, ok := matchFillerPattern(str)
	if !ok {
		return "", false, nil
	}
	cval, err := getVal(res, flatSvc, replacedBy)
	if err != nil {
		return "", false, err
	}
	return cval, true, nil
}

// matchPattern takes in a string and tests it against the pattern
// if match is successful it returns the string after removing the pattern
// specific details from it
func matchFillerPattern(str string) (string, bool) {
	res := FillerRegex.FindString(str)
	if res == "" {
		return res, false
	}

	return strings.TrimSuffix(strings.TrimPrefix(res, "$(#ref."), ")"), true
}

func getVal(key string, flatSvc map[string]interface{}, replacedBy map[string]string) (string, error) {
	val, found := flatSvc[key]
	if !found {
		newval, ok := replacedBy[strings.Split(key, ".")[1]]
		if !ok {
			return "", fmt.Errorf("invalid reference query: %s", key)
		}
		val = newval
	}

	cval, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("resolved reference query [%s] does not return string", key)
	}
	return cval, nil
}
