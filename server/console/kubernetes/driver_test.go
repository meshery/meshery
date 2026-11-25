package kubernetes

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/meshery/meshery/server/console"
	"github.com/meshery/meshkit/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
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

// labelledPod is a pod stamped with the labels an owner workload selects on.
func labelledPod(name, namespace string, phase corev1.PodPhase, labels map[string]string, containers ...string) *corev1.Pod {
	p := pod(name, namespace, phase, containers...)
	p.Labels = labels
	return p
}

func deployment(name, namespace string, selector map[string]string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec:       appsv1.DeploymentSpec{Selector: &metav1.LabelSelector{MatchLabels: selector}},
	}
}

// A console named after an owner workload must resolve to one of the pods that
// workload runs. Before this, a workload's name was passed straight to
// `pods/<name>`, so every deployment row failed with "pods not found".
func TestCapabilitiesResolvesOwnerWorkloads(t *testing.T) {
	selector := map[string]string{"app": "meshsync"}
	running := labelledPod("meshery-meshsync-7d9f", "meshery", corev1.PodRunning, selector, "meshsync")

	t.Run("deployment", func(t *testing.T) {
		d := &driver{clientset: fake.NewSimpleClientset(
			deployment("meshery-meshsync", "meshery", selector),
			running,
			// A pod from an unrelated workload in the same namespace must not be
			// picked up by the selector list.
			labelledPod("unrelated-abc", "meshery", corev1.PodRunning, map[string]string{"app": "other"}, "other"),
		)}

		caps, err := d.Capabilities(context.Background(), console.TargetRef{
			Resource: "Deployment", Namespace: "meshery", Name: "meshery-meshsync",
		})
		if err != nil {
			t.Fatalf("Capabilities for a deployment: %v", err)
		}
		if !caps.Terminal || !caps.Logs {
			t.Errorf("caps = %+v, want a running pod's console kinds", caps)
		}
		if len(caps.Containers) != 1 || caps.Containers[0] != "meshsync" {
			t.Errorf("Containers = %v, want the resolved pod's containers", caps.Containers)
		}
	})

	t.Run("statefulset", func(t *testing.T) {
		d := &driver{clientset: fake.NewSimpleClientset(
			&appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{Name: "db", Namespace: "default"},
				Spec:       appsv1.StatefulSetSpec{Selector: &metav1.LabelSelector{MatchLabels: selector}},
			},
			labelledPod("db-0", "default", corev1.PodRunning, selector, "postgres"),
		)}

		caps, err := d.Capabilities(context.Background(), console.TargetRef{
			Resource: "statefulset", Namespace: "default", Name: "db",
		})
		if err != nil {
			t.Fatalf("Capabilities for a statefulset: %v", err)
		}
		if caps.DefaultContainer != "postgres" {
			t.Errorf("DefaultContainer = %q, want %q", caps.DefaultContainer, "postgres")
		}
	})

	t.Run("job", func(t *testing.T) {
		d := &driver{clientset: fake.NewSimpleClientset(
			&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{Name: "migration", Namespace: "default"},
				Spec:       batchv1.JobSpec{Selector: &metav1.LabelSelector{MatchLabels: selector}},
			},
			labelledPod("migration-x1", "default", corev1.PodFailed, selector, "migrate"),
		)}

		caps, err := d.Capabilities(context.Background(), console.TargetRef{
			Resource: "Job", Namespace: "default", Name: "migration",
		})
		if err != nil {
			t.Fatalf("Capabilities for a job: %v", err)
		}
		// A finished job has no process to exec into, but its logs are exactly
		// why a user opens the console.
		if caps.Terminal {
			t.Error("Terminal = true for a failed pod, want false")
		}
		if !caps.Logs {
			t.Error("Logs = false for a failed pod, want true")
		}
	})
}

// A workload that exists but runs no pods is an unsupported console with a
// reason, not a 404: the target the user clicked does exist.
func TestCapabilitiesWorkloadWithNoPods(t *testing.T) {
	d := &driver{clientset: fake.NewSimpleClientset(
		deployment("meshery-meshsync", "meshery", map[string]string{"app": "meshsync"}),
	)}

	caps, err := d.Capabilities(context.Background(), console.TargetRef{
		Resource: "deployment", Namespace: "meshery", Name: "meshery-meshsync",
	})
	if err != nil {
		t.Fatalf("Capabilities for a deployment with no pods: %v", err)
	}
	if caps.Terminal || caps.Logs {
		t.Errorf("caps = %+v, want no supported console kinds", caps)
	}
	if caps.Reason == "" {
		t.Error("Reason is empty; the UI needs something to show")
	}
}

