package oci

import (
	"encoding/json"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/opencontainers/go-digest"
	"oras.land/oras-go/v2"
	specs "github.com/opencontainers/image-spec/specs-go"
)

// generate OCI Descriptor for the given design
func generateDescriptor(mediaType string, blob []byte, target oras.Target) (desc ocispec.Descriptor) {
	desc = ocispec.Descriptor{ // Generate descriptor based on the media type and blob content
		MediaType: mediaType,
		Digest:    digest.FromBytes(blob), // Calculate digest
		Size:      int64(len(blob)),       // Include blob size
	}
	return desc
}

func generateManifestContent(config ocispec.Descriptor, layers ...ocispec.Descriptor) ([]byte, error) {
	content := ocispec.Manifest{
		Config:    config, // Set config blob
		Layers:    layers, // Set layer blobs
		Versioned: specs.Versioned{SchemaVersion: 2},
	}
	return json.Marshal(content) // Get json content
}

