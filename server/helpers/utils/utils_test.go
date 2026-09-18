package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta3/component"
)

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		name string
		in   string
		sep  string
		want []string
	}{
		{
			name: "empty input returns nil",
			in:   "",
			sep:  ",",
			want: nil,
		},
		{
			name: "single value, no separator present",
			in:   "https://cloud.meshery.io",
			sep:  ",",
			want: []string{"https://cloud.meshery.io"},
		},
		{
			name: "two comma-separated values",
			in:   "https://cloud.meshery.io,https://cloud.acme.io",
			sep:  ",",
			want: []string{"https://cloud.meshery.io", "https://cloud.acme.io"},
		},
		{
			name: "trims whitespace around entries",
			in:   "  https://a.example  ,\thttps://b.example\n",
			sep:  ",",
			want: []string{"https://a.example", "https://b.example"},
		},
		{
			name: "drops empty entries from trailing or doubled separators",
			in:   "https://a.example,,https://b.example,",
			sep:  ",",
			want: []string{"https://a.example", "https://b.example"},
		},
		{
			name: "whitespace-only entries are dropped",
			in:   "https://a.example,   ,https://b.example",
			sep:  ",",
			want: []string{"https://a.example", "https://b.example"},
		},
		{
			name: "space-delimited ADAPTER_URLS form",
			in:   "meshery-istio:10000 meshery-linkerd:10001 meshery-consul:10002",
			sep:  ", \t\n\r",
			want: []string{"meshery-istio:10000", "meshery-linkerd:10001", "meshery-consul:10002"},
		},
		{
			name: "mixed comma and whitespace delimiters",
			in:   "a,b c\td\ne",
			sep:  ", \t\n\r",
			want: []string{"a", "b", "c", "d", "e"},
		},
		{
			name: "runs of mixed delimiters collapse",
			in:   "a , ,\tb",
			sep:  ", \t",
			want: []string{"a", "b"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SplitAndTrim(tc.in, tc.sep)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("SplitAndTrim(%q, %q) = %#v, want %#v", tc.in, tc.sep, got, tc.want)
			}
		})
	}
}

func TestValidateRegistrySVGPathComponent(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "real model name", value: "kubernetes", wantErr: false},
		{name: "real model name with hyphen", value: "cert-manager", wantErr: false},
		{name: "single character", value: "a", wantErr: false},
		{name: "digits and dots", value: "v1.2.3", wantErr: false},
		{name: "underscore", value: "my_model", wantErr: false},
		{name: "empty string is rejected", value: "", wantErr: true},
		{name: "bare traversal is rejected", value: "..", wantErr: true},
		{name: "issue's exact traversal payload is rejected", value: "../../../tmp/meshery-traversal", wantErr: true},
		{name: "embedded traversal segment is rejected", value: "istio/../../etc", wantErr: true},
		{name: "forward slash is rejected", value: "a/b", wantErr: true},
		{name: "backslash is rejected", value: `a\b`, wantErr: true},
		{name: "leading dot is rejected", value: ".hidden", wantErr: true},
		{name: "trailing hyphen is rejected", value: "trailing-", wantErr: true},
		{name: "leading hyphen is rejected", value: "-leading", wantErr: true},
		{name: "absolute path is rejected", value: "/etc/passwd", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateRegistrySVGPathComponent(tc.value)
			if (err != nil) != tc.wantErr {
				t.Errorf("validateRegistrySVGPathComponent(%q) error = %v, wantErr %v", tc.value, err, tc.wantErr)
			}
		})
	}
}

func TestValidateSVGContent(t *testing.T) {
	tests := []struct {
		name    string
		svg     string
		wantErr bool
	}{
		{name: "empty content is allowed", svg: "", wantErr: false},
		{name: "clean icon SVG is allowed", svg: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 2L2 7"/></svg>`, wantErr: false},
		{name: "issue's exact script payload is rejected", svg: `<svg xmlns="http://www.w3.org/2000/svg"><script>alert(document.domain)</script></svg>`, wantErr: true},
		{name: "script tag with attributes is rejected", svg: `<script type="text/javascript">evil()</script>`, wantErr: true},
		{name: "uppercase script tag is rejected", svg: `<SCRIPT>evil()</SCRIPT>`, wantErr: true},
		{name: "onload handler is rejected", svg: `<svg onload="evil()"></svg>`, wantErr: true},
		{name: "onerror handler is rejected", svg: `<svg><image href="x" onerror="evil()"/></svg>`, wantErr: true},
		{name: "foreignObject is rejected", svg: `<svg><foreignObject><body onload="evil()"/></foreignObject></svg>`, wantErr: true},
		{name: "javascript URI is rejected", svg: `<svg><a href="javascript:evil()">click</a></svg>`, wantErr: true},
		{name: "HTML-entity-encoded javascript URI is rejected", svg: `<svg><a href="&#106;avascript:evil()">click</a></svg>`, wantErr: true},
		{name: "hex-entity-encoded javascript URI is rejected", svg: `<svg><a href="&#x6a;avascript:evil()">click</a></svg>`, wantErr: true},
		{name: "mixed-case event handler is rejected", svg: `<svg OnLoad="evil()"></svg>`, wantErr: true},
		{name: "iframe is rejected", svg: `<svg><iframe src="https://evil.example"/></svg>`, wantErr: true},
		{name: "embed is rejected", svg: `<svg><embed src="https://evil.example"/></svg>`, wantErr: true},
		{name: "object is rejected", svg: `<svg><object data="https://evil.example"/></svg>`, wantErr: true},
		{name: "duplicate attribute is rejected as malformed rather than silently picking one value", svg: `<svg onload="safe()" onload="evil()"></svg>`, wantErr: true},
		{name: "unescaped ampersand is rejected as malformed rather than tolerated like a browser would", svg: `<svg><title>Fish & Chips</title></svg>`, wantErr: true},
		{name: "clean SVG with a style block is allowed", svg: `<svg xmlns="http://www.w3.org/2000/svg"><style>.a{fill:red}</style><path class="a" d="M12 2L2 7"/></svg>`, wantErr: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSVGContent(tc.svg)
			if (err != nil) != tc.wantErr {
				t.Errorf("validateSVGContent(%q) error = %v, wantErr %v", tc.svg, err, tc.wantErr)
			}
		})
	}
}

