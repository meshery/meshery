package stages

import (
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
func Filler(data *Data, err error, next ChainStageNextFunction) {
	// Flatten the service map to perform queries
	flatSvc := map[string]interface{}{}
	utils.FlattenMap("", utils.ToMapStringInterface(data.Pattern), flatSvc)

	fill(data.Pattern, flatSvc)

	if next != nil {
		next(data, nil)
	}
}

func fill(p *core.Pattern, flatSvc map[string]interface{}) {
	for _, v := range p.Services {
		fillDependsOn(v, flatSvc)
		fillNamespace(v, flatSvc)
		fillSettings(v, flatSvc)
		fillTraits(v, flatSvc)
		fillType(v, flatSvc)
	}
}

func fillDependsOn(svc *core.Service, flatSvc map[string]interface{}) {
	for i, d := range svc.DependsOn {
		k, ok := matchPattern(d)
		if !ok {
			continue
		}

		val, ok := flatSvc[k].(string)
		if !ok {
			d = ""
			continue
		}

		svc.DependsOn[i] = val
	}
}

func fillNamespace(svc *core.Service, flatSvc map[string]interface{}) {
	nsKey, ok := matchPattern(svc.Namespace)
	if !ok {
		svc.Namespace = ""
		return
	}

	nsVal, ok := flatSvc[nsKey].(string)
	if !ok {
		svc.Namespace = ""
		return
	}

	svc.Namespace = nsVal
}

func fillType(svc *core.Service, flatSvc map[string]interface{}) {
	tKey, ok := matchPattern(svc.Type)
	if !ok {
		svc.Type = ""
		return
	}

	tVal, ok := flatSvc[tKey].(string)
	if !ok {
		svc.Type = ""
		return
	}

	svc.Type = tVal
}

func fillSettings(svc *core.Service, flatSvc map[string]interface{}) {
	svc.Settings = fillMap(svc.Settings, flatSvc)
}

func fillTraits(svc *core.Service, flatSvc map[string]interface{}) {
	svc.Traits = fillMap(svc.Traits, flatSvc)
}

func fillMap(mp map[string]interface{}, flatSvc map[string]interface{}) map[string]interface{} {
	var _fillMap func(mp map[string]interface{}) map[string]interface{}

	_fillMap = func(mp map[string]interface{}) map[string]interface{} {
		for k, v := range mp {
			switch cNode := v.(type) {
			case string:
				res, ok := matchPattern(cNode)
				if !ok {
					continue
				}

				val, ok := flatSvc[res].(string)
				if !ok {
					continue
				}

				mp[k] = val
			case map[string]interface{}:
				mp[k] = _fillMap(cNode)
			}
		}

		return mp
	}

	return _fillMap(mp)
}

// matchPattern takes in a string and tests it against the pattern
// if match is successfull it returns the string after removing the pattern
// specific details from it
func matchPattern(str string) (string, bool) {
	res := FillerRegex.FindString(str)
	if res == "" {
		return res, false
	}

	return strings.TrimSuffix(strings.TrimPrefix(res, "$(#ref."), ")"), true
}
