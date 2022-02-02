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

type servicewrapper struct {
	parentname string
	svc        *core.Service
	name       string
	imported   map[string]bool
}
type servicestack struct {
	sws []*servicewrapper
}

func (ss *servicestack) push(s *servicewrapper) {
	ss.sws = append(ss.sws, s)
}
func (ss *servicestack) pop() (s *servicewrapper) {
	if ss.isEmpty() {
		return nil
	}
	svc := ss.sws[len(ss.sws)-1]
	ss.sws = ss.sws[0 : len(ss.sws)-1]
	return svc
}
func (ss *servicestack) isEmpty() bool {
	if len(ss.sws) == 0 {
		return true
	}
	return false
}
func (ss *servicestack) runOnStack(fn func(svc *core.Service)) {
	for _, s := range ss.sws {
		fn(s.svc)
	}
}
func Import(prov ServiceInfoProvider, act ServiceActionProvider) ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}
		data.Lock.Lock()
		//maintain an auxilliary list of dependencies
		//---
		nonImportingServiceStack := &servicestack{sws: nil}
		importingServiceStack := &servicestack{}
		vars := map[string]interface{}{}
		for _, svc := range data.Pattern.Services {
			importingServiceStack.push(&servicewrapper{
				parentname: "",
				svc:        svc,
				name:       svc.Name,
				imported:   map[string]bool{},
			})
		}
		err = process(importingServiceStack, nonImportingServiceStack, vars)
		if err != nil {
			act.Terminate(err)
			return
		}
		data.Pattern.Services = stackToServices(nonImportingServiceStack)
		data.Pattern.Vars = vars
		patternYaml, _ := yaml.Marshal(data.Pattern)
		*data.Pattern, err = core.NewPatternFile(patternYaml)
		if err != nil {
			act.Terminate(err)
			return
		}
		data.Lock.Unlock()
		b, _ := yaml.Marshal(data.Pattern)
		fmt.Println("AFTER FIRST STAGE: \n", string(b))
		if next != nil {
			next(data, nil)
		}
	}
}

func stackToServices(s *servicestack) (services map[string]*core.Service) {
	services = make(map[string]*core.Service)
	if s.isEmpty() {
		fmt.Println("FUCK")
		return
	}
	for _, svcwrap := range s.sws {
		services[svcwrap.name] = svcwrap.svc
	}
	return
}
func process(imp *servicestack, nonimp *servicestack, vars map[string]interface{}) error {

	for !imp.isEmpty() {
		sw := imp.pop()
		url, ok := matchImportPattern(sw.svc.Type)
		if !ok {
			fmt.Println("NOT GOING AHEAD FOR ", sw.name)
			nonimp.push(sw)
			continue
		}
		fmt.Println("GOING AHEAD FOR ", sw.name)
		p, err := getPatternFromLocation(url)
		if err != nil {
			return err
		}
		var svcws []*servicewrapper
		var oldtonew = make(map[string]string)
		for name, svc := range p.Services {
			fmt.Println("FOR ", sw.name)
			svcw := &servicewrapper{svc: svc}
			svcw.name = svc.Name + getHash(svc)
			for _, svco := range p.Services {
				replaceInDependsOn(name, svcw.name, &svco.DependsOn) // adding resolving ref
			}
			svcw.parentname = sw.name
			svcw.imported = sw.imported
			svcw.svc = svc
			svcw.svc.Name = svcw.name
			svcw.svc.DependsOn = append(svcw.svc.DependsOn, sw.svc.DependsOn...)
			oldtonew[name] = svcw.name
			svcws = append(svcws, svcw)
		}
		for old, new := range oldtonew {
			err := changeReferenceInPattern(old, new, svcws)
			if err != nil {
				return err
			}
		}
		for k, v := range sw.svc.Settings {
			vars[k] = v
		}
		pushSvcToStack(svcws, imp, nonimp)
	}

	return nil
}

func changeReferenceInPattern(old string, new string, sw []*servicewrapper) error {
	fmt.Println("replacing ", old, " with ", new)
	for _, svcwrap := range sw {
		s := svcwrap.svc
		yamlp, err := yaml.Marshal(s)
		if err != nil {
			return err
		}
		yamls := string(yamlp)
		newsvc := core.Service{}
		yamls = strings.ReplaceAll(yamls, "$(#ref.services."+old+".", "$(#ref.services."+new+".")
		err = yaml.Unmarshal([]byte(yamls), &newsvc)
		if err != nil {
			return err
		}
		svcwrap.svc = &newsvc
	}
	return nil
}
func pushSvcToStack(sw []*servicewrapper, imst *servicestack, nonimst *servicestack) {
	for _, s := range sw {
		imst.runOnStack(func(svc *core.Service) {
			replaceInDependsOn(s.parentname, s.name, &svc.DependsOn)
		})
		nonimst.runOnStack(func(svc *core.Service) {
			replaceInDependsOn(s.parentname, s.name, &svc.DependsOn)
		})
	}

	for _, s := range sw {
		if _, ok := matchImportPattern(s.svc.Type); !ok {
			fmt.Println("PUSHING HERE ", s.name)
			nonimst.push(s)
			continue
		}
		imst.push(s)
	}
}
func replaceInDependsOn(old string, new string, do *[]string) {
	for i, d := range *do {
		if d == old {
			(*do)[i] = new
		}
	}
}
func expandImportOnServices(pattern *core.Pattern, parentimported map[string]bool, replacedBy map[string]string) error {
	servicenodes := map[string]*node{}
	for key, svc := range pattern.Services {
		imported := make(map[string]bool)
		for k := range parentimported {
			imported[k] = true
		}
		servicenodes[key] = &node{
			name: key,
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
		svcs, err = expandImportedPatternToServices(key, svc, loc, pattern.Vars, imported, replacedBy)
		if err != nil {
			return err
		}
		delete(pattern.Services, key)
		for name, s := range svcs {
			nod := servicenodes[key]
			nod.imports = append(nod.imports, node{name: name})
			servicenodes[key] = nod
			fmt.Println(servicenodes[key].imports)
			pattern.Services[name] = s
		}
	}
	for _, s := range pattern.Services {
		for k, n := range servicenodes {
			i := contains(s.DependsOn, k)
			if i != -1 {
				n.updateDependency(&s.DependsOn, i)
			}
		}

	}
	return nil
}
func expandImportedPatternToServices(name string, svc *core.Service, loc string, vars map[string]interface{}, imported map[string]bool, replacedBy map[string]string) (map[string]*core.Service, error) {
	if imported[loc] {
		return nil, errors.New("[Service " + svc.Name + "]Circular Import detected for URL " + loc)
	}

	pattern, err := getPatternFromLocation(loc)
	if err != nil {
		return nil, err
	}
	imported[loc] = true
	for oldName, oldsvc := range pattern.Services { //change the names of the services for uniqueness
		newname := strings.ToLower(oldName) + getHash(oldsvc)
		oldsvc.Name = newname
		oldsvc.DependsOn = svc.DependsOn
		pattern.Services[newname] = oldsvc
		replacedBy[oldName] = newname
		delete(pattern.Services, oldName)
	}
	err = expandImportOnServices(&pattern, imported, replacedBy)
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
		if resp.StatusCode != http.StatusOK {
			return p, errors.New(fmt.Sprintf("Got non ok response %d for URL: %s", resp.StatusCode, loc))
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

//As the name of services can change, we need a way to ensure that the depends-on constraint is not broken
// So during the start of evaluation of each pattern file intializes slice of this below data structure. Where each node is a service.
// Normal services are nodes which point to null.
// The service which imports another service points to slice of nodes, and recursively so on and so forth.
// Just after returning all services from patterns, we will use this data structure to resolve depends-on field with the following logic:
// For each service: For each depends on svc: We start from the Node[svc]. We go to each of its children and run UpdateDepency by passing depends-on
// Update dependency work the following way: If this node point to no further nodes- append the svc name in depends-on,otherwise recursively call UpdateDependency

type node struct {
	name    string
	imports []node
}

func (n *node) updateDependency(dependson *[]string, i int) {
	fmt.Println("name is ", n.name, " and import is ", n.imports)
	if len(n.imports) == 0 { //base case
		fmt.Println("ADDING IN DEPENDS ON: ", n.name)
		*dependson = append(*dependson, n.name)
		return
	}
	fmt.Println("WILL BE REMOVING ", (*dependson)[i])
	*dependson = append((*dependson)[:i], (*dependson)[i+1:]...)
	fmt.Println("AFTER DELETION ", *dependson)
	fmt.Println("n.imports is ", n.imports)
	for _, n := range n.imports {
		fmt.Println("WILL EXECUTE ON NODE ", n.name, " WITH DEPENDS ON AS ", *dependson)
		n.updateDependency(dependson, i)
	}
}
func contains(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}

	return -1
}
