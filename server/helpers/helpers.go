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

// kubeconfigCredentialFileFields maps each kubeconfig field that names a file on
// disk to the inline field that replaces it. These three are the only fields in
// the kubeconfig schema with a base64 "-data" counterpart, so they are the only
// ones FlattenMinifyKubeConfig may resolve against the server's filesystem.
var kubeconfigCredentialFileFields = map[string]string{
	"certificate-authority": "certificate-authority-data",
	"client-certificate":    "client-certificate-data",
	"client-key":            "client-key-data",
}

// maxInlinedCredentialSize caps a single inlined credential. Real CA bundles and
// client certificates are a few KB; the cap keeps a path naming a huge file from
// being read into memory in full.
const maxInlinedCredentialSize = 1 << 20 // 1 MiB

// FlattenMinifyKubeConfig inlines a kubeconfig's file-path credential references
// into their base64 "-data" form, so a config written by a tool that points at
// local files (minikube writes `client-certificate: /home/u/.minikube/client.crt`)
// becomes self-contained.
//
// The kubeconfig walked here is USER-SUPPLIED: it arrives as the `k8sfile`
// multipart upload on POST /api/system/kubernetes and
// POST /api/system/kubernetes/contexts. Path resolution is therefore confined to
// the three fields in kubeconfigCredentialFileFields. Resolving *any* value that
// merely looked like an existing path - the previous behaviour - made this an
// arbitrary-file-read primitive against the server: a crafted kubeconfig could
// name /var/run/secrets/kubernetes.io/serviceaccount/token (or any other file the
// server process can read) under an arbitrary key, and the inlined bytes came
// back out in the contexts those handlers return and in the credential persisted
// to the remote provider. Any new field added here must be one a kubeconfig
// legitimately points at a file with.
func FlattenMinifyKubeConfig(config []byte) ([]byte, error) {
	cfg := map[interface{}]interface{}{}

	if err := yaml.Unmarshal(config, &cfg); err != nil {
		return config, err
	}

	NestedMapExplorer(cfg, func(key interface{}, value interface{}) (interface{}, interface{}) {
		// Keys reach this callback as the map key for a scalar field and as the
		// slice index for a bare list element; only a named field can be one of
		// the credential fields.
		keyStr, ok := key.(string)
		if !ok {
			return key, value
		}

		dataKey, isCredentialFile := kubeconfigCredentialFileFields[keyStr]
		if !isCredentialFile {
			return key, value
		}

		strV, ok := value.(string)
		if !ok || strV == "" {
			return key, value
		}

		data, err := readInlinableCredential(strV)
		if err != nil {
			// A reference that cannot be read is left untouched rather than
			// failing the whole config: the kubeconfig may be destined for a
			// host where it does resolve, and the context's own connection
			// check reports the failure with far better context than aborting
			// the import here would.
			return key, value
		}

		return dataKey, base64.StdEncoding.EncodeToString([]byte(data))
	})

	return yaml.Marshal(cfg)
}

// readInlinableCredential reads the credential file a kubeconfig field points at.
//
// It accepts only an absolute-looking path to a regular file under
// maxInlinedCredentialSize. os.Stat alone is satisfied by a directory and by a
// device node such as /dev/zero, whose read never returns; and a path without a
// separator would otherwise be resolved relative to the server's working
// directory rather than the kubeconfig's own.
func readInlinableCredential(path string) (string, error) {
	if !strings.Contains(path, string(filepath.Separator)) {
		return "", fmt.Errorf("%q is not a file path", path)
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%q is not a regular file", path)
	}
	if info.Size() > maxInlinedCredentialSize {
		return "", fmt.Errorf("%q is %d bytes, over the %d byte limit for an inlined kubeconfig credential", path, info.Size(), maxInlinedCredentialSize)
	}

	return ResolveFSRef(path)
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
