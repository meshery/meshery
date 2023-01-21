package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

// RecursiveCastMapStringInterfaceToMapStringInterface will convert a
// map[string]interface{} recursively => map[string]interface{}
func RecursiveCastMapStringInterfaceToMapStringInterface(in map[string]interface{}) map[string]interface{} {
	res := ConvertMapInterfaceMapString(in)
	out, ok := res.(map[string]interface{})
	if !ok {
		fmt.Println("failed to cast")
	}

	return out
}

// ConvertMapInterfaceMapString converts map[interface{}]interface{} => map[string]interface{}
//
// It will also convert []interface{} => []string
func ConvertMapInterfaceMapString(v interface{}) interface{} {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v2 := range x {
			switch k2 := k.(type) {
			case string:
				m[k2] = ConvertMapInterfaceMapString(v2)
			default:
				m[fmt.Sprint(k)] = ConvertMapInterfaceMapString(v2)
			}
		}
		v = m

	case []interface{}:
		for i, v2 := range x {
			x[i] = ConvertMapInterfaceMapString(v2)
		}

	case map[string]interface{}:
		for k, v2 := range x {
			x[k] = ConvertMapInterfaceMapString(v2)
		}
	}

	return v
}

// FlattenMap flattens the given map and writes the flattened map in the dest
func FlattenMap(prefix string, src map[string]interface{}, dest map[string]interface{}) {
	if len(prefix) > 0 {
		prefix += "."
	}

	for k, v := range src {
		switch cnode := v.(type) {
		case map[string]interface{}:
			if strings.ContainsAny(k, ".") {
				FlattenMap(prefix+"["+k+"]", cnode, dest)
			} else {
				FlattenMap(prefix+k, cnode, dest)
			}
		case []interface{}:
			for i, v := range cnode {
				switch ccNode := v.(type) {
				case map[string]interface{}:
					FlattenMap(prefix+k+"."+strconv.Itoa(i), ccNode, dest)
				default:
					dest[prefix+k+"."+strconv.Itoa(i)] = v
				}
			}
		default:
			dest[prefix+k] = v
		}
	}
}

// ToMapStringInterface takes in data of type interface and returns
// a map[string]interface{} from that data
//
// If the conversion fails then returns an empty map
func ToMapStringInterface(mp interface{}) map[string]interface{} {
	byt, err := json.Marshal(mp)
	if err != nil {
		return map[string]interface{}{}
	}

	res := map[string]interface{}{}
	if err := json.Unmarshal(byt, &res); err != nil {
		return map[string]interface{}{}
	}

	return res
}

func IsClosed(ch chan struct{}) bool {
	if ch == nil {
		return true
	}
	select {
	case <-ch:
		return true
	default:
	}
	return false
}

const UI = "../../ui/public/static/img/meshmodels" //Relative to cmd/main.go
var UISVGPaths = make([]string, 1)
var hashCheckSVG = make(map[string]string)
var mx sync.Mutex

func writeHashCheckSVG(key string, val string) {
	mx.Lock()
	hashCheckSVG[key] = val
	mx.Unlock()
}

func WriteSVGsOnFileSystem(comp *v1alpha1.ComponentDefinition) {
	successCreatingDirectory := false
	defer func(s bool) {
		if successCreatingDirectory {
			UISVGPaths = append(UISVGPaths, filepath.Join(UI, comp.Model.Name))
		}
	}(successCreatingDirectory)
	if comp.Metadata["svg_color"] != "" {
		path := filepath.Join(UI, comp.Model.Name, "color")
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true

		x, ok := comp.Metadata["svg_color"].(string)
		if ok {
			hash := md5.Sum([]byte(x))
			hashString := hex.EncodeToString(hash[:])
			pathsvg := hashCheckSVG[hashString]
			if pathsvg != "" { // the image has already been loaded, point the component to that path
				comp.Metadata["svg_color"] = pathsvg
				goto White
			} else {
				f, err := os.Create(filepath.Join(path, comp.Kind+"-"+comp.Model.DisplayName+"-color.svg"))
				if err != nil {
					fmt.Println(err)
					return
				}
				_, err = f.WriteString(x)
				if err != nil {
					fmt.Println(err)
					return
				}
				comp.Metadata["svg_color"] = getRelativePathForAPI(filepath.Join(comp.Model.Name, "color", comp.Kind+"-"+comp.Model.DisplayName+"-color.svg")) //Replace the actual SVG with path to SVG
				writeHashCheckSVG(hashString, comp.Metadata["svg_color"].(string))
			}
		}
	}
White:
	if comp.Metadata["svg_white"] != "" {
		path := filepath.Join(UI, comp.Model.Name, "white")
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true

		x, ok := comp.Metadata["svg_white"].(string)
		if ok {
			hash := md5.Sum([]byte(x))
			hashString := hex.EncodeToString(hash[:])
			pathsvg := hashCheckSVG[hashString]
			if pathsvg != "" { // the image has already been loaded, point the component to that path
				comp.Metadata["svg_white"] = pathsvg
				return
			}
			f, err := os.Create(filepath.Join(path, comp.Kind+"-"+comp.Model.DisplayName+"-white.svg"))
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = f.WriteString(x)
			if err != nil {
				fmt.Println(err)
				return
			}
			comp.Metadata["svg_white"] = getRelativePathForAPI(filepath.Join(comp.Model.Name, "white", comp.Kind+"-"+comp.Model.DisplayName+"-white.svg")) //Replace the actual SVG with path to SVG
			writeHashCheckSVG(hashString, comp.Metadata["svg_white"].(string))
		}
	}
}

func DeleteSVGsFromFileSystem() {
	for _, path := range UISVGPaths {
		os.RemoveAll(path)
	}
}
func getRelativePathForAPI(path string) string {
	ui := strings.TrimPrefix(UI, "../../")
	return filepath.Join(ui, path)
}
