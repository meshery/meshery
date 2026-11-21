package models

import (
	"sync"
	"testing"
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

func TestRemoteProviderPackageSnapshotRemainsConsistent(t *testing.T) {
	provider := newTestRemoteProvider(t, "https://remote.example.com")

	properties := []ProviderProperties{
		{
			ProviderName:   "Provider One",
			ProviderURL:    "https://one.example.com",
			PackageVersion: "1.0.0",
		},
		{
			ProviderName:   "Provider Two",
			ProviderURL:    "https://two.example.com",
			PackageVersion: "2.0.0",
		},
	}

	provider.SetProviderProperties(properties[0])

	var wg sync.WaitGroup
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
			}
		}()
	}

	wg.Wait()
}
