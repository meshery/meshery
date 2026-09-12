package utils

import (
	"encoding/json"
	"fmt"
	mathrand "math/rand"
	"strconv"
	"strings"

	"github.com/meshery/meshkit/encoding"
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

// GetRandomAlphabetsOfDigit generates a random string of a given length
// using lowercase alphabets.
func GetRandomAlphabetsOfDigit(length int) (s string) {
	charSet := "abcdedfghijklmnopqrstuvwxyz"
	for i := 0; i < length; i++ {
		random := mathrand.Intn(len(charSet))
		randomChar := charSet[random]
		s += string(randomChar)
	}
	return
}

// IsDesignInAlpha2Format reports whether a stored design predates the versioned
// design schema and therefore must be migrated by convertV1alpha2ToV1beta3.
//
// Every design schema from v1beta1 onward stamps a versioned
// "designs.meshery.io/<version>" schemaVersion (v1beta1, v1beta3, ...); the
// legacy v1alpha2 format carries none (it is identified instead by a top-level
// `services` map). Detection therefore keys off the canonical schemaVersion
// prefix rather than whitelisting a single "current" version.
//
// Whitelisting one version is what regressed imported-design rendering: the
// import pipeline (NewPatternFileFromK8sManifest) stamps v1beta3, but this
// guard only recognized v1beta1, so every freshly imported design failed the
// check, was misclassified as alpha2, and convertV1alpha2ToV1beta3 reparsed it
// as a v1alpha2.PatternFile (which has `services`, not `components`) -
// producing a design with zero components that rendered as an empty canvas.
// Matching the prefix keeps current and future design schema versions out of
// the lossy migration path.
func IsDesignInAlpha2Format(patternFile string) (bool, error) {
	// Decode only schemaVersion into a targeted struct instead of the whole
	// file into a map[string]interface{}, which would allocate the entire
	// nested design tree just to read one field. Both the json and yaml tags
	// are required: encoding.Unmarshal tries encoding/json first and falls back
	// to gopkg.in/yaml.v3, so a json-only tag would miss the camelCase
	// schemaVersion key in a YAML-encoded pattern file and misclassify it.
	var design struct {
		SchemaVersion string `json:"schemaVersion" yaml:"schemaVersion"`
	}
	if err := encoding.Unmarshal([]byte(patternFile), &design); err != nil {
		return true, err
	}

	if strings.HasPrefix(design.SchemaVersion, "designs.meshery.io/") {
		return false, nil
	}
	return true, nil
}
