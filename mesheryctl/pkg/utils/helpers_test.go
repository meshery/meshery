package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/constants"
	"github.com/meshery/schemas/models/core"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd/api"
)

type mockCloser struct {
	closeFunc func() error
}

func (m mockCloser) Close() error {
	return m.closeFunc()

}

type trackingTransport struct {
	base http.RoundTripper
	body *trackingBody
}

func (t *trackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	resp, err := base.RoundTrip(req)
	if err == nil && resp != nil && resp.Body != nil {
		t.body = &trackingBody{ReadCloser: resp.Body}
		resp.Body = t.body
	}
	return resp, err
}

type trackingBody struct {
	io.ReadCloser
	closed bool
}

func (b *trackingBody) Close() error {
	b.closed = true
	return b.ReadCloser.Close()
}

func installTrackingTransport(t *testing.T) *trackingTransport {
	t.Helper()
	transport := &trackingTransport{base: http.DefaultTransport}
	previous := http.DefaultTransport
	http.DefaultTransport = transport
	t.Cleanup(func() { http.DefaultTransport = previous })
	return transport
}

func assertBodyClosed(t *testing.T, transport *trackingTransport, msg string) {
	t.Helper()
	if transport.body == nil || !transport.body.closed {
		t.Fatal(msg)
	}
}

func setupAuthenticatedContext(t *testing.T, endpoint string) {
	t.Helper()
	origToken := TokenFlag
	t.Cleanup(func() { TokenFlag = origToken })
	TokenFlag = GetToken(t)
	origEndpoint := viper.GetString("contexts.local.endpoint")
	t.Cleanup(func() { viper.Set("contexts.local.endpoint", origEndpoint) })
	viper.Set("contexts.local.endpoint", endpoint)
	tokenData, err := os.ReadFile(GetToken(t))
	if err != nil {
		return
	}
	tempAuth := filepath.Join(t.TempDir(), "auth.json")
	if err := os.WriteFile(tempAuth, tokenData, 0o600); err != nil {
		return
	}
	viper.Set("tokens", []map[string]string{{"name": "default", "location": tempAuth}})
}

func TestIsLocalProvider(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"Local", true},
		{"local", true},
		{"None", true}, // legacy alias
		{"NONE", true},
		{"Meshery", false},
		{"", false},
	}
	for _, c := range cases {
		if got := IsLocalProvider(c.in); got != c.want {
			t.Errorf("IsLocalProvider(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestSafeClose(t *testing.T) {

	SetupMeshkitLoggerTesting(t, false)

	// testcases for SafeClose(co io.Closer)
	t.Run("SafeClose", func(t *testing.T) {
		// define a io.Closer for testing
		expectedErr := errors.New("close error")
		mc := &mockCloser{
			closeFunc: func() error {
				return expectedErr
			},
		}
		SafeClose(mc)

	})
}

func TestBackupConfigFile(t *testing.T) {
	cfgFile := "testdata/config.yaml"
	tmpFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpFile.Write(data); err != nil {
		t.Fatal(err)
	}
	_ = tmpFile.Close()

	// Act
	BackupConfigFile(tmpFile.Name())

	// Check backup in same dir as tmpFile
	dir := filepath.Dir(tmpFile.Name())
	backupFilePath := filepath.Join(dir, "config.bak.yaml")
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		t.Errorf("BackupConfigFile failed: backup file %s does not exist", backupFilePath)
	} else {
		// optional cleanup
		defer func() { _ = os.Remove(backupFilePath) }()
	}
}
func TestStringWithCharset(t *testing.T) {
	// checking the length, since this function returns random strings everytime
	strLength := 10

	gotString := StringWithCharset(strLength)

	if len(gotString) != strLength {
		t.Errorf("StringWithCharset got = %v want = %v", len(gotString), strLength)
	}
}

func TestPrereq(t *testing.T) {
	_, _, err := prereq()
	if err != nil {
		t.Errorf("prereq error = %v", err)
	}
}

func TestSetFileLocation(t *testing.T) {
	originalMesheryFolder := MesheryFolder
	originalDockerComposeFile := DockerComposeFile
	originalAuthConfigFile := AuthConfigFile
	originalDefaultConfigPath := DefaultConfigPath
	defer func() {
		MesheryFolder = originalMesheryFolder
		DockerComposeFile = originalDockerComposeFile
		AuthConfigFile = originalAuthConfigFile
		DefaultConfigPath = originalDefaultConfigPath
	}()
	err := SetFileLocation()
	if err != nil {
		t.Errorf("SetFileLocation error = %v", err)
	}
}

func TestNavigateToBrowser(t *testing.T) {
	// Save the original function.
	originalOpenURL := openURL

	// Restore it after this test finishes.
	t.Cleanup(func() {
		openURL = originalOpenURL
	})

	const testURL = "https://meshery.io"

	t.Run("success", func(t *testing.T) {
		var receivedURL string

		// Replace the real browser opener with a fake.
		openURL = func(url string) error {
			receivedURL = url
			return nil
		}

		err := NavigateToBrowser(testURL)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if receivedURL != testURL {
			t.Errorf("got URL %q, want %q",
				receivedURL, testURL)
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("browser unavailable")
		var receivedURL string

		// Simulate the browser opener failing.
		openURL = func(url string) error {
			receivedURL = url
			return expectedErr
		}

		err := NavigateToBrowser(testURL)

		if receivedURL != testURL {
			t.Errorf("got URL %q, want %q",
				receivedURL, testURL)
		}

		if !errors.Is(err, expectedErr) {
			t.Errorf("got error %v, want %v",
				err, expectedErr)
		}
	})
}

func TestContentTypeIsHTML(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		expectedOutput bool
	}{
		{
			name: "correct content-type",
			response: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"text/html"},
				},
			},
			expectedOutput: true,
		},
		{
			name: "empty content-type",
			response: &http.Response{
				Header: http.Header{
					"Content-Type": []string{},
				},
			},
			expectedOutput: false,
		},
		{
			name: "incorrect content-type",
			response: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"multipart/form-data"},
				},
			},
			expectedOutput: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContentTypeIsHTML(tt.response)
			if got != tt.expectedOutput {
				t.Errorf("ContentTypeIsHTML error = %v want = %v", got, tt.expectedOutput)
			}
		})
	}
}

func TestAskForConfirmation(t *testing.T) {
	tests := []struct {
		name     string
		question string
		input    string
		want     bool
		wantErr  bool
	}{
		{
			name:     "test with input 'yes'",
			question: "question?",
			input:    "yes\n",
			want:     true,
		},
		{
			name:     "test with input 'no'",
			question: "question?",
			input:    "no\n",
			want:     false,
		},
		{
			name:     "unrecognized input is asked again until a valid answer",
			question: "question?",
			input:    "maybe\ny\n",
			want:     true,
		},
		{
			name:     "stdin closed without an answer returns an error",
			question: "question?",
			input:    "",
			wantErr:  true,
		},
		{
			name:     "stdin closed after an unrecognized answer returns an error",
			question: "question?",
			input:    "maybe\n",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mocking stdio
			// see https://stackoverflow.com/a/64518829
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			_, err = w.WriteString(tt.input)
			if err != nil {
				t.Error(err)
			}
			_ = w.Close()
			stdin := os.Stdin
			defer func() { os.Stdin = stdin }()
			os.Stdin = r

			got, err := AskForConfirmation(tt.question)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("AskForConfirmation error = nil, want an error")
				}
				if got {
					t.Errorf("AskForConfirmation got = true on error, want false")
				}
				return
			}
			if err != nil {
				t.Fatalf("AskForConfirmation unexpected error = %v", err)
			}
			if got != tt.want {
				t.Errorf("AskForConfirmation got = %v want = %v", got, tt.want)
			}
		})
	}
}

func TestCreateConfigFile(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Error(err, "failed to get users home directory")
	}
	originalDefaultConfigPath := DefaultConfigPath
	defer func() { DefaultConfigPath = originalDefaultConfigPath }()
	DefaultConfigPath = filepath.Join(home, "config.yaml")
	err = CreateConfigFile()
	if err != nil {
		t.Errorf("CreateConfigFile error = %v", err)
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr error
	}{
		{
			name: "Correct URL",
			url:  "https://www.meshery.io",
		},
		{
			name:    "Unsupported scheme",
			url:     "mqtt://www.meshery.io",
			wantErr: ErrParsingUrl(fmt.Errorf("mqtt is not a supported protocol")),
		},
		{
			name:    "invalid URL",
			url:     "meshery.io",
			wantErr: ErrParsingUrl(fmt.Errorf("parse \"meshery.io\": invalid URI for request")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if err != nil {
				AssertMeshkitErrorsEqual(t, err, tt.wantErr)
			}
		})
	}
}

