package kubernetes

import (
	"context"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/console"
	"github.com/meshery/meshkit/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func pod(name, namespace string, phase corev1.PodPhase, containers ...string) *corev1.Pod {
	spec := corev1.PodSpec{}
	for _, c := range containers {
		spec.Containers = append(spec.Containers, corev1.Container{Name: c})
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec:       spec,
		Status:     corev1.PodStatus{Phase: phase},
	}
}

func TestCapabilitiesRunningPod(t *testing.T) {
	d := &driver{clientset: fake.NewSimpleClientset(pod("web", "default", corev1.PodRunning, "app", "sidecar"))}

	caps, err := d.Capabilities(context.Background(), console.TargetRef{
		Resource: "pod", Namespace: "default", Name: "web",
	})
	if err != nil {
		t.Fatalf("Capabilities: %v", err)
	}
	if !caps.Terminal {
		t.Error("Terminal = false for a running pod, want true")
	}
	if !caps.Logs {
		t.Error("Logs = false, want true")
	}
	if caps.DefaultContainer != "app" {
		t.Errorf("DefaultContainer = %q, want the first container %q", caps.DefaultContainer, "app")
	}
	if len(caps.Containers) != 2 {
		t.Errorf("Containers = %v, want both containers", caps.Containers)
	}
}

// TestCapabilitiesNonRunningPod: logs remain available (a crashed pod's logs are
// exactly what a user comes for) but there is no process to exec into.
func TestCapabilitiesNonRunningPod(t *testing.T) {
	for _, phase := range []corev1.PodPhase{corev1.PodPending, corev1.PodFailed, corev1.PodSucceeded} {
		t.Run(string(phase), func(t *testing.T) {
			d := &driver{clientset: fake.NewSimpleClientset(pod("web", "default", phase, "app"))}

			caps, err := d.Capabilities(context.Background(), console.TargetRef{
				Resource: "pod", Namespace: "default", Name: "web",
			})
			if err != nil {
				t.Fatalf("Capabilities: %v", err)
			}
			if caps.Terminal {
				t.Errorf("Terminal = true for a %s pod, want false", phase)
			}
			if !caps.Logs {
				t.Errorf("Logs = false for a %s pod, want true", phase)
			}
			if caps.Reason == "" {
				t.Error("Reason is empty; an unsupported kind must explain itself")
			}
		})
	}
}

func TestCapabilitiesIncludesInitAndEphemeralContainers(t *testing.T) {
	p := pod("web", "default", corev1.PodRunning, "app")
	p.Spec.InitContainers = []corev1.Container{{Name: "init-db"}}
	p.Spec.EphemeralContainers = []corev1.EphemeralContainer{{
		EphemeralContainerCommon: corev1.EphemeralContainerCommon{Name: "debugger"},
	}}
	d := &driver{clientset: fake.NewSimpleClientset(p)}

	caps, err := d.Capabilities(context.Background(), console.TargetRef{
		Resource: "pod", Namespace: "default", Name: "web",
	})
	if err != nil {
		t.Fatalf("Capabilities: %v", err)
	}

	want := map[string]bool{"app": true, "init-db": true, "debugger": true}
	if len(caps.Containers) != len(want) {
		t.Fatalf("Containers = %v, want %v", caps.Containers, want)
	}
	for _, name := range caps.Containers {
		if !want[name] {
			t.Errorf("unexpected container %q", name)
		}
	}
}

func TestCapabilitiesHonoursDefaultContainerAnnotation(t *testing.T) {
	p := pod("web", "default", corev1.PodRunning, "app", "sidecar")
	p.Annotations = map[string]string{defaultContainerAnnotation: "sidecar"}
	d := &driver{clientset: fake.NewSimpleClientset(p)}

	caps, err := d.Capabilities(context.Background(), console.TargetRef{
		Resource: "pod", Namespace: "default", Name: "web",
	})
	if err != nil {
		t.Fatalf("Capabilities: %v", err)
	}
	if caps.DefaultContainer != "sidecar" {
		t.Errorf("DefaultContainer = %q, want the annotated %q", caps.DefaultContainer, "sidecar")
	}
}

// A stale annotation naming a container that no longer exists must not select
// a container the pod does not have.
func TestCapabilitiesIgnoresStaleDefaultContainerAnnotation(t *testing.T) {
	p := pod("web", "default", corev1.PodRunning, "app")
	p.Annotations = map[string]string{defaultContainerAnnotation: "removed"}
	d := &driver{clientset: fake.NewSimpleClientset(p)}

	caps, err := d.Capabilities(context.Background(), console.TargetRef{
		Resource: "pod", Namespace: "default", Name: "web",
	})
	if err != nil {
		t.Fatalf("Capabilities: %v", err)
	}
	if caps.DefaultContainer != "app" {
		t.Errorf("DefaultContainer = %q, want a fallback to %q", caps.DefaultContainer, "app")
	}
}

