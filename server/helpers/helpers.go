package helpers

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/layer5io/meshkit/utils"
	"gopkg.in/yaml.v2"
)

func MergeStringMaps(maps ...map[string]string) map[string]string {
	res := map[string]string{}

	for _, mp := range maps {
		for k, v := range mp {
			res[k] = v
		}
	}

	return res
}

func ResolveFSRef(path string) (string, error) {
	return utils.ReadFileSource(fmt.Sprintf("file://%s", path))
}

func FlattenMinifyKubeConfig(config []byte) ([]byte, error) {
	cfg := map[interface{}]interface{}{}

	if err := yaml.Unmarshal(config, &cfg); err != nil {
		return config, err
	}

	// Find all references to file path and try to resolve them
	NestedMapExplorer(cfg, func(key interface{}, value interface{}) (interface{}, interface{}) {
		strV, ok := value.(string)
		if !ok {
			return key, value
		}

		// Check if the value is a filepath
		if !strings.Contains(strV, string(filepath.Separator)) {
			return key, value
		}

		// Assume that the value is a filepath, verify that with stat
		if _, err := os.Stat(strV); err != nil {
			return key, value
		}

		// Valid filepath found => get the file data
		data, err := ResolveFSRef(strV)
		if err != nil {
			return key, value
		}

		// Encode data as base64
		return fmt.Sprintf("%s-data", key), base64.StdEncoding.EncodeToString([]byte(data))
	})

	return yaml.Marshal(cfg)
}

func NestedMapExplorer(
	mp map[interface{}]interface{},
	fn func(key interface{}, value interface{}) (interface{}, interface{}),
) {
	for k, v := range mp {
		switch cNode := v.(type) {
		case map[interface{}]interface{}:
			NestedMapExplorer(cNode, fn)
		case []interface{}:
			for i, el := range cNode {
				switch ccNode := el.(type) {
				case map[interface{}]interface{}:
					NestedMapExplorer(ccNode, fn)
				default:
					_, nv := fn(i, el)
					cNode[i] = nv
				}
			}
		default:
			delete(mp, k)
			key, val := fn(k, cNode)
			mp[key] = val
		}
	}
}

func HostnameToPascalCase(input string) string {
    parts := strings.Split(input, ".")

    for i, part := range parts {
        parts[i] = CapitalizeFirstLetter(part)
    }

    pascalCaseHostname := strings.Join(parts, " ")

    return pascalCaseHostname
}

func CapitalizeFirstLetter(word string) string {
    if word == "" {
        return ""
    }

    runes := []rune(word)
    if unicode.IsLower(runes[0]) {
        runes[0] = unicode.ToUpper(runes[0])
    }

    return string(runes)
}