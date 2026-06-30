package config

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// var update = flag.Bool("update", false, "update golden files")

var tests = []string{"ab2q3^$serhj", "klg434%$^", "m234@#$n", "op$#%G", "q$@#$r", "st4@#$23", "uv$@#$", "$@#TGGB#wx", "y@#$FG@$z", "1234"}

func TestGetCommitSHA(t *testing.T) {
	for _, test := range tests {
		version := Version{"", test, ""}
		got := version.GetCommitSHA()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestGetBuild(t *testing.T) {
	for _, test := range tests {
		version := Version{test, "", ""}
		got := version.GetBuild()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetLocation(t *testing.T) {
	token := Token{}
	for _, test := range tests {
		token.SetLocation(test)
		got := token.GetLocation()
		want, err := os.UserHomeDir()
		MesheryFolder := ".meshery"
		path := filepath.Join(want, MesheryFolder, test)
		want = path
		if err != nil {
			t.Errorf("Fail")
		}
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetVersion(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetVersion(test)
		got := meshConfigCtx.GetVersion()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetName(t *testing.T) {
	token := Token{}
	for _, test := range tests {
		token.SetName(test)
		got := token.Name
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetChannel(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetChannel(test)
		got := meshConfigCtx.GetChannel()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetComponents(t *testing.T) {
	dummy := []string{"abc", "def", "ghi", "jkl", "mno", "pqr"}
	meshConfigCtx := Context{"", "", "", dummy, "", "", "", "", nil}
	got := meshConfigCtx.GetComponents()
	want := dummy
	for i, j := range got {
		if j != want[i] {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetPlatform(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetPlatform(test)
		got := meshConfigCtx.GetPlatform()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetToken(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetToken(test)
		got := meshConfigCtx.GetToken()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetEndpoint(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetEndpoint(test)
		got := meshConfigCtx.GetEndpoint()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestGetCurrentContextName(t *testing.T) {
	for _, test := range tests {
		mesheryctlconfig := MesheryCtlConfig{nil, test, nil}
		got := mesheryctlconfig.GetCurrentContextName()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestSetContext(t *testing.T) {
	for _, test := range tests {
		mesheryctlconfig := MesheryCtlConfig{nil, test, nil}
		err := UpdateContextInConfig(nil, test)
		if err != nil {
			fmt.Print("Fail") //Internal:need to be fixed
		}
		got := mesheryctlconfig.GetCurrentContextName()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestSetEndpoint(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetEndpoint(test)
		got := meshConfigCtx.GetEndpoint()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetToken(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetToken(test)
		got := meshConfigCtx.GetToken()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetPlatform(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetPlatform(test)
		got := meshConfigCtx.GetPlatform()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetChannel(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetChannel(test)
		got := meshConfigCtx.GetChannel()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestSetVersion(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetVersion(test)
		got := meshConfigCtx.GetVersion()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetComponents(t *testing.T) {
	dummy := []string{"abc", "def", "ghi", "jkl", "mno", "pqr"}
	meshConfigCtx := Context{"", "", "", dummy, "", "", "", "", nil}
	got := meshConfigCtx.GetComponents()
	want := dummy
	for i, j := range got {
		if j != want[i] {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetProvider(t *testing.T) {
	for _, test := range tests {
		meshConfigCtx := Context{"", "", "", nil, "", "", test, "", nil}
		got := meshConfigCtx.GetProvider()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetProvider(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetProvider(test)
		got := meshConfigCtx.GetProvider()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetOperatorStatus(t *testing.T) {
	for _, test := range tests {
		meshConfigCtx := Context{"", "", "", nil, "", "", "", test, nil}
		got := meshConfigCtx.GetOperatorStatus()
		want := test

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}
}

func TestSetOperatorStatus(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetOperatorStatus(test)
		got := meshConfigCtx.GetOperatorStatus()
		want := test

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}
}

func TestValidateVersion(t *testing.T) {
	t.Run("empty version defaults to latest", func(t *testing.T) {
		ctx := &Context{Version: ""}
		if err := ctx.ValidateVersion(); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
		if ctx.Version != "latest" {
			t.Fatalf("expected 'latest', got %q", ctx.Version)
		}
	})

	t.Run("latest version returns nil", func(t *testing.T) {
		ctx := &Context{Version: "latest"}
		if err := ctx.ValidateVersion(); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("timeout when server is unreachable", func(t *testing.T) {
		slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(15 * time.Second)
		}))
		defer slow.Close()

		orig := validateVersionURL
		validateVersionURL = func(_, _, _ string) string { return slow.URL }
		defer func() { validateVersionURL = orig }()

		ctx := &Context{Version: "v0.7.0-test"}
		err := ctx.ValidateVersion()
		if err == nil {
			t.Fatal("expected timeout error")
		}
		if !strings.Contains(err.Error(), "timeout") {
			t.Fatalf("expected timeout error, got: %v", err)
		}
	})
}