// func TestReadToken(t *testing.T) {
// }

func TestTruncateID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{name: "long id truncated to 8 chars", id: "1234567890", want: "12345678"},
		{name: "exact 8 char id unchanged", id: "12345678", want: "12345678"},
		{name: "short fallback unchanged", id: "N/A", want: "N/A"},
		{name: "empty string unchanged", id: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateID(tt.id); got != tt.want {
				t.Errorf("TruncateID(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	ts := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
	tests := []struct {
		name     string
		ts       *time.Time
		layout   string
		fallback string
		want     string
	}{
		{name: "nil timestamp returns fallback", ts: nil, layout: "01-02-2006", fallback: "N/A", want: "N/A"},
		{name: "non-nil date formatted", ts: &ts, layout: "01-02-2006", fallback: "N/A", want: "01-02-2025"},
		{name: "non-nil datetime formatted", ts: &ts, layout: "01-02-2006 15:04:05", fallback: "N/A", want: "01-02-2025 03:04:05"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTime(tt.ts, tt.layout, tt.fallback); got != tt.want {
				t.Errorf("FormatTime() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatUUID(t *testing.T) {
	id := core.Uuid{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef, 0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}
	tests := []struct {
		name     string
		id       *core.Uuid
		fallback string
		want     string
	}{
		{name: "nil uuid returns fallback", id: nil, fallback: "N/A", want: "N/A"},
		{name: "nil uuid returns null fallback", id: nil, fallback: "null", want: "null"},
		{name: "non-nil uuid formatted", id: &id, fallback: "N/A", want: "12345678-90ab-cdef-1234-567890abcdef"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatUUID(tt.id, tt.fallback); got != tt.want {
				t.Errorf("FormatUUID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStringContainedInSlice(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		slice []string
		want  int
	}{
		{
			name:  "test with present string",
			str:   "data2",
			slice: []string{"data1", "data2"},
			want:  1,
		},
		{
			name:  "test with absent string",
			str:   "data3",
			slice: []string{"data1", "data2"},
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringContainedInSlice(tt.str, tt.slice)
			if got != tt.want {
				t.Errorf("StringContainedInSlice got = %v want = %v", got, tt.want)
			}
		})
	}
}

func TestStringInSlice(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		slice []string
		want  bool
	}{
		{
			name:  "test with string present",
			str:   "data1",
			slice: []string{"data1", "data2"},
			want:  true,
		},
		{
			name:  "test with string absent",
			str:   "data3",
			slice: []string{"data1", "data2"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringInSlice(tt.str, tt.slice)
			if got != tt.want {
				t.Errorf("StringInSlice got = %v want = %v", got, tt.want)
			}
		})
	}
}

func TestAskForInput(t *testing.T) {
	input := "data1"

	// mocking stdio
	// see https://stackoverflow.com/a/64518829
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.WriteString(input + "\n")
	if err != nil {
		t.Error(err)
	}
	_ = w.Close()
	stdin := os.Stdin
	defer func() { os.Stdin = stdin }()
	os.Stdin = r

	got, err := AskForInput("Prompt", []string{"data1", "data2"})
	if err != nil {
		t.Fatalf("AskForInput unexpected error = %v", err)
	}
	if got != input {
		t.Errorf("AskForInput got = %v want = %v", got, input)
	}
}

func TestAskForInputInvalidThenValid(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = w.WriteString("nope\ndata2\n"); err != nil {
		t.Fatal(err)
	}
	_ = w.Close()
	stdin := os.Stdin
	defer func() { os.Stdin = stdin }()
	os.Stdin = r

	got, err := AskForInput("Prompt", []string{"data1", "data2"})
	if err != nil {
		t.Fatalf("AskForInput unexpected error = %v", err)
	}
	if got != "data2" {
		t.Errorf("AskForInput got = %v want = data2", got)
	}
}

func TestAskForInputStdinClosed(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_ = w.Close()
	stdin := os.Stdin
	defer func() { os.Stdin = stdin }()
	os.Stdin = r

	got, err := AskForInput("Prompt", []string{"data1", "data2"})
	if err == nil {
		t.Fatalf("AskForInput error = nil, want an error; got %q", got)
	}
	if got != "" {
		t.Errorf("AskForInput got = %q on error, want an empty string", got)
	}
}

// func TestPrintToTableInStringFormat(t *testing.T) {
// 	want := NewGoldenFile(t, "PrintToTableInStringFormat.golden", fixturesDir).Load()
// 	got := PrintToTableInStringFormat([]string{"firstheader", "secondheader"}, [][]string{{"data1", "data2"}, {"data3", "data4"}})
// 	if got != want {
// 		t.Errorf("PrintToTableInStringFormat got = %v want = %v", got, want)
// 	}
// }

func TestCreateDefaultSpinner(t *testing.T) {
	// only checking for Suffix and FinalMSG
	got := CreateDefaultSpinner("suffix", "message")
	want := struct {
		Suffix   string
		FinalMsg string
	}{
		Suffix:   " suffix", // The leading space is intentional
		FinalMsg: "message\n",
	}
	if want.Suffix != got.Suffix {
		t.Errorf("CreateDefaultSpinner got = %v want = %v", got.Suffix, want.Suffix)
	}
	if want.FinalMsg != got.FinalMSG {
		t.Errorf("CreateDefaultSpinner got = %v want = %v", got.FinalMSG, want.FinalMsg)
	}
}

func TestContainsStringPrefix(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		str   string
		want  bool
	}{
		{
			name:  "str is present in the slice",
			slice: []string{"data1", "data2"},
			str:   "data2",
			want:  true,
		},
		{
			name:  "str is not present in the slice",
			slice: []string{"data1", "data2"},
			str:   "data3",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsStringPrefix(tt.slice, tt.str)
			if got != tt.want {
				t.Errorf("ContainsStringPrefix got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidId(t *testing.T) {
	const (
		idA = "dc76a1b2-0000-4000-8000-000000000001"
		idB = "d4e5f6a7-0000-4000-8000-000000000002"
		idC = "1a2b3c4d-0000-4000-8000-000000000003"
	)
	manyIDs := make([]string, 12)
	for i := range manyIDs {
		manyIDs[i] = fmt.Sprintf("e0000000-0000-4000-8000-0000000000%02d", i)
	}
	tests := []struct {
		name         string
		ids          []string
		input        string
		wantID       string
		wantIsID     bool
		wantErr      error
		wantErrText  string
		wantRequests int
	}{
		{
			name:         "given an exact ID when ValidId then that ID is returned",
			ids:          []string{idA, idB, idC},
			input:        idB,
			wantID:       idB,
			wantIsID:     true,
			wantRequests: 1,
		},
		{
			name:         "given a prefix of one ID when ValidId then the full ID is returned",
			ids:          []string{idA, idB, idC},
			input:        "dc",
			wantID:       idA,
			wantIsID:     true,
			wantRequests: 1,
		},
		{
			name:         "given a prefix of several IDs when ValidId then an ambiguity error lists them",
			ids:          []string{idA, idB, idC},
			input:        "d",
			wantErr:      ErrAmbiguousID("d", []string{idA, idB}),
			wantRequests: 1,
		},
		{
			name:         "given a prefix of more than 10 IDs when ValidId then the error lists only the first 10",
			ids:          manyIDs,
			input:        "e",
			wantErr:      ErrAmbiguousID("e", manyIDs),
			wantErrText:  `"e" matches 12 IDs, including: ` + strings.Join(manyIDs[:10], ", "),
			wantRequests: 1,
		},
		{
			name:         "given input matching no ID when ValidId then the input is returned unchanged",
			ids:          []string{idA, idB, idC},
			input:        "my-design",
			wantID:       "my-design",
			wantIsID:     false,
			wantRequests: 1,
		},
		{
			name:         "given an exact ID that is also a prefix of another ID when ValidId then the exact ID wins",
			ids:          []string{idA + "0", idA},
			input:        idA,
			wantID:       idA,
			wantIsID:     true,
			wantRequests: 1,
		},
		{
			name:         "given an empty ID list when ValidId then the input is returned unchanged",
			ids:          []string{},
			input:        "d",
			wantID:       "d",
			wantIsID:     false,
			wantRequests: 1,
		},
		{
			name:         "given empty input when ValidId then it is rejected without a request",
			ids:          []string{idA, idB, idC},
			input:        "",
			wantErr:      ErrInvalidArgument(errors.New("name or ID cannot be empty")),
			wantRequests: 0,
		},
		{
			name:         "given whitespace input when ValidId then it is rejected without a request",
			ids:          []string{idA, idB, idC},
			input:        "   ",
			wantErr:      ErrInvalidArgument(errors.New("name or ID cannot be empty")),
			wantRequests: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				if got := r.URL.Query().Get("pagesize"); got != "10000" {
					t.Errorf("GetID sent pagesize=%q, want %q", got, "10000")
				}
				patterns := make([]string, 0, len(tt.ids))
				for _, id := range tt.ids {
					patterns = append(patterns, fmt.Sprintf(`{"id":%q}`, id))
				}
				_, _ = fmt.Fprintf(w, `{"patterns":[%s]}`, strings.Join(patterns, ","))
			}))
			defer server.Close()
			setupAuthenticatedContext(t, server.URL)

			gotID, gotIsID, err := ValidId(server.URL, tt.input, "pattern")

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				AssertMeshkitErrorsEqual(t, err, tt.wantErr)
				if tt.wantErrText != "" && !strings.HasPrefix(err.Error(), tt.wantErrText+" | ") {
					t.Errorf("ValidId error = %q, want it to start with %q", err.Error(), tt.wantErrText)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if gotID != tt.wantID || gotIsID != tt.wantIsID {
					t.Errorf("ValidId got = (%q, %v), want (%q, %v)", gotID, gotIsID, tt.wantID, tt.wantIsID)
				}
			}
			if requests != tt.wantRequests {
				t.Errorf("ValidId made %d requests, want %d", requests, tt.wantRequests)
			}
		})
	}
}

func TestSetOverrideValues(t *testing.T) {
	testChannel := "testChannel"

	tests := []struct {
		name                string
		ctx                 *config.Context
		mesheryImageVersion string
		customCallbackURL   string
		customProviderURL   string
		want                map[string]interface{}
	}{
		{
			name: "Context contains no components and no meshery image version",
			ctx: &config.Context{
				Components: nil,
				Channel:    testChannel,
			},
			mesheryImageVersion: "",
			want: map[string]interface{}{
				"meshery-istio": map[string]interface{}{
					"enabled": false,
				},
				"meshery-cilium": map[string]interface{}{
					"enabled": false,
				},
				"meshery-linkerd": map[string]interface{}{
					"enabled": false,
				},
				"meshery-consul": map[string]interface{}{
					"enabled": false,
				},
				"meshery-kuma": map[string]interface{}{
					"enabled": false,
				},
				"meshery-nsm": map[string]interface{}{
					"enabled": false,
				},
				"meshery-traefik-mesh": map[string]interface{}{
					"enabled": false,
				},
				"image": map[string]interface{}{
					"tag": testChannel + "-",
				},
			},
		},
		{
			name: "Context contains part of all available components and meshery image version",
			ctx: &config.Context{
				Components: []string{"meshery-istio", "meshery-nsm"},
				Channel:    testChannel,
			},
			mesheryImageVersion: "testImageVersion",
			want: map[string]interface{}{
				"meshery-istio": map[string]interface{}{
					"enabled": true,
				},
				"meshery-cilium": map[string]interface{}{
					"enabled": false,
				},
				"meshery-linkerd": map[string]interface{}{
					"enabled": false,
				},
				"meshery-consul": map[string]interface{}{
					"enabled": false,
				},
				"meshery-kuma": map[string]interface{}{
					"enabled": false,
				},
				"meshery-nsm": map[string]interface{}{
					"enabled": true,
				},
				"meshery-traefik-mesh": map[string]interface{}{
					"enabled": false,
				},
				"image": map[string]interface{}{
					"tag": testChannel + "-testImageVersion",
				},
			},
		},
		{
			name: "Context contains all available components and meshery image version",
			ctx: &config.Context{
				Components: []string{"meshery-istio", "meshery-linkerd", "meshery-consul", "meshery-kuma",
					"meshery-nsm", "meshery-traefik-mesh", "meshery-cilium"},
				Channel: testChannel,
			},
			mesheryImageVersion: "testImageVersion",
			want: map[string]interface{}{
				"meshery-istio": map[string]interface{}{
					"enabled": true,
				},
				"meshery-cilium": map[string]interface{}{
					"enabled": true,
				},
				"meshery-linkerd": map[string]interface{}{
					"enabled": true,
				},
				"meshery-consul": map[string]interface{}{
					"enabled": true,
				},
				"meshery-kuma": map[string]interface{}{
					"enabled": true,
				},
				"meshery-nsm": map[string]interface{}{
					"enabled": true,
				},
				"meshery-traefik-mesh": map[string]interface{}{
					"enabled": true,
				},
				"image": map[string]interface{}{
					"tag": testChannel + "-testImageVersion",
				},
			},
		},
		{
			name: "Context contains custom env variables and disabled operator and provider, custom callback url and provider url is specified",
			ctx: &config.Context{
				Components: nil,
				Channel:    testChannel,
				Operator:   "disabled",
				Provider:   "Custom",
				EnvVars: map[string]any{
					"debug":                            "0",
					"playground":                       true,
					"meshsync_default_deployment_mode": "embedded",
					"custom_var_int":                   1,
					"custom_var_int_string":            "1",
					"custom_var_bool":                  true,
					"custom_var_bool_string":           "true",
				},
			},
			customCallbackURL:   "http://localhost:9081/api/system/oauth/callback",
			customProviderURL:   "http://localhost:9876",
			mesheryImageVersion: "",
			want: map[string]interface{}{
				"meshery-istio": map[string]interface{}{
					"enabled": false,
				},
				"meshery-cilium": map[string]interface{}{
					"enabled": false,
				},
				"meshery-linkerd": map[string]interface{}{
					"enabled": false,
				},
				"meshery-consul": map[string]interface{}{
					"enabled": false,
				},
				"meshery-kuma": map[string]interface{}{
					"enabled": false,
				},
				"meshery-nsm": map[string]interface{}{
					"enabled": false,
				},
				"meshery-traefik-mesh": map[string]interface{}{
					"enabled": false,
				},
				"image": map[string]interface{}{
					"tag": testChannel + "-",
				},
				// non string values (numeric and bool) must go as string in quotes
				// otherwise helm will return unmarshaling error
				"env": map[string]any{
					"DISABLE_OPERATOR":                 "\"true\"",
					constants.ProviderENV:              "Custom",
					constants.CallbackURLENV:           "http://localhost:9081/api/system/oauth/callback",
					constants.ProviderURLsENV:          "http://localhost:9876",
					"DEBUG":                            "\"0\"",
					"PLAYGROUND":                       "\"true\"",
					"MESHSYNC_DEFAULT_DEPLOYMENT_MODE": "embedded",
					"CUSTOM_VAR_INT":                   "\"1\"",
					"CUSTOM_VAR_INT_STRING":            "\"1\"",
					"CUSTOM_VAR_BOOL":                  "\"true\"",
					"CUSTOM_VAR_BOOL_STRING":           "\"true\"",
				},
			},
		},
	}

	for _, tt := range tests {
		got := SetOverrideValues(tt.ctx, tt.mesheryImageVersion, tt.customCallbackURL, tt.customProviderURL)
		eq := reflect.DeepEqual(got, tt.want)
		if !eq {
			t.Errorf("SetOverrideValues %s got = %v want = %v", tt.name, got, tt.want)
		}
	}
}

// MockClient implements the client interface for testing
type MockClient struct {
	GetKubeConfigFunc func() (*api.Config, error)
}

// GetKubeConfig implements the method from meshkitkube.Client interface
func (m *MockClient) GetKubeConfig() (*api.Config, error) {
	if m.GetKubeConfigFunc != nil {
		return m.GetKubeConfigFunc()
	}
	return nil, fmt.Errorf("GetKubeConfig not implemented")
}

func getCurrentK8sContextForTest(client *MockClient) (string, error) {
	if client == nil {
		return "", fmt.Errorf("kubernetes client is nil")
	}
	config, err := client.GetKubeConfig()
	if err != nil {
		return "", err
	}
	return config.CurrentContext, nil
}

func TestGetCurrentK8sContext(t *testing.T) {
	tests := []struct {
		name          string
		client        *MockClient
		wantContext   string
		wantErrString string
	}{
		{
			name:          "nil client",
			client:        nil,
			wantContext:   "",
			wantErrString: "kubernetes client is nil",
		},
		{
			name: "client returns error",
			client: &MockClient{
				GetKubeConfigFunc: func() (*api.Config, error) {
					return nil, fmt.Errorf("test error")
				},
			},
			wantContext:   "",
			wantErrString: "test error",
		},
		{
			name: "successful context retrieval",
			client: &MockClient{
				GetKubeConfigFunc: func() (*api.Config, error) {
					return &api.Config{
						CurrentContext: "minikube",
					}, nil
				},
			},
			wantContext:   "minikube",
			wantErrString: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, err := getCurrentK8sContextForTest(tt.client)

			// Check context
			if context != tt.wantContext {
				t.Errorf("getCurrentK8sContextForTest() context = %v, want %v", context, tt.wantContext)
			}

			// Check error
			if tt.wantErrString == "" {
				if err != nil {
					t.Errorf("getCurrentK8sContextForTest() error = %v, want nil", err)
				}
			} else {
				if err == nil {
					t.Errorf("getCurrentK8sContextForTest() error = nil, want error containing %v", tt.wantErrString)
				} else if !strings.Contains(err.Error(), tt.wantErrString) {
					t.Errorf("getCurrentK8sContextForTest() error = %v, want error containing %v", err.Error(), tt.wantErrString)
				}
			}
		})
	}
}

// TestSubcommandNames pins the contract that makes an "invalid subcommand"
// message self-maintaining: the advertised list is derived from the same slice
// the command validates against, so adding a subcommand cannot leave the error
// message advertising a stale set.
func TestSubcommandNames(t *testing.T) {
	tests := []struct {
		name      string
		available []*cobra.Command
		expected  string
	}{
		{
			name:      "no subcommands",
			available: nil,
			expected:  "",
		},
		{
			name:      "single subcommand",
			available: []*cobra.Command{{Use: "view"}},
			expected:  "view",
		},
		{
			name: "names are sorted, not declaration-ordered",
			available: []*cobra.Command{
				{Use: "view"}, {Use: "generate"}, {Use: "list"}, {Use: "search"},
			},
			expected: "generate, list, search, view",
		},
		{
			name:      "usage line is reduced to the command name",
			available: []*cobra.Command{{Use: "view [model-name]"}, {Use: "list [flags]"}},
			expected:  "list, view",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubcommandNames(tt.available); got != tt.expected {
				t.Errorf("SubcommandNames() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestTrimFileExtension checks that only the final extension is removed.
// Regression test: strings.Trim(name, ext) strips a character set from both
// ends, turning "my-app.yaml" into "-app". Use strings.TrimSuffix instead.
func TestTrimFileExtension(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"leading chars match extension", "my-app.yaml", "my-app"},
		{"mostly extension letters", "mail.yaml", "mail"},
		{"json extension", "sample.json", "sample"},
		{"multiple dots", "app.v2.yaml", "app.v2"},
		{"no extension", "Test Design", "Test Design"},
		{"empty name", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimFileExtension(tt.in); got != tt.want {
				t.Errorf("TrimFileExtension(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestUploadFileWithParams(t *testing.T) {
	origToken := TokenFlag
	defer func() { TokenFlag = origToken }()
	TokenFlag = GetToken(t)

	tmpDir := t.TempDir()
	sampleFile := filepath.Join(tmpDir, "sample.txt")
	err := os.WriteFile(sampleFile, []byte("sample content for upload"), 0o600)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	params := map[string]string{
		"token": "valid-token",
		"org":   "meshery",
	}

	req, err := UploadFileWithParams("http://localhost:9081/api/upload", params, "file", sampleFile)
	if err != nil {
		t.Fatalf("UploadFileWithParams returned unexpected error: %v", err)
	}
	if req == nil {
		t.Fatal("UploadFileWithParams returned nil request")
	}
	if req.Method != http.MethodPost {
		t.Errorf("expected POST method, got %s", req.Method)
	}
	if !strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data") {
		t.Errorf("expected multipart/form-data content type, got %s", req.Header.Get("Content-Type"))
	}

	nonExistent := filepath.Join(tmpDir, "does-not-exist.txt")
	_, err = UploadFileWithParams("http://localhost:9081/api/upload", params, "file", nonExistent)
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestDeleteConfiguration_BodyClosed(t *testing.T) {
	_ = SetupMeshkitLoggerTesting(t, false)

	t.Run("closes body on success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{"status":"deleted"}`)
		}))
		defer server.Close()
		setupAuthenticatedContext(t, server.URL)
		transport := installTrackingTransport(t)
		if err := DeleteConfiguration(server.URL, "test-id", "pattern"); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		assertBodyClosed(t, transport, "expected response body to be closed on delete configuration")
	})
}
