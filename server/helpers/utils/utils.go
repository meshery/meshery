package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils"
)

const (
	HelmChartURL          = "https://meshery.io/charts/"
	HelmChartName         = "meshery"
	HelmChartOperatorName = "meshery-operator"
	MesheryFolder         = ".meshery"
	ManifestsFolder       = "manifests"
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
func writeSVGHelper(metadata map[string]interface{}, dirname, filename string) {
	filename = strings.ToLower(filename)
	successCreatingDirectory := false
	defer func() {
		if successCreatingDirectory {
			UISVGPaths = append(UISVGPaths, filepath.Join(UI, dirname))
		}
	}()
	if metadata["svgColor"] != "" {
		path := filepath.Join(UI, dirname, "color")
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true

		x, ok := metadata["svgColor"].(string)
		if ok {
			hash := md5.Sum([]byte(x))
			hashString := hex.EncodeToString(hash[:])
			pathsvg := hashCheckSVG[hashString]
			if pathsvg != "" { // the image has already been loaded, point the component to that path
				metadata["svgColor"] = pathsvg
				goto White
			}
			f, err := os.Create(filepath.Join(path, filename+"-color.svg"))
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = f.WriteString(x)
			if err != nil {
				fmt.Println(err)
				return
			}
			metadata["svgColor"] = getRelativePathForAPI(filepath.Join(dirname, "color", filename+"-color.svg")) //Replace the actual SVG with path to SVG
			writeHashCheckSVG(hashString, metadata["svgColor"].(string))
		}
	}
White:
	if metadata["svgWhite"] != "" {
		path := filepath.Join(UI, dirname, "white")
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true

		x, ok := metadata["svgWhite"].(string)
		if ok {
			hash := md5.Sum([]byte(x))
			hashString := hex.EncodeToString(hash[:])
			pathsvg := hashCheckSVG[hashString]
			if pathsvg != "" { // the image has already been loaded, point the component to that path
				metadata["svgWhite"] = pathsvg
				goto Complete
			}
			f, err := os.Create(filepath.Join(path, filename+"-white.svg"))
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = f.WriteString(x)
			if err != nil {
				fmt.Println(err)
				return
			}
			metadata["svgWhite"] = getRelativePathForAPI(filepath.Join(dirname, "white", filename+"-white.svg")) //Replace the actual SVG with path to SVG
			writeHashCheckSVG(hashString, metadata["svgWhite"].(string))
		}
	}
Complete:
	if metadata["svgComplete"] != "" {
		path := filepath.Join(UI, dirname, "complete")
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true

		x, ok := metadata["svgComplete"].(string)
		if ok {
			hash := md5.Sum([]byte(x))
			hashString := hex.EncodeToString(hash[:])
			pathsvg := hashCheckSVG[hashString]
			if pathsvg != "" { // the image has already been loaded, point the component to that path
				metadata["svgComplete"] = pathsvg
				return
			}
			f, err := os.Create(filepath.Join(path, filename+"-complete.svg"))
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = f.WriteString(x)
			if err != nil {
				fmt.Println(err)
				return
			}
			metadata["svgComplete"] = getRelativePathForAPI(filepath.Join(dirname, "complete", filename+"-complete.svg")) //Replace the actual SVG with path to SVG
			writeHashCheckSVG(hashString, metadata["svgComplete"].(string))
		}
	}
}
func WriteSVGsOnFileSystem(comp *v1alpha1.ComponentDefinition) {
	if comp.Metadata == nil {
		comp.Metadata = make(map[string]interface{})
	}
	if comp.Model.Metadata == nil {
		comp.Model.Metadata = make(map[string]interface{})
	}
	writeSVGHelper(comp.Metadata, comp.Model.Name, comp.Kind)             //Write SVG on components
	writeSVGHelper(comp.Model.Metadata, comp.Model.Name, comp.Model.Name) //Write SVG on models
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

func SliceContains(elements []string, name string) bool {
	for _, ele := range elements {
		if strings.Contains(strings.ToLower(name), ele) {
			return true
		}
	}
	return false
}

func GetPlatform() string {
	// local means running Meshery Server locally
	platform := "local"

	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount"); err == nil &&
		os.Getenv("KUBERNETES_SERVICE_HOST") != "" &&
		os.Getenv("KUBERNETES_SERVICE_PORT") != "" {

		// kubernetes means running Meshery Server inside a Kubernetes cluster
		platform = "kubernetes"
	}

	if isRunningInContainer() {
		// docker means running Meshery Server inside a Docker container
		platform = "docker"
	}

	return platform
}

// isRunningInContainer returns true if the process is running inside a container
// this code is taken from https://github.com/moby/libnetwork/blob/master/drivers/bridge/setup_bridgenetfiltering.go
func isRunningInContainer() bool {
	_, err := os.Stat("/.dockerenv")
	return !os.IsNotExist(err)
}

func SanitizeFileName(fileName string) string {
	extensionIndex := strings.LastIndex(fileName, ".")
	tempPath := strings.Split(fileName, "")

	finalPath := tempPath[:extensionIndex]
	suffixPath := strings.Join(tempPath[(extensionIndex+1):len(fileName)], "")
	finalPath = append(finalPath, "-*.", suffixPath)
	return strings.Join(finalPath, "")
}

func GetComponentFieldPathFromK8sFieldPath(path string) (newpath string) {
	if strings.HasPrefix(path, "metadata.") {
		path = strings.TrimPrefix(path, "metadata.")
		paths := strings.Split(path, ".")
		if len(paths) != 0 {
			if paths[0] == "name" || paths[0] == "namespace" || paths[0] == "labels" || paths[0] == "annotations" {
				return paths[0]
			}
		}
		return
	}
	return fmt.Sprintf("%s.%s", "settings", path)
}

// Prunes the diff part present in the k8s response message.
// Diff corresponds to the previous change and applied change, and doesn't contain any info which can be helpful to the user.
// If we want we can show this in a CodeEditor component.
func FormatK8sMessage(message string) string {
	exp, err := regexp.Compile(`(/?[a-zA-Z]).*\n([-,+])+`)
	if err != nil {
		return message
	}
	index := exp.FindStringIndex(message)
	if index == nil {
		return message
	}
	// If index is not nil, there will always be an array of length 2.
	// 0th index since we want the start index of matched string.
	return message[:index[0]]
}

func MarshalAndUnmarshal[k any, v any](val k) (unmarshalledvalue v, err error) {
	data, err := utils.Marshal(val)
	if err != nil {
		return
	}

	err = utils.Unmarshal(data, &unmarshalledvalue)
	if err != nil {
		return
	}
	return
}

func ReplaceSpacesAndConvertToLowercase(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", ""))
}