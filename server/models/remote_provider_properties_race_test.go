package models

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/spf13/viper"
)

func TestRemoteProviderProviderPropertiesConcurrentAccess(t *testing.T) {
	provider := newTestRemoteProvider(t, "https://remote.example.com")

	properties := []ProviderProperties{
		{
			ProviderName:        "Provider One",
			ProviderType:        RemoteProviderType,
			PackageVersion:      "1.0.0",
			ProviderURL:         "https://one.example.com",
			ProviderDescription: []string{"first provider"},
			Capabilities: Capabilities{
				{Feature: SyncPrefs, Endpoint: "/one/preferences"},
			},
		},
		{
			ProviderName:        "Provider Two",
			ProviderType:        RemoteProviderType,
			PackageVersion:      "2.0.0",
			ProviderURL:         "https://two.example.com",
			ProviderDescription: []string{"second provider"},
			Capabilities: Capabilities{
				{Feature: SyncPrefs, Endpoint: "/two/preferences"},
			},
		},
	}

	provider.SetProviderProperties(properties[0])

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 2000; i++ {
			provider.SetProviderProperties(properties[i%len(properties)])
		}
	}()

	for reader := 0; reader < 4; reader++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 2000; i++ {
				_ = provider.Name()
				_ = provider.GetProviderURL()
				_ = provider.GetProviderType()
				_ = provider.Description()
				_ = provider.PackageLocation()

				snapshot := provider.GetProviderProperties()
				_ = snapshot.ProviderName
				_ = snapshot.Capabilities.IsSupported(SyncPrefs)

				_ = provider.supportsFeature(SyncPrefs)
				_, _ = provider.endpointForFeature(SyncPrefs)

			}
		}()
	}

	wg.Wait()
}

func TestRemoteProviderProviderPropertiesDefensiveCopy(t *testing.T) {
	provider := newTestRemoteProvider(t, "https://remote.example.com")

	link, show, beta, external := true, true, false, true
	childLink := false

	input := ProviderProperties{
		ProviderName:        "Original",
		ProviderDescription: []string{"original description"},
		Capabilities:        Capabilities{{Feature: SyncPrefs, Endpoint: "/original"}},
		Redirects:           map[string]string{"login": "/original"},
		Extensions: Extensions{
			Navigator: NavigatorExtensions{
				{
					Title:  "Parent",
					Link:   &link,
					Show:   &show,
					IsBeta: &beta,
					Href:   Href{URI: "/parent", External: &external},
					Children: NavigatorExtensions{
						{Title: "Child", Link: &childLink},
					},
				},
			},
			Account: AccountExtensions{
				{
					Title: "Account parent",
					Children: AccountExtensions{
						{Title: "Account child"},
					},
				},
			},
			UserPrefs:    UserPrefsExtensions{{Component: "original"}},
			GraphQL:      GraphQLExtensions{{Component: "original"}},
			Collaborator: CollaboratorExtensions{{Component: "original"}},
		},
	}

	provider.SetProviderProperties(input)

	// Mutate the caller-owned input after the setter returns.
	input.ProviderDescription[0] = "mutated input"
	input.Capabilities[0].Endpoint = "/mutated"
	input.Redirects["login"] = "/mutated"
	input.Extensions.Navigator[0].Title = "Mutated parent"
	*input.Extensions.Navigator[0].Link = false
	*input.Extensions.Navigator[0].Href.External = false
	input.Extensions.Navigator[0].Children[0].Title = "Mutated child"
	input.Extensions.Account[0].Children[0].Title = "Mutated account child"
	input.Extensions.UserPrefs[0].Component = "mutated"
	input.Extensions.GraphQL[0].Component = "mutated"
	input.Extensions.Collaborator[0].Component = "mutated"

	got := provider.GetProviderProperties()
	if got.ProviderDescription[0] != "original description" {
		t.Fatalf("provider description changed through input: %q", got.ProviderDescription[0])
	}
	if got.Capabilities[0].Endpoint != "/original" {
		t.Fatalf("capability changed through input: %q", got.Capabilities[0].Endpoint)
	}
	if got.Redirects["login"] != "/original" {
		t.Fatalf("redirect changed through input: %q", got.Redirects["login"])
	}
	if got.Extensions.Navigator[0].Title != "Parent" || !*got.Extensions.Navigator[0].Link {
		t.Fatal("navigator parent or Link pointer changed through input")
	}
	if !*got.Extensions.Navigator[0].Href.External {
		t.Fatal("navigator Href.External changed through input")
	}
	if got.Extensions.Navigator[0].Children[0].Title != "Child" {
		t.Fatal("navigator child changed through input")
	}
	if got.Extensions.Account[0].Children[0].Title != "Account child" {
		t.Fatal("account child changed through input")
	}
	if got.Extensions.UserPrefs[0].Component != "original" ||
		got.Extensions.GraphQL[0].Component != "original" ||
		got.Extensions.Collaborator[0].Component != "original" {
		t.Fatal("extension slice changed through input")
	}

	// Mutate the returned snapshot and verify the provider remains unchanged.
	*got.Extensions.Navigator[0].Show = false
	*got.Extensions.Navigator[0].Children[0].Link = true
	got.Extensions.Account[0].Children[0].Title = "Mutated returned snapshot"
	got.ProviderDescription[0] = "mutated returned snapshot"
	got.Capabilities[0].Endpoint = "/mutated returned snapshot"
	got.Redirects["login"] = "/mutated returned snapshot"

	again := provider.GetProviderProperties()
	if again.ProviderDescription[0] != "original description" ||
		again.Capabilities[0].Endpoint != "/original" ||
		again.Redirects["login"] != "/original" {
		t.Fatal("provider properties changed through returned snapshot")
	}
	if !*again.Extensions.Navigator[0].Show ||
		*again.Extensions.Navigator[0].Children[0].Link {
		t.Fatal("navigator pointers changed through returned snapshot")
	}
	if again.Extensions.Account[0].Children[0].Title != "Account child" {
		t.Fatal("account child changed through returned snapshot")
	}
}

func TestCloneProviderPropertiesPreservesNilAndEmptySlices(t *testing.T) {
	properties := ProviderProperties{
		ProviderDescription: []string{},
		Capabilities:        Capabilities{},
		Extensions: Extensions{
			Navigator:    NavigatorExtensions{},
			UserPrefs:    UserPrefsExtensions{},
			GraphQL:      GraphQLExtensions{},
			Account:      AccountExtensions{},
			Collaborator: CollaboratorExtensions{},
		},
	}

	clone := cloneProviderProperties(properties)

	if clone.ProviderDescription == nil ||
		clone.Capabilities == nil ||
		clone.Extensions.Navigator == nil ||
		clone.Extensions.UserPrefs == nil ||
		clone.Extensions.GraphQL == nil ||
		clone.Extensions.Account == nil ||
		clone.Extensions.Collaborator == nil {
		t.Fatal("clone changed a non-nil empty slice to nil")
	}

	if cloneProviderProperties(ProviderProperties{}).ProviderDescription != nil {
		t.Fatal("clone changed a nil slice to non-nil")
	}
}

func TestRemoteProviderDescriptionPreservesNilAndEmptySlices(t *testing.T) {
	provider := newTestRemoteProvider(t, "https://remote.example.com")

	provider.SetProviderProperties(ProviderProperties{ProviderDescription: []string{}})
	description := provider.Description()
	if description == nil {
		t.Fatal("Description changed a non-nil empty slice to nil")
	}

	provider.SetProviderProperties(ProviderProperties{})
	if description = provider.Description(); description != nil {
		t.Fatal("Description changed a nil slice to non-nil")
	}
}

func makeProviderPackageArchive(name string) []byte {
	var archive bytes.Buffer
	gzipWriter := gzip.NewWriter(&archive)
	tarWriter := tar.NewWriter(gzipWriter)
	content := []byte(name)
	_ = tarWriter.WriteHeader(&tar.Header{
		Name: "extension.txt",
		Mode: 0600,
		Size: int64(len(content)),
	})
	_, _ = tarWriter.Write(content)
	_ = tarWriter.Close()
	_ = gzipWriter.Close()
	return archive.Bytes()
}

func TestRemoteProviderPackageSnapshotRemainsConsistent(t *testing.T) {
	provider := newTestRemoteProvider(t, "https://remote.example.com")
	homeDir := t.TempDir()
	originalHomeDir := remoteProviderHomeDir
	remoteProviderHomeDir = homeDir
	t.Cleanup(func() {
		remoteProviderHomeDir = originalHomeDir
	})
	origSkip := viper.Get(SKIP_DOWNLOAD_EXTENSIONS_ENV)
	defer viper.Set(SKIP_DOWNLOAD_EXTENSIONS_ENV, origSkip)
	viper.Set(SKIP_DOWNLOAD_EXTENSIONS_ENV, false)

	var oneRequests atomic.Int64
	var twoRequests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/one.tar.gz":
			oneRequests.Add(1)
		case "/two.tar.gz":
			twoRequests.Add(1)
		default:
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/gzip")
		_, _ = w.Write(makeProviderPackageArchive(r.URL.Path))
	}))
	defer server.Close()

	properties := []ProviderProperties{
		{
			ProviderName:   "Provider One",
			ProviderURL:    "https://one.example.com",
			PackageVersion: "1.0.0",
			PackageURL:     server.URL + "/one.tar.gz",
		},
		{
			ProviderName:   "Provider Two",
			ProviderURL:    "https://two.example.com",
			PackageVersion: "2.0.0",
			PackageURL:     server.URL + "/two.tar.gz",
		},
	}

	for i, expected := range properties {
		provider.SetProviderProperties(expected)
		snapshot := provider.GetProviderProperties()
		location := snapshot.PackageLocation()
		snapshot.DownloadProviderExtensionPackage(provider.Log)
		content, err := os.ReadFile(filepath.Join(location, "extension.txt"))
		if err != nil {
			t.Fatalf("provider %d package downloaded to the wrong snapshot location: %v", i, err)
		}
		expectedContent := "/one.tar.gz"
		if i == 1 {
			expectedContent = "/two.tar.gz"
		}
		if string(content) != expectedContent {
			t.Fatalf("provider %d package content=%q, want %q", i, content, expectedContent)
		}
	}

	provider.SetProviderProperties(properties[0])
	packageRoot := filepath.Join(homeDir, ".meshery", "provider")
	if err := os.RemoveAll(packageRoot); err != nil {
		t.Fatalf("remove provider package cache before concurrent downloads: %v", err)
	}

	var wg sync.WaitGroup
	var packageDownloadMu sync.Mutex
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000; i++ {
			provider.SetProviderProperties(properties[i%len(properties)])
		}
	}()

	for reader := 0; reader < 4; reader++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 5000; i++ {
				snapshot := provider.GetProviderProperties()
				validOne := snapshot.ProviderName == "Provider One" &&
					snapshot.ProviderURL == "https://one.example.com" &&
					snapshot.PackageVersion == "1.0.0"
				validTwo := snapshot.ProviderName == "Provider Two" &&
					snapshot.ProviderURL == "https://two.example.com" &&
					snapshot.PackageVersion == "2.0.0"
				if !validOne && !validTwo {
					t.Errorf("package snapshot mixed provider fields: name=%q url=%q version=%q",
						snapshot.ProviderName, snapshot.ProviderURL, snapshot.PackageVersion)
					return
				}
				packageDownloadMu.Lock()
				location := snapshot.PackageLocation()
				snapshot.DownloadProviderExtensionPackage(provider.Log)
				content, err := os.ReadFile(filepath.Join(location, "extension.txt"))
				if err != nil {
					packageDownloadMu.Unlock()
					t.Errorf("package snapshot was not downloaded to %s: %v", location, err)
					return
				}
				expectedContent := "/one.tar.gz"
				if validTwo {
					expectedContent = "/two.tar.gz"
				}
				if string(content) != expectedContent {
					packageDownloadMu.Unlock()
					t.Errorf("package snapshot destination=%s content=%q, want %q",
						location, content, expectedContent)
					return
				}
				packageDownloadMu.Unlock()
			}
		}()
	}

	wg.Wait()

	for i, expected := range properties {
		location := filepath.Join(packageRoot, expected.ProviderName, expected.PackageVersion)
		content, err := os.ReadFile(filepath.Join(location, "extension.txt"))
		if err != nil {
			t.Fatalf("provider %d package snapshot was not downloaded to %s: %v", i, location, err)
		}
		expectedContent := "/one.tar.gz"
		if i == 1 {
			expectedContent = "/two.tar.gz"
		}
		if string(content) != expectedContent {
			t.Fatalf("provider %d package destination=%s content=%q, want %q", i, location, content, expectedContent)
		}
	}

	if oneRequests.Load() == 0 || twoRequests.Load() == 0 {
		t.Fatalf("package download did not observe both provider package URLs: one=%d two=%d",
			oneRequests.Load(), twoRequests.Load())
	}
}
