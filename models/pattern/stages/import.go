package stages

import (
	"crypto/sha1"
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

func Import(prov ServiceInfoProvider, act ServiceActionProvider) ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}
		data.Lock.Lock()
		//maintain an auxilliary list of dependencies
		//---

		imported := make(map[string]bool, 0)
		dependencies := make(map[string][]string) //key is the service name and value is the array of all services that depends on it
		replacedBy := make(map[string]string, 0)
		err = expandImportOnServices(data.Pattern, imported, dependencies, replacedBy)
		if err != nil {
			act.Terminate(err)
			return
		}
		data.Other["replacedBy"] = replacedBy
		data.Lock.Unlock()
		b, _ := yaml.Marshal(data.Pattern)
		fmt.Println("AFTER FIRST STAGE: \n", string(b))
		if next != nil {
			next(data, nil)
		}
	}
}

func expandImportOnServices(pattern *core.Pattern, parentimported map[string]bool, dependencies map[string][]string, replacedBy map[string]string) error {
	for key, svc := range pattern.Services {
		imported := make(map[string]bool)
		for k := range parentimported {
			imported[k] = true
		}
		loc, ok := matchImportPattern(svc.Type)
		if !ok {
			continue
		}
		var svcs map[string]*core.Service
		var err error
		if pattern.Vars == nil {
			pattern.Vars = make(map[string]interface{})
		}
		svcs, err = expandImportedPatternToServices(key, svc, loc, pattern.Vars, imported, dependencies, replacedBy)
		if err != nil {
			return err
		}
		delete(pattern.Services, key)
		for name, s := range svcs {
			pattern.Services[name] = s
		}
	}
	for _, svc := range pattern.Services {
		svc.DependsOn = nil
	}
	for _, svc := range pattern.Services {
		for i, d := range svc.DependsOn {
			if replacedBy[d] != "" {
				svc.DependsOn[i] = replacedBy[d]
			}
		}
	}
	return nil
}
func expandImportedPatternToServices(name string, svc *core.Service, loc string, vars map[string]interface{}, imported map[string]bool, dependencies map[string][]string, replacedBy map[string]string) (map[string]*core.Service, error) {
	if imported[loc] {
		return nil, errors.New("[Service " + svc.Name + "]Circular Import detected for URL " + loc)
	}

	pattern, err := getPatternFromLocation(loc)
	if err != nil {
		return nil, err
	}
	imported[loc] = true
	for oldName, oldsvc := range pattern.Services { //change the names of the services for uniqueness
		newname := strings.ToLower(pattern.Name) + getHash(oldsvc)
		oldsvc.Name = newname
		oldsvc.DependsOn = svc.DependsOn
		pattern.Services[newname] = oldsvc
		replacedBy[oldName] = newname
		delete(pattern.Services, oldName)
	}
	err = expandImportOnServices(&pattern, imported, dependencies, replacedBy)
	if err != nil {
		return nil, err
	}
	for k, v := range pattern.Vars {
		if svc.Settings[k] != nil {
			vars[k] = svc.Settings[k]
		} else {
			vars[k] = v
		}
	}
	return pattern.Services, nil
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

func dependenciesToDependsOn(dependencies map[string][]string, svc map[string]*core.Service) {
	for d, dep := range dependencies {
		for _, y := range dep {
			if svc[y] != nil {
				svc[y].DependsOn = append(svc[y].DependsOn, d)
			}
		}
	}
}

func dependsOnToDependencies(dependencies map[string][]string, svcs map[string]*core.Service) {
	for _, svc := range svcs {
		for _, dsvc := range svc.DependsOn {
			dependencies[dsvc] = append(dependencies[dsvc], svc.Name)
		}
	}
}

func replaceDependsOnInAllSvcsWithNewSvcName(old, new string, svcs map[string]*core.Service) {
	for _, svc := range svcs {
		for i, d := range svc.DependsOn {
			if d == old {
				fmt.Println("old=", old, " new=", new)
				svc.DependsOn[i] = new
			}
		}
	}
}
