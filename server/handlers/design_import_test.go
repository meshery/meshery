package handlers

import (
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

func TestFetchFileFromURL_RejectsSSRF(t *testing.T) {
	restrictedURLs := []struct {
		name string
		url  string
	}{
		{"IPv4 Loopback", "http://127.0.0.1:9081/api"},
		{"Localhost hostname", "http://localhost:9081/api"},
		{"IPv6 Loopback", "http://[::1]:9081/api"},
		{"AWS/GCP IMDS", "http://169.254.169.254/latest/meta-data/"},
		{"IPv4-mapped IPv6 IMDS", "http://[::ffff:169.254.169.254]/latest/meta-data/"},
		{"GCP metadata hostname", "http://metadata.google.internal/computeMetadata/v1/"},
		{"AWS instance-data hostname", "http://instance-data/latest/meta-data/"},
		{"Private Class A", "http://10.0.0.1/design.yaml"},
		{"Private Class B", "http://172.16.0.1/design.yaml"},
		{"Private Class C", "http://192.168.1.1/design.yaml"},
		{"File scheme", "file:///etc/passwd"},
		{"FTP scheme", "ftp://example.com/design.yaml"},
		{"Empty host", "http:///design.yaml"},
	}

	for _, tc := range restrictedURLs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fetchFileFromURL(tc.url)
			if err == nil {
				t.Fatalf("fetchFileFromURL(%q) should have failed with SSRF error, but succeeded", tc.url)
			}

			// Ensure error code is models.ErrDoRequestCode
			code := errors.GetCode(err)
			if code != models.ErrDoRequestCode {
				t.Errorf("expected error code %q, got %q (err: %v)", models.ErrDoRequestCode, code, err)
			}
		})
	}
}

func TestResolveImportVariant_URLImportSSRF(t *testing.T) {
	const restrictedURL = "http://169.254.169.254/latest/meta-data/"
	designName := "Malicious Design"

	body := pattern.MesheryPatternImportRequestBody{}
	if err := body.FromMesheryPatternImportURLPayload(pattern.MesheryPatternImportURLPayload{
		Name: &designName,
		Url:  restrictedURL,
	}); err != nil {
		t.Fatalf("failed to construct import request body: %v", err)
	}

	variant, err := resolveImportVariant(body)
	if err == nil {
		t.Fatal("expected resolveImportVariant to reject restricted URL, but it succeeded")
	}

	// Verify the variant still carries the URL so the downstream event logs the source
	if variant.File.FileName != restrictedURL {
		t.Errorf("expected variant.File.FileName to preserve %q, got %q", restrictedURL, variant.File.FileName)
	}
	if variant.Name != designName {
		t.Errorf("expected variant.Name to preserve %q, got %q", designName, variant.Name)
	}

	if !strings.Contains(err.Error(), "restricted") && !strings.Contains(err.Error(), "forbidden") && !strings.Contains(err.Error(), "unable to make") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestConvertFileToDesignRejectsUnregisteredResources(t *testing.T) {
	rm, _ := newTestRegistryDB(t)
	_, _, err := ConvertFileToDesign(FileToImport{
		FileName: "deployment.yaml",
		Data:     []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: example\n"),
	}, rm, newTestLogger(t))
	if err == nil || !strings.Contains(err.Error(), "Deployment") || !strings.Contains(err.Error(), "apps/v1") {
		t.Fatalf("expected an error identifying the unregistered resource, got %v", err)
	}
}

func TestConvertFileToDesignRejectsPartiallyRegisteredManifest(t *testing.T) {
	rm, _ := newTestRegistryDB(t)
	seedTestComponent(t, rm) // Registers batch/v1 Job.
	log := newTestLogger(t)
	registered := "apiVersion: batch/v1\nkind: Job\nmetadata:\n  name: example-job\n"

	design, _, err := ConvertFileToDesign(FileToImport{
		FileName: "job.yaml",
		Data:     []byte(registered),
	}, rm, log)
	if err != nil || len(design.Components) != 1 {
		t.Fatalf("expected registered Job to import as one component, got %d components and error %v", len(design.Components), err)
	}

	_, _, err = ConvertFileToDesign(FileToImport{
		FileName: "mixed.yaml",
		Data:     []byte(registered + "---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: example-deployment\n"),
	}, rm, log)
	if err == nil || !strings.Contains(err.Error(), "Deployment") || !strings.Contains(err.Error(), "apps/v1") {
		t.Fatalf("expected an error identifying the unregistered Deployment, got %v", err)
	}
}
