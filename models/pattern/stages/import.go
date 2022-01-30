package stages

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

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
		imported := make(map[string]bool, 0)
		err = expandImportOnServices(data.Pattern, imported)
		if err != nil {
			act.Terminate(err)
			return
		}
		data.Lock.Unlock()

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
			pattern.Services[name] = s
		}
	}
	return nil
}
func expandImportedPatternToServices(name string, svc *core.Service, loc string, vars map[string]interface{}, imported map[string]bool) (map[string]*core.Service, map[string]interface{}, error) {
	if imported[loc] {
		fmt.Println("TRUE ", loc)
		return nil, nil, errors.New("[Service " + svc.Name + "]Circular Import detected for URL " + loc)
	}

	pattern, err := getPatternFromLocation(loc)
	if err != nil {
		return nil, nil, err
	}
	imported[loc] = true
	fmt.Println("LOC as true", loc)
	for _, svc := range pattern.Services {
		svc.Name = strings.ToLower(pattern.Name) + svc.Name + getHash(5)

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

func getHash(length int) string {
	b := make([]byte, length)
	charset := "abcdefghijklmnopqrstuvwxyz"
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
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

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

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