// withTempUI points the package-level UI var at a fresh temp directory for
// the duration of the test, restoring the original value on cleanup, so
// these tests never write through the real (relative) production path.
func withTempUI(t *testing.T) string {
	t.Helper()
	original := UI
	tempDir := t.TempDir()
	UI = tempDir
	t.Cleanup(func() { UI = original })
	return tempDir
}

// TestWriteSVGsOnFileSystem_RejectsExploitPayload reproduces the exact
// repro from the issue this fix addresses: a Model.Name of
// "../../../tmp/meshery-traversal" together with an SVG containing an
// inline <script>. Before the fix, this wrote an attacker-chosen file
// outside the meshmodels directory with attacker-chosen, unsanitized
// content. It must now be rejected outright, with nothing written to disk.
func TestWriteSVGsOnFileSystem_RejectsExploitPayload(t *testing.T) {
	tempDir := withTempUI(t)

	maliciousSVG := `<svg xmlns="http://www.w3.org/2000/svg"><script>alert(document.domain)</script></svg>`
	comp := &component.ComponentDefinition{
		Model: &model.ModelDefinition{
			Name: "../../../tmp/meshery-traversal",
		},
		Component: component.Component{
			Kind: "probe",
		},
		Styles: &core.ComponentStyles{
			SvgColor:    maliciousSVG,
			SvgWhite:    maliciousSVG,
			SvgComplete: maliciousSVG,
		},
	}

	err := WriteSVGsOnFileSystem(comp)
	if err == nil {
		t.Fatal("WriteSVGsOnFileSystem returned no error for a path-traversal Model.Name; the exploit was not blocked")
	}

	// Confirm nothing escaped: no file was written anywhere outside tempDir.
	// filepath.Join lexically cleans ".." segments, so a naive check against
	// the joined path would not have caught the original bug; walk the
	// actual parent of tempDir instead.
	parent := filepath.Dir(tempDir)
	entries, readErr := os.ReadDir(parent)
	if readErr != nil {
		t.Fatalf("reading %s: %v", parent, readErr)
	}
	for _, e := range entries {
		if strings.Contains(e.Name(), "meshery-traversal") {
			t.Errorf("found %q in %s: the traversal payload escaped the confined directory", e.Name(), parent)
		}
	}
}

func TestWriteSVGsOnFileSystem_RejectsScriptContentWithSafeName(t *testing.T) {
	withTempUI(t)

	comp := &component.ComponentDefinition{
		Model: &model.ModelDefinition{
			Name: "istio",
		},
		Component: component.Component{
			Kind: "probe",
		},
		Styles: &core.ComponentStyles{
			SvgColor: `<svg onload="alert(document.domain)"></svg>`,
		},
	}

	if err := WriteSVGsOnFileSystem(comp); err == nil {
		t.Fatal("WriteSVGsOnFileSystem returned no error for SVG content with an onload handler")
	}
}

// TestWriteSVGsOnFileSystem_AllowsLegitimateRegistration confirms the fix
// does not regress the real, working case: a plain slug name with clean
// SVG content still gets written to the expected location, matching what
// every current adapter and ArtifactHub-generated component sends.
func TestWriteSVGsOnFileSystem_AllowsLegitimateRegistration(t *testing.T) {
	tempDir := withTempUI(t)

	cleanSVG := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 2L2 7"/></svg>`
	comp := &component.ComponentDefinition{
		Model: &model.ModelDefinition{
			Name: "cert-manager",
		},
		Component: component.Component{
			Kind: "Probe",
		},
		Styles: &core.ComponentStyles{
			SvgColor: cleanSVG,
		},
	}

	if err := WriteSVGsOnFileSystem(comp); err != nil {
		t.Fatalf("WriteSVGsOnFileSystem returned an error for a legitimate registration: %v", err)
	}

	written := filepath.Join(tempDir, "cert-manager", "color", "probe-color.svg")
	data, err := os.ReadFile(written)
	if err != nil {
		t.Fatalf("expected file %s to exist: %v", written, err)
	}
	if string(data) != cleanSVG {
		t.Errorf("written content = %q, want %q", string(data), cleanSVG)
	}
}
