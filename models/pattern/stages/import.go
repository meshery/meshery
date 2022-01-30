package stages

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/layer5io/meshery/models/pattern/core"
	"gopkg.in/yaml.v2"
)

const ImportPattern = `\$\(#use\s.+\)`

var ImportRegex *regexp.Regexp

var dependencies = make(map[string][]string) //key is the service name and value is the array of all services that depends on it
func Import(prov ServiceInfoProvider, act ServiceActionProvider) ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}
		data.Lock.Lock()
		//maintain an auxilliary list of dependencies
		for _, svc := range data.Pattern.Services {
			for _, dsvc := range svc.DependsOn {
				dependencies[dsvc] = append(dependencies[dsvc], svc.Name)
			}
		}

		imported := make(map[string]bool, 0)
		err = expandImportOnServices(data.Pattern, imported)
		if err != nil {
			act.Terminate(err)
			return
		}
		for _, s := range data.Pattern.Services {
			s.DependsOn = nil
		}
		//refill the dependencies
		for d, dep := range dependencies {
			for _, y := range dep {
				if data.Pattern.Services[y] != nil {
					data.Pattern.Services[y].DependsOn = append(data.Pattern.Services[y].DependsOn, d)
				}
			}
		}
		data.Lock.Unlock()
		b, _ := yaml.Marshal(data.Pattern)
		fmt.Println("AFTER FIRST STAGE: \n", string(b))
		if next != nil {
			next(data, nil)
		}
	}
}

func expandImportOnServices(pattern *core.Pattern, parentimported map[string]bool) error {
	for key, svc := range pattern.Services {
		imported := make(map[string]bool)
		for k := range parentimported {
			imported[k] = true
		}
		loc, ok := matchImportPattern(svc.Type)
		if !ok {
			continue
		}
		dep := dependencies[svc.Name]
		delete(dependencies, svc.Name)
		var svcs map[string]*core.Service
		var err error
		if pattern.Vars == nil {
			pattern.Vars = make(map[string]interface{})
		}
		svcs, _, err = expandImportedPatternToServices(key, svc, loc, pattern.Vars, imported)
		if err != nil {
			return err
		}
		delete(pattern.Services, key)
		for name, s := range svcs {
			dependencies[name] = dep
			pattern.Services[name] = s
		}
	}
	return nil
}
func expandImportedPatternToServices(name string, svc *core.Service, loc string, vars map[string]interface{}, imported map[string]bool) (map[string]*core.Service, map[string]interface{}, error) {
	if imported[loc] {
		return nil, nil, errors.New("[Service " + svc.Name + "]Circular Import detected for URL " + loc)
	}

	pattern, err := getPatternFromLocation(loc)
	if err != nil {
		return nil, nil, err
	}
	imported[loc] = true
	for oldName, svc := range pattern.Services {
		name := strings.ToLower(pattern.Name) + svc.Name + getHash(svc)
		svc.Name = name
		pattern.Services[name] = svc
		delete(pattern.Services, oldName)
	}

	err = expandImportOnServices(&pattern, imported)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range pattern.Vars {
		if svc.Settings[k] != nil {
			vars[k] = svc.Settings[k]
		} else {
			vars[k] = v
		}
	}
	return pattern.Services, vars, nil
}

func getPatternFromLocation(loc string) (p core.Pattern, err error) {
	if strings.HasPrefix(loc, "https://") {
		resp, err := http.Get(loc)
		if err != nil {
			return p, err
		}
		pat, err := io.ReadAll(resp.Body)
		if err != nil {
			return p, err
		}
		return core.NewPatternFile(pat)
	}
	return p, nil
}

func getHash(s *core.Service) string {
	b, _ := yaml.Marshal(s)
	h := sha1.New()
	h.Write(b)
	bs := h.Sum(nil)
	str := string(fmt.Sprintf("%x\n", bs))
	return str[0:8]
}

// matchPattern takes in a string and tests it against the pattern
// if match is successful it returns the string after removing the pattern
// specific details from it
func matchImportPattern(str string) (string, bool) {
	res := ImportRegex.FindString(str)
	if res == "" {
		return res, false
	}

	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(res, "$(#use"), ")")), true
}

func init() {
	var err error
	ImportRegex, err = regexp.Compile(ImportPattern)
	if err != nil {
		log.Fatal("failed to compile filler pattern regex")
	}
}

func tryToMarshalAndFAILTHENYAML(i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		fmt.Println("FAILED ", err.Error())
		b, err := yaml.Marshal(i)
		if err != nil {

		}
		fmt.Println("YAML\n ", string(b))
	}
	fmt.Println("JSON\n", string(b))
}
func convert(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			fmt.Println("THIS MOFOOO ", k, "=", v)
			res[fmt.Sprint(k)] = convert(v2)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	tryToMarshalAndFAILTHENYAML(res)
	return res
}
