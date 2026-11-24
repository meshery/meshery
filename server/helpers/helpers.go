package helpers

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/meshery/meshkit/utils"
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

// kubeConfigFileReferenceKeys are the kubeconfig fields that name a file on the
// host running the client instead of embedding its content.
var kubeConfigFileReferenceKeys = map[string]struct{}{
	"certificate-authority": {},
	"client-certificate":    {},
	"client-key":            {},
	"tokenFile":             {},
}

// RejectKubeConfigFileReferences returns ErrKubeConfigFileReference when config
// names a file instead of embedding its content (certificate-authority,
// client-certificate, client-key or tokenFile, at any depth, whether or not the
// file exists).
//
// Use it on every kubeconfig a user uploads. The Meshery Server host is not the
// uploader's machine, so such a path is either dead or an attempt to read the
// server's own files, and the result must not reveal which. Do not use it on
// the server's own kubeconfig, whose paths are legitimate and are resolved by
// FlattenMinifyKubeConfig.
func RejectKubeConfigFileReferences(config []byte) error {
	// Decode exactly the way the consumers do (models.InternalKubeConfig, which
	// skips unknown top-level keys) and also generically. The two decodes can
	// disagree on a crafted document, so a reference found by either rejects it.
	// A document the consumers' own decode refuses is passed through: nothing
	// downstream reads a file for it.
	var typed struct {
		Clusters []map[interface{}]interface{} `yaml:"clusters"`
		Users    []map[interface{}]interface{} `yaml:"users"`
	}
	if err := yaml.Unmarshal(config, &typed); err != nil {
		return nil
	}
	for _, c := range typed.Clusters {
		if hasKubeConfigFileReference(c) {
			return ErrKubeConfigFileReference()
		}
	}
	for _, u := range typed.Users {
		if hasKubeConfigFileReference(u) {
			return ErrKubeConfigFileReference()
		}
	}
	var doc interface{}
	if err := yaml.Unmarshal(config, &doc); err == nil && hasKubeConfigFileReference(doc) {
		return ErrKubeConfigFileReference()
	}
	return nil
}

func hasKubeConfigFileReference(node interface{}) bool {
	switch n := node.(type) {
	case map[interface{}]interface{}:
		for k, v := range n {
			if key, ok := k.(string); ok {
				if _, isRef := kubeConfigFileReferenceKeys[key]; isRef && v != nil && v != "" {
					return true
				}
			}
			if hasKubeConfigFileReference(v) {
				return true
			}
		}
	case []interface{}:
		for _, el := range n {
			if hasKubeConfigFileReference(el) {
				return true
			}
		}
	}
	return false
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