// A deployment that was deleted since the UI listed it is a missing target, not
// an unsupported one.
func TestCapabilitiesMissingWorkload(t *testing.T) {
	d := &driver{clientset: fake.NewSimpleClientset()}

	_, err := d.Capabilities(context.Background(), console.TargetRef{
		Resource: "deployment", Namespace: "meshery", Name: "ghost",
	})
	if err == nil {
		t.Fatal("Capabilities for a missing deployment succeeded")
	}
	if code := errors.GetCode(err); code != console.ErrTargetNotFoundCode {
		t.Errorf("error code = %q, want %q", code, console.ErrTargetNotFoundCode)
	}
}

// A workload with no selector names no pods. Listing every pod in the namespace
// would attach the console to somebody else's workload, so the console is
// refused instead.
func TestWorkloadWithoutSelectorNamesNoPods(t *testing.T) {
	d := &driver{clientset: fake.NewSimpleClientset(
		deployment("bare", "default", nil),
		labelledPod("someone-else", "default", corev1.PodRunning, map[string]string{"app": "other"}, "app"),
	)}

	caps, err := d.Capabilities(context.Background(), console.TargetRef{
		Resource: "deployment", Namespace: "default", Name: "bare",
	})
	if err != nil {
		t.Fatalf("Capabilities: %v", err)
	}
	if caps.Logs || caps.Terminal {
		t.Errorf("caps = %+v, want a deployment with no pods to admit nothing", caps)
	}
	if len(caps.Containers) != 0 {
		t.Errorf("Containers = %v, want none", caps.Containers)
	}
}

func TestSelectPod(t *testing.T) {
	labels := map[string]string{"app": "app"}
	at := func(offset time.Duration) metav1.Time {
		return metav1.NewTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Add(offset))
	}

	pods := []corev1.Pod{
		*labelledPod("old-running", "default", corev1.PodRunning, labels, "app"),
		*labelledPod("new-pending", "default", corev1.PodPending, labels, "app"),
	}
	pods[0].CreationTimestamp = at(0)
	pods[1].CreationTimestamp = at(time.Hour)

	if got := selectPod(pods); got.Name != "old-running" {
		t.Errorf("selectPod = %q, want the only running pod", got.Name)
	}

	pods[1].Status.Phase = corev1.PodRunning
	if got := selectPod(pods); got.Name != "new-pending" {
		t.Errorf("selectPod = %q, want the most recently created running pod", got.Name)
	}

	// Equal creation times must resolve reproducibly, not by the order the API
	// server happened to return the pods in.
	pods[1].CreationTimestamp = at(0)
	pods[1].Name = "aaa-running"
	if got := selectPod(pods); got.Name != "aaa-running" {
		t.Errorf("selectPod = %q, want the name tie-break to decide", got.Name)
	}
}

func TestPodForResolvesWorkloadToPod(t *testing.T) {
	selector := map[string]string{"app": "meshsync"}
	d := &driver{clientset: fake.NewSimpleClientset(
		deployment("meshery-meshsync", "meshery", selector),
		labelledPod("meshery-meshsync-7d9f", "meshery", corev1.PodRunning, selector, "meshsync"),
	)}

	pod, err := d.podFor(context.Background(), "deployment", console.TargetRef{
		Resource: "Deployment", Namespace: "meshery", Name: "meshery-meshsync",
	})
	if err != nil {
		t.Fatalf("podFor: %v", err)
	}
	if pod == nil || pod.Name != "meshery-meshsync-7d9f" {
		t.Errorf("podFor = %v, want the pod the deployment selects", pod)
	}
}

// A workload with no pods must be refused rather than silently tailing the logs
// of an unrelated pod in the same namespace.
func TestStreamLogsRefusesWorkloadWithNoPods(t *testing.T) {
	d := &driver{clientset: fake.NewSimpleClientset(
		deployment("meshery-meshsync", "meshery", map[string]string{"app": "meshsync"}),
		labelledPod("someone-else", "meshery", corev1.PodRunning, map[string]string{"app": "other"}, "other"),
	)}

	_, err := d.StreamLogs(context.Background(), console.TargetRef{
		Resource: "deployment", Namespace: "meshery", Name: "meshery-meshsync",
	}, console.LogOptions{})
	if err == nil {
		t.Fatal("StreamLogs opened a console against a deployment with no pods")
	}
	if code := errors.GetCode(err); code != console.ErrUnsupportedConsoleCode {
		t.Errorf("error code = %q, want %q", code, console.ErrUnsupportedConsoleCode)
	}
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
