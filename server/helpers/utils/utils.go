package utils

import (
	"crypto/md5"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

const (
	HelmChartURL          = "https://meshery.io/charts/"
	HelmChartName         = "meshery"
	HelmChartOperatorName = "meshery-operator"
	MesheryFolder         = ".meshery"
	ManifestsFolder       = "manifests"
	registryLocation      = ".meshery/models"
	DefVersion            = "1.0.0"
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

func IsClosed[K any](ch chan K) bool {
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
func writeSVGHelper(svgColor, svgWhite, svgComplete string, dirname, filename string) (svgColorPath, svgWhitePath, svgCompletePath string) {
	filename = strings.ToLower(filename)
	successCreatingDirectory := false
	defer func() {
		if successCreatingDirectory {
			UISVGPaths = append(UISVGPaths, filepath.Join(UI, dirname))
		}
	}()
	if svgColor != "" {
		path := filepath.Join(UI, dirname, "color")
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true

		hash := md5.Sum([]byte(svgColor))
		hashString := hex.EncodeToString(hash[:])
		pathsvg := hashCheckSVG[hashString]
		if pathsvg != "" { // the image has already been loaded, point the component to that path
			svgColorPath = pathsvg
			goto White
		}
		f, err := os.Create(filepath.Join(path, filename+"-color.svg"))
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = f.WriteString(svgColor)
		if err != nil {
			fmt.Println(err)
			return
		}
		svgColorPath = getRelativePathForAPI(filepath.Join(dirname, "color", filename+"-color.svg")) //Replace the actual SVG with path to SVG
		writeHashCheckSVG(hashString, svgColor)

	}
White:
	if svgWhite != "" {
		path := filepath.Join(UI, dirname, "white")
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true

		hash := md5.Sum([]byte(svgWhite))
		hashString := hex.EncodeToString(hash[:])
		pathsvg := hashCheckSVG[hashString]
		if pathsvg != "" { // the image has already been loaded, point the component to that path
			svgWhitePath = pathsvg
			goto Complete
		}
		f, err := os.Create(filepath.Join(path, filename+"-white.svg"))
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = f.WriteString(svgWhite)
		if err != nil {
			fmt.Println(err)
			return
		}
		svgWhitePath = getRelativePathForAPI(filepath.Join(dirname, "white", filename+"-white.svg")) //Replace the actual SVG with path to SVG
		writeHashCheckSVG(hashString, svgWhite)

	}
Complete:
	if svgComplete != "" {
		path := filepath.Join(UI, dirname, "complete")
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true

		hash := md5.Sum([]byte(svgComplete))
		hashString := hex.EncodeToString(hash[:])
		pathsvg := hashCheckSVG[hashString]
		if pathsvg != "" { // the image has already been loaded, point the component to that path
			svgCompletePath = pathsvg
			return
		}
		f, err := os.Create(filepath.Join(path, filename+"-complete.svg"))
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = f.WriteString(svgComplete)
		if err != nil {
			fmt.Println(err)
			return
		}
		svgCompletePath = getRelativePathForAPI(filepath.Join(dirname, "complete", filename+"-complete.svg")) //Replace the actual SVG with path to SVG
		writeHashCheckSVG(hashString, svgComplete)

	}
	return
}
func WriteSVGsOnFileSystem(comp *component.ComponentDefinition) {

	if comp.Styles != nil {
		writeSVGHelper(
			comp.Styles.SvgColor,
			comp.Styles.SvgWhite,
			comp.Styles.SvgComplete,
			comp.Model.Name,
			comp.Component.Kind) //Write SVG on components
	}

	if comp.Model.Metadata != nil {
		svgComplete := ""
		if comp.Model.Metadata.SvgComplete != nil {
			svgComplete = *comp.Model.Metadata.SvgComplete
		}
		writeSVGHelper(
			comp.Model.Metadata.SvgColor,
			comp.Model.Metadata.SvgWhite,
			svgComplete,
			comp.Model.Name,
			comp.Model.Name) //Write SVG on models
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

	err = encoding.Unmarshal([]byte(data), &unmarshalledvalue)
	if err != nil {
		return
	}
	return
}

type JSONMap map[string]interface{}

// Value converts the JSON map to a database value.
func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan converts the database value to a JSON map.
func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// ApplyFilters applies dynamic filters to the GORM query
func ApplyFilters(query *gorm.DB, filter string, dynamicKeys []string) *gorm.DB {
	if filter != "" {
		filterArr := strings.Split(filter, " ")
		filterKey := filterArr[0]
		filterVal := strings.Join(filterArr[1:], " ")

		switch filterKey {
		case "deleted_at":
			// Handle deleted_at filter
			if filterVal == "Deleted" {
				query = query.Where("deleted_at IS NOT NULL")
			} else {
				query = query.Where("deleted_at IS NULL")
			}
		default:
			// Handle dynamic keys
			for _, key := range dynamicKeys {
				if filterKey == key {
					query = query.Where(fmt.Sprintf("%s = ?", filterKey), filterVal)
					break
				}
			}
		}
	}

	return query
}

func FormatToTitleCase(s string) string {
	c := cases.Title(language.English)
	return c.String(s)
}
func ExtractFile(filePath string, destDir string) error {

	if utils.IsTarGz(filePath) {
		return utils.ExtractTarGz(destDir, filePath)
	} else if utils.IsZip(filePath) {
		return utils.ExtractZip(destDir, filePath)
	}
	return utils.ErrExtractType
}
func ConvertToJSONCompatible(data interface{}) interface{} {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for key, value := range v {
			m[key.(string)] = ConvertToJSONCompatible(value)
		}
		return m
	case []interface{}:
		for i, item := range v {
			v[i] = ConvertToJSONCompatible(item)
		}
	}
	return data
}
func ReplaceSVGData(model *model.ModelDefinition) error {
	// Function to read SVG data from file
	readSVGData := func(path string) (string, error) {
		path = "../../" + path
		svgData, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(svgData), nil
	}

	// Replace SVG paths with actual data in metadata
	metadata := model.Metadata
	if metadata.SvgColor != "" {
		svgData, err := readSVGData(metadata.SvgColor)
		if err == nil {
			metadata.SvgColor = svgData
		} else {
			return err
		}
	}
	if metadata.SvgWhite != "" {
		svgData, err := readSVGData(metadata.SvgWhite)
		if err == nil {
			metadata.SvgWhite = svgData
		} else {
			return err
		}
	}
	components, ok := model.Components.([]component.ComponentDefinition)
	if !ok {
		return fmt.Errorf("invalid type for Components field")
	}
	// Replace SVG paths with actual data in components
	for i := range components {
		compStyle := components[i].Styles
		if compStyle != nil {
			svgColor, err := readSVGData(compStyle.SvgColor)
			if err == nil {
				compStyle.SvgColor = svgColor
			} else {
				return err
			}
			svgWhite, err := readSVGData(compStyle.SvgWhite)
			if err == nil {
				compStyle.SvgWhite = svgWhite
			} else {
				return err
			}
		}
		components[i].Styles = compStyle
	}
	model.Components = components
	return nil
}
func CreateVersionedDirectoryForModelAndComp(version, modelName string) (string, string, error) {
	modelLocation := filepath.Join(os.Getenv("HOME"), registryLocation)
	modelDirPath := filepath.Join(modelLocation, modelName, version, DefVersion)
	err := utils.CreateDirectory(modelDirPath)
	if err != nil {
		return "", "", err
	}

	compDirPath := filepath.Join(modelDirPath, "components")
	err = utils.CreateDirectory(compDirPath)
	return modelDirPath, compDirPath, err
}
