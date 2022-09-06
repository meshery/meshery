package stages

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/layer5io/meshery/server/models/pattern/core"
	"gopkg.in/yaml.v2"
)

const ImportPattern = `\$\(#use\s.+\)`

var ImportRegex *regexp.Regexp

/*
Two stacks are created containing service wrappers. One stack holds all the services which have no further imports.
We pop services from mixed stack, process it, expand it into n services(reset dependsOn using parentname) and push the services into one of the stacks. If the service has no imports, we push it onto nonimporting stack.
Otherwise we push the service onto importingstack. This happens until the importingstack is empty.
*/
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
		//If no imports are there, completely bypass this stage
		for _, svc := range data.Pattern.Services {
			if _, ok := matchImportPattern(svc.Type); ok { //If we found any service with imports
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
				//At the end of this processing, importingServiceStack will be empty and nonImportingServiceStack will have all the services in it with proper dependson and references set
				err = process(importingServiceStack, nonImportingServiceStack, vars)
				if err != nil {
					act.Terminate(err)
					return
				}
				data.Pattern.Services = stackToServices(nonImportingServiceStack)
				data.Pattern.Vars = vars
				patternYaml, err := yaml.Marshal(data.Pattern)
				if err != nil {
					if err != nil {
						act.Terminate(err)
						return
					}
				}
				*data.Pattern, err = core.NewPatternFile(patternYaml)
				if err != nil {
					act.Terminate(err)
					return
				}
				break
			}
		}
		data.Lock.Unlock()
		if next != nil {
			next(data, nil)
		}
	}
}

func stackToServices(s *servicestack) (services map[string]*core.Service) {
	services = make(map[string]*core.Service)
	if s.isEmpty() {
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
			nonimp.push(sw)
			continue
		}
		p, err := getPatternFromLocation(url)
		if err != nil {
			return err
		}
		var svcws []*servicewrapper
		var oldtonew = make(map[string]string)
		for name, svc := range p.Services {
			svcw := &servicewrapper{svc: svc}
			svcw.name = svc.Name + getHashOfService(svc)
			for _, svco := range p.Services {
				replaceInDependsOn(name, svcw.name, &svco.DependsOn)
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
			hashedKey := k + getHashOfPattern(&p) //change variable names to avoid exported variable name conflicts during import
			err := changeVarsInPattern(k, hashedKey, svcws)
			if err != nil {
				return fmt.Errorf("could not import variables from imported pattern")
			}
			vars[hashedKey] = v
		}
		pushSvcToStack(svcws, imp, nonimp)
	}

	return nil
}

func changeReferenceInPattern(old string, new string, sw []*servicewrapper) error {
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
func changeVarsInPattern(old string, new string, sw []*servicewrapper) error {
	for _, svcwrap := range sw {
		s := svcwrap.svc
		yamlp, err := yaml.Marshal(s)
		if err != nil {
			return err
		}
		yamls := string(yamlp)
		newsvc := core.Service{}
		yamls = strings.ReplaceAll(yamls, "$(#ref.vars."+old+")", "$(#ref.vars."+new+")")
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
			addInDependsOn(s.parentname, s.name, &svc.DependsOn)
		})
		nonimst.runOnStack(func(svc *core.Service) {
			addInDependsOn(s.parentname, s.name, &svc.DependsOn)
		})
	}

	for _, s := range sw {
		imst.runOnStack(func(svc *core.Service) {
			removeInDependsOn(s.parentname, &svc.DependsOn)
		})
		nonimst.runOnStack(func(svc *core.Service) {
			removeInDependsOn(s.parentname, &svc.DependsOn)
		})
	}
	for _, s := range sw {
		if _, ok := matchImportPattern(s.svc.Type); !ok {
			nonimst.push(s)
			continue
		}
		imst.push(s)
	}
}
func addInDependsOn(old string, new string, do *[]string) {
	for _, d := range *do {
		if d == old {
			(*do) = append((*do), new)
		}
	}
}
func replaceInDependsOn(old string, new string, do *[]string) {
	for i, d := range *do {
		if d == old {
			(*do)[i] = new
		}
	}
}
func removeInDependsOn(old string, do *[]string) {
	for i, d := range *do {
		if d == old {
			(*do) = append((*do)[0:i], (*do)[i+1:]...)
		}
	}
}
func replaceParentInDependsOn(parent string, childrens []string, do *[]string) {
	for i, d := range *do {
		if d == parent {
			(*do) = append((*do)[0:i], (*do)[i+1:]...)
			(*do) = append((*do), childrens...)
			break
		}
	}
}
func getPatternFromLocation(loc string) (p core.Pattern, err error) {
	if strings.HasPrefix(loc, "https://") {
		resp, err := http.Get(loc)
		if err != nil {
			return p, err
		}
		if resp.StatusCode != http.StatusOK {
			return p, fmt.Errorf("got non ok HTTP response %d for URL: %s", resp.StatusCode, loc)
		}
		pat, err := io.ReadAll(resp.Body)
		if err != nil {
			return p, err
		}
		return core.NewPatternFile(pat)
	}
	return p, nil
}

func getHashOfService(s *core.Service) string {
	b, _ := yaml.Marshal(s)
	h := sha1.New()
	h.Write(b)
	bs := h.Sum(nil)
	str := string(fmt.Sprintf("%x\n", bs))
	return str[0:8]
}

func getHashOfPattern(s *core.Pattern) string {
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
