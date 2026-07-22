package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta3/component"
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
	RegistryLocation      = ".meshery/models"
	DefVersion            = "1.0.0"
)

// SplitAndTrim splits s on any rune that appears in delims, trims whitespace
// from each resulting field, and discards empty entries. Use this when reading
// delimited values from an environment variable via viper.GetString:
// viper.GetStringSlice does not split a single delimited env-var value into
// multiple slice entries when AutomaticEnv is enabled, so the whole string
// flows through as one element. Pass the full set of expected separator
// characters — e.g. ", \t\n\r" — to accept either comma-separated or
// whitespace-separated configurations, which both forms appear across the
// Meshery manifests in install/.
func SplitAndTrim(s, delims string) []string {
	if s == "" {
		return nil
	}
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return strings.ContainsRune(delims, r)
	})
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if f = strings.TrimSpace(f); f != "" {
			out = append(out, f)
		}
	}
	return out
}

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
	res, err := utils.MarshalAndUnmarshal[interface{}, map[string]interface{}](mp)
	if err != nil {
		return map[string]interface{}{}
	}
	return res
}

const UI = "../../ui/public/static/img/meshmodels" //Relative to cmd/main.go
var UISVGPaths = make([]string, 1)

func writeSVGHelper(svgColor, svgWhite, svgComplete string, dirname, filename string) (svgColorPath, svgWhitePath, svgCompletePath string) {
	filename = strings.ToLower(filename)
	successCreatingDirectory := false
	defer func() {
		if successCreatingDirectory {
			UISVGPaths = append(UISVGPaths, filepath.Join(UI, dirname))
		}
	}()

	variants := []struct {
		data   string
		suffix string
		out    *string
	}{
		{svgColor, "color", &svgColorPath},
		{svgWhite, "white", &svgWhitePath},
		{svgComplete, "complete", &svgCompletePath},
	}
	for _, variant := range variants {
		if variant.data == "" {
			continue
		}
		path := filepath.Join(UI, dirname, variant.suffix)
		if err := os.MkdirAll(path, 0777); err != nil {
			fmt.Println(err)
			return
		}
		successCreatingDirectory = true

		svgFileName := filename + "-" + variant.suffix + ".svg"
		f, err := os.Create(filepath.Join(path, svgFileName))
		if err != nil {
			fmt.Println(err)
			return
		}
		if _, err := f.WriteString(variant.data); err != nil {
			fmt.Println(err)
			return
		}
		*variant.out = getRelativePathForAPI(filepath.Join(dirname, variant.suffix, svgFileName)) //Replace the actual SVG with path to SVG
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
		if err := os.RemoveAll(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Println(err)
		}
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
	return fileName[:extensionIndex] + "-*" + fileName[extensionIndex:]
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
	// replaceIfSet overwrites *field with the SVG data found at *field's
	// current value (a path), leaving it untouched if empty.
	replaceIfSet := func(field *string) error {
		if *field == "" {
			return nil
		}
		svgData, err := readSVGData(*field)
		if err != nil {
			return err
		}
		*field = svgData
		return nil
	}

	// Replace SVG paths with actual data in metadata
	metadata := model.Metadata
	if err := replaceIfSet(&metadata.SvgColor); err != nil {
		return err
	}
	if err := replaceIfSet(&metadata.SvgWhite); err != nil {
		return err
	}

	components, ok := model.Components.([]component.ComponentDefinition)
	if !ok {
		return fmt.Errorf("invalid type for Components field")
	}
	// Replace SVG paths with actual data in components
	for i := range components {
		compStyle := components[i].Styles
		if compStyle != nil {
			if err := replaceIfSet(&compStyle.SvgColor); err != nil {
				return err
			}
			if err := replaceIfSet(&compStyle.SvgWhite); err != nil {
				return err
			}
		}
		components[i].Styles = compStyle
	}
	model.Components = components
	return nil
}
func CreateVersionedDirectoryForModelAndComp(version, modelName string) (string, string, error) {
	modelLocation := filepath.Join(os.Getenv("HOME"), RegistryLocation)
	modelDirPath := filepath.Join(modelLocation, modelName, version, DefVersion)
	err := utils.CreateDirectory(modelDirPath)
	if err != nil {
		return "", "", err
	}

	compDirPath := filepath.Join(modelDirPath, "components")
	err = utils.CreateDirectory(compDirPath)
	return modelDirPath, compDirPath, err
}
func CopyDirectory(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		info, err := entry.Info()
		if err != nil {
			return err
		}

		if info.IsDir() {
			if err := os.MkdirAll(dstPath, info.Mode()); err != nil {
				return err
			}
			if err := CopyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := srcFile.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dstFile.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Ensure the destination file has the same permissions as the source
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}