func TestCapabilitiesNonPodResource(t *testing.T) {
	d := &driver{clientset: fake.NewSimpleClientset()}

	caps, err := d.Capabilities(context.Background(), console.TargetRef{
		Resource: "configmap", Namespace: "default", Name: "settings",
	})
	if err != nil {
		t.Fatalf("Capabilities on a non-pod resource returned an error: %v", err)
	}
	if caps.Terminal || caps.Logs {
		t.Errorf("caps = %+v, want no supported console kinds", caps)
	}
	if caps.Reason == "" {
		t.Error("Reason is empty; the UI needs something to show")
	}
}

func TestCapabilitiesMissingPod(t *testing.T) {
	d := &driver{clientset: fake.NewSimpleClientset()}

	_, err := d.Capabilities(context.Background(), console.TargetRef{
		Resource: "pod", Namespace: "default", Name: "ghost",
	})
	if err == nil {
		t.Fatal("Capabilities for a missing pod succeeded")
	}
	if code := errors.GetCode(err); code != console.ErrTargetNotFoundCode {
		t.Errorf("error code = %q, want %q", code, console.ErrTargetNotFoundCode)
	}
}

func TestCapabilitiesRequiresNamespaceAndName(t *testing.T) {
	d := &driver{clientset: fake.NewSimpleClientset()}

	for _, target := range []console.TargetRef{
		{Resource: "pod", Name: "web"},
		{Resource: "pod", Namespace: "default"},
	} {
		_, err := d.Capabilities(context.Background(), target)
		if err == nil {
			t.Fatalf("Capabilities(%+v) succeeded, want an invalid-target error", target)
		}
		if code := errors.GetCode(err); code != console.ErrInvalidTargetCode {
			t.Errorf("error code = %q, want %q", code, console.ErrInvalidTargetCode)
		}
	}
}

func TestResolveContainer(t *testing.T) {
	caps := console.Capabilities{Containers: []string{"app", "sidecar"}, DefaultContainer: "app"}

	t.Run("defaults when unspecified", func(t *testing.T) {
		got, err := resolveContainer(console.TargetRef{}, caps)
		if err != nil || got != "app" {
			t.Errorf("resolveContainer = (%q, %v), want (\"app\", nil)", got, err)
		}
	})

	t.Run("honours an explicit container", func(t *testing.T) {
		got, err := resolveContainer(console.TargetRef{Container: "sidecar"}, caps)
		if err != nil || got != "sidecar" {
			t.Errorf("resolveContainer = (%q, %v), want (\"sidecar\", nil)", got, err)
		}
	})

	t.Run("rejects an unknown container", func(t *testing.T) {
		_, err := resolveContainer(console.TargetRef{Container: "nope"}, caps)
		if err == nil {
			t.Fatal("resolveContainer accepted a container the pod does not have")
		}
		if code := errors.GetCode(err); code != console.ErrInvalidTargetCode {
			t.Errorf("error code = %q, want %q", code, console.ErrInvalidTargetCode)
		}
	})

	t.Run("rejects a pod with no containers", func(t *testing.T) {
		_, err := resolveContainer(console.TargetRef{}, console.Capabilities{})
		if err == nil {
			t.Fatal("resolveContainer succeeded against a pod with no containers")
		}
	})
}

// The driver registers itself on import, which is what lets the HTTP layer stay
// free of Kubernetes.
func TestDriverRegistersItself(t *testing.T) {
	for _, kind := range console.Default.Kinds() {
		if kind == ConnectionKind {
			return
		}
	}
	t.Errorf("Default registry kinds = %v, want it to include %q", console.Default.Kinds(), ConnectionKind)
}

// The default shell must actually fall back: `exec bash || exec sh` does not,
// because a failed exec terminates the shell.
func TestDefaultShellCommandFallsBack(t *testing.T) {
	if len(defaultShellCommand) != 3 {
		t.Fatalf("defaultShellCommand = %v, want a three-part sh -c invocation", defaultShellCommand)
	}
	script := defaultShellCommand[2]
	if !strings.Contains(script, "[ -x /bin/bash ]") {
		t.Errorf("script %q does not test for bash before exec'ing it", script)
	}
	if !strings.Contains(script, "exec /bin/sh") {
		t.Errorf("script %q has no /bin/sh fallback", script)
	}
}
