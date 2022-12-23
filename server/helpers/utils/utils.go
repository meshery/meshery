package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

const UI = "../../ui/public/static/img" //Relative to cmd/main.go
var UISVGPaths = make([]string, 1)

func WriteSVGsOnFileSystem(comp v1alpha1.ComponentDefinition) {
	successCreatingDirectory := false
	defer func(s bool) {
		if successCreatingDirectory {
			UISVGPaths = append(UISVGPaths, filepath.Join(UI, comp.Model.Name))
		}
	}(successCreatingDirectory)
	if comp.Metadata["SVG_Color"] != "" {
		path := filepath.Join(UI, comp.Model.Name, "color")
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true
		f, err := os.Create(filepath.Join(path, comp.Model.DisplayName+"-color.svg"))
		if err != nil {
			fmt.Println(err)
			return
		}
		x, ok := comp.Metadata["SVG_Color"].(string)
		if ok {
			_, err = f.WriteString(x)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}
	if comp.Metadata["SVG_White"] != "" {
		path := filepath.Join(UI, comp.Model.Name, "white")
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true
		f, err := os.Create(filepath.Join(path, comp.Model.DisplayName+"-white.svg"))
		if err != nil {
			fmt.Println(err)
			return
		}
		x, ok := comp.Metadata["SVG_White"].(string)
		if ok {
			_, err = f.WriteString(x)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func DeleteSVGsFromFileSystem() {
	for _, path := range UISVGPaths {
		os.RemoveAll(path)
	}
}
