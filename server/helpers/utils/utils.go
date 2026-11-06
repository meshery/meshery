package utils

import (
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/meshery/meshkit/encoding"
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

// UI is a var, not a const, so tests can point it at a temp directory
// instead of writing through the real (relative) path.
var UI = "../../ui/public/static/img/meshmodels" //Relative to cmd/main.go
var UISVGPaths = make([]string, 1)

// registrySVGPathComponentPattern is the only shape a caller-supplied value
// (Model.Name, Component.Kind) may take before it is used to build a
// filesystem path under UI. Real registrants only ever send simple slugs
// (istio, kubernetes, cert-manager); anything else, in particular path
// separators and "..", is rejected outright rather than stripped, since a
// stripped-but-still-attacker-influenced value gives no real guarantee.
var registrySVGPathComponentPattern = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9._-]{0,254}[a-zA-Z0-9])?$`)

// dangerousSVGElements are element local names (case-insensitive, namespace
// prefix stripped) that must never appear in a caller-supplied SVG asset:
// script and foreignObject can carry or embed arbitrary script; iframe,
// embed, and object can load and execute another document entirely.
var dangerousSVGElements = map[string]bool{
	"script":        true,
	"foreignobject": true,
	"iframe":        true,
	"embed":         true,
	"object":        true,
}

// dangerousSVGAttrValuePattern matches a javascript: (or vbscript:) URI once
// an attribute value has been whitespace-collapsed, catching the common
// obfuscation of splitting the scheme with tabs/newlines that some legacy
// URI parsers tolerated.
var dangerousSVGAttrValuePattern = regexp.MustCompile(`(?i)^\s*(javascript|vbscript)\s*:`)

// containsDangerousSVGContent structurally walks svg as XML rather than
// pattern-matching the raw bytes, so a script or event handler split across
// unusual whitespace, attribute order, or HTML-entity-encoded characters
// (for example &#106;avascript:) cannot slip past a single regex the way it
// could with the string this replaced. html.UnescapeString normalizes
// entities in every attribute value before it is checked. If the content
// does not even parse as well-formed XML, it is rejected outright: real SVG
// icon assets are well-formed, and a browser's own SVG/HTML parser is far
// more forgiving of malformed markup than encoding/xml, so tolerating
// malformed input here would open exactly the obfuscation gap this
// function exists to close.
func containsDangerousSVGContent(svg string) bool {
	dec := xml.NewDecoder(strings.NewReader(svg))
	for {
		tok, err := dec.Token()
		if err != nil {
			return err != io.EOF
		}
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		if dangerousSVGElements[strings.ToLower(start.Name.Local)] {
			return true
		}
		for _, attr := range start.Attr {
			name := strings.ToLower(attr.Name.Local)
			if strings.HasPrefix(name, "on") {
				return true
			}
			value := html.UnescapeString(attr.Value)
			if dangerousSVGAttrValuePattern.MatchString(value) {
				return true
			}
		}
	}
}

// InvalidRegistrySVGAssetError signals that a caller-supplied Model.Name,
// Component.Kind, or SVG content was rejected by validation, as opposed to
// a filesystem or I/O failure encountered while writing an already-valid
// asset. A distinct type lets a caller of WriteSVGsOnFileSystem map it to
// HTTP 400 while treating any other failure as HTTP 500, so a disk-full or
// permission error is never reported back as if the request were bad.
type InvalidRegistrySVGAssetError string

func (e InvalidRegistrySVGAssetError) Error() string { return string(e) }

func validateRegistrySVGPathComponent(value string) error {
	if !registrySVGPathComponentPattern.MatchString(value) {
		return InvalidRegistrySVGAssetError(fmt.Sprintf("invalid registry path component %q: must be a plain alphanumeric slug (letters, digits, '.', '_', '-')", value))
	}
	return nil
}

func validateSVGContent(svg string) error {
	if svg == "" {
		return nil
	}
	if containsDangerousSVGContent(svg) {
		return InvalidRegistrySVGAssetError("SVG content contains a script, event handler, foreignObject, or javascript: URI, which is not permitted")
	}
	return nil
}

func writeSVGHelper(svgColor, svgWhite, svgComplete string, dirname, filename string) (svgColorPath, svgWhitePath, svgCompletePath string, err error) {
	if err = validateRegistrySVGPathComponent(dirname); err != nil {
		return "", "", "", err
	}
	if err = validateRegistrySVGPathComponent(filename); err != nil {
		return "", "", "", err
	}

	filename = strings.ToLower(filename)

	// UI itself is expected to already exist (it ships with the ui/
	// static assets), but OpenRoot requires that, unlike the MkdirAll
	// this replaces; create it first if some deployment mode doesn't
	// have it, so this isn't a new failure mode versus before.
	if _, statErr := os.Stat(UI); os.IsNotExist(statErr) {
		if mkErr := os.MkdirAll(UI, 0777); mkErr != nil {
			return "", "", "", mkErr
		}
	}

	// os.Root confines every MkdirAll/Create below to UI, following the
	// same pattern SafeOpenFile (server/handlers/utils.go) already uses
	// for reads: a dirname that somehow still resolved outside UI (a bug
	// in validateRegistrySVGPathComponent, or a future caller that skips
	// it) fails here instead of writing outside the intended directory,
	// and unlike a check-then-write, there is no window between
	// validation and the write in which the target could be swapped for
	// a symlink.
	root, err := os.OpenRoot(UI)
	if err != nil {
		return "", "", "", err
	}
	defer func() { _ = root.Close() }()

	successCreatingDirectory := false
	defer func() {
		if successCreatingDirectory {
			UISVGPaths = append(UISVGPaths, filepath.Join(UI, dirname))
		}
	}()
	if svgColor != "" {
		if err = validateSVGContent(svgColor); err != nil {
			return "", "", "", err
		}
		path := filepath.Join(dirname, "color")
		if err = root.MkdirAll(path, 0777); err != nil {
			fmt.Println(err)
			return "", "", "", err
		}
		successCreatingDirectory = true

		f, ferr := root.Create(filepath.Join(path, filename+"-color.svg"))
		if ferr != nil {
			fmt.Println(ferr)
			return "", "", "", ferr
		}
		_, ferr = f.WriteString(svgColor)
		if ferr == nil {
			ferr = f.Close()
		} else {
			_ = f.Close()
		}
		if ferr != nil {
			fmt.Println(ferr)
			return "", "", "", ferr
		}
		svgColorPath = getRelativePathForAPI(filepath.Join(dirname, "color", filename+"-color.svg")) //Replace the actual SVG with path to SVG

	}

	if svgWhite != "" {
		if err = validateSVGContent(svgWhite); err != nil {
			return "", "", "", err
		}
		path := filepath.Join(dirname, "white")
		if err = root.MkdirAll(path, 0777); err != nil {
			fmt.Println(err)
			return "", "", "", err
		}
		successCreatingDirectory = true

		f, ferr := root.Create(filepath.Join(path, filename+"-white.svg"))
		if ferr != nil {
			fmt.Println(ferr)
			return "", "", "", ferr
		}
		_, ferr = f.WriteString(svgWhite)
		if ferr == nil {
			ferr = f.Close()
		} else {
			_ = f.Close()
		}
		if ferr != nil {
			fmt.Println(ferr)
			return "", "", "", ferr
		}
		svgWhitePath = getRelativePathForAPI(filepath.Join(dirname, "white", filename+"-white.svg")) //Replace the actual SVG with path to SVG

	}
	if svgComplete != "" {
		if err = validateSVGContent(svgComplete); err != nil {
			return "", "", "", err
		}
		path := filepath.Join(dirname, "complete")
		if err = root.MkdirAll(path, 0777); err != nil {
			fmt.Println(err)
			return "", "", "", err
		}
		successCreatingDirectory = true

		f, ferr := root.Create(filepath.Join(path, filename+"-complete.svg"))
		if ferr != nil {
			fmt.Println(ferr)
			return "", "", "", ferr
		}
		_, ferr = f.WriteString(svgComplete)
		if ferr == nil {
			ferr = f.Close()
		} else {
			_ = f.Close()
		}
		if ferr != nil {
			fmt.Println(ferr)
			return "", "", "", ferr
		}
		svgCompletePath = getRelativePathForAPI(filepath.Join(dirname, "complete", filename+"-complete.svg")) //Replace the actual SVG with path to SVG

	}
	return svgColorPath, svgWhitePath, svgCompletePath, nil
}

// WriteSVGsOnFileSystem returns an InvalidRegistrySVGAssetError when
// comp.Model.Name, comp.Component.Kind, or an SVG field was rejected by
// validation, and a plain error for any other failure (a filesystem or I/O
// problem writing an already-valid asset). Callers should check for
// InvalidRegistrySVGAssetError to tell the two apart: the former is the
// caller's fault (HTTP 400), the latter is the server's (HTTP 500).
func WriteSVGsOnFileSystem(comp *component.ComponentDefinition) error {

	if comp.Styles != nil {
		_, _, _, err := writeSVGHelper(
			comp.Styles.SvgColor,
			comp.Styles.SvgWhite,
			comp.Styles.SvgComplete,
			comp.Model.Name,
			comp.Component.Kind) //Write SVG on components
		if err != nil {
			return err
		}
	}

	if comp.Model.Metadata != nil {
		svgComplete := ""
		if comp.Model.Metadata.SvgComplete != nil {
			svgComplete = *comp.Model.Metadata.SvgComplete
		}
		_, _, _, err := writeSVGHelper(
			comp.Model.Metadata.SvgColor,
			comp.Model.Metadata.SvgWhite,
			svgComplete,
			comp.Model.Name,
			comp.Model.Name) //Write SVG on models
		if err != nil {
			return err
		}
	}
	return nil
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

// FormatK8sMessage prunes the diff part present in the k8s response message.
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

func WriteEscaped(w http.ResponseWriter, data []byte, contentType string) (int, error) {
	if contentType == "" {
		contentType = "text/plain; charset=utf-8"
	}
	w.Header().Set("Content-Type", contentType)
	escaped := template.HTMLEscapeString(string(data))
	return w.Write([]byte(escaped))
}
