// Package kubernetes implements the consoles driver for connections of kind
// "kubernetes": interactive terminals backed by the pods/exec subresource, and
// log tails backed by pods/log.
//
// A console always attaches to a pod. A target may name an owner workload
// (Deployment, StatefulSet, ...) instead; it is resolved to one of the pods that
// workload runs before any console is opened.
//
// Importing this package registers the driver with the consoles registry. The
// HTTP layer never imports it directly; it is pulled in for its side effect
// wherever consoles are served.
package kubernetes

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strings"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/console"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/httpstream"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	k8sexec "k8s.io/client-go/util/exec"
)

const (
	// ConnectionKind is the connections.Connection.Kind this driver serves.
	ConnectionKind = "kubernetes"

	// ResourcePod is the resource a console attaches to directly. Terminals and
	// logs are properties of a running container, and a pod is the smallest thing
	// that has one.
	ResourcePod = "pod"

	// The owner workloads a console target may name. Each is resolved to one of
	// the pods it runs, so a console opened on a Deployment lands on a pod rather
	// than asking the API server for a pod called after the deployment.
	//
	// CronJob is deliberately absent: its pods belong to the Jobs it creates, and
	// naming one of those Jobs (or its pod) is the honest way to reach them.
	ResourceDeployment            = "deployment"
	ResourceStatefulSet           = "statefulset"
	ResourceDaemonSet             = "daemonset"
	ResourceJob                   = "job"
	ResourceReplicaSet            = "replicaset"
	ResourceReplicationController = "replicationcontroller"

	// defaultContainerAnnotation is the well-known annotation kubectl honours
	// to pick a container when the user does not name one.
	defaultContainerAnnotation = "kubectl.kubernetes.io/default-container"
)

// defaultShellCommand is used when the client does not specify a command.
//
// `exec` replaces the shell, and a failed `exec` terminates it, so the usual
// `exec bash || exec sh` does not actually fall back. Testing for bash first
// and only then exec'ing it does.
var defaultShellCommand = []string{
	"/bin/sh",
	"-c",
	"export TERM=xterm-256color; [ -x /bin/bash ] && exec /bin/bash || exec /bin/sh",
}

func init() {
	console.Register(factory{})
}

// factory builds a driver per Kubernetes connection.
type factory struct{}

func (factory) Kind() string { return ConnectionKind }

// NewDriver resolves the connection's credential into a Kubernetes client.
//
// The credential is fetched with the requesting user's token, so the console
// runs as that user against the cluster and Kubernetes RBAC governs what it may
// do. Meshery adds no authorization of its own beyond deciding who may reach
// the connection: a user who cannot `create pods/exec` gets a 403 from the API
// server, exactly as they would with kubectl.
func (factory) NewDriver(_ context.Context, cc console.ConnectionContext) (console.Driver, error) {
	k8sContext, err := models.K8sContextFromConnection(cc.Provider, cc.Token, cc.Connection)
	if err != nil {
		return nil, err
	}
	client, err := k8sContext.GenerateKubeHandler()
	if err != nil {
		return nil, err
	}
	return &driver{clientset: client.KubeClient, restConfig: client.RestConfig}, nil
}

// driver serves consoles against one Kubernetes cluster.
//
// It holds kubernetes.Interface rather than the concrete clientset so that the
// resource lookups behind Capabilities can be exercised against a fake.
type driver struct {
	clientset  k8sclient.Interface
	restConfig rest.Config
}

var (
	_ console.TerminalDriver = (*driver)(nil)
	_ console.LogStreamer    = (*driver)(nil)
)

// Capabilities inspects the live pod to decide what it admits.
func (d *driver) Capabilities(ctx context.Context, target console.TargetRef) (console.Capabilities, error) {
	_, caps, err := d.attach(ctx, target)
	return caps, err
}

// attach resolves the target onto the concrete pod a console attaches to, and
// reports what that pod admits. Capabilities, OpenTerminal and StreamLogs all go
// through it, so a workload means the same pod on the handshake as it does on the
// exec or log request.
//
// A resource the driver cannot address at all — a ConfigMap, a Service — is not
// an error: it answers with no supported kind and the reason, which the panel
// renders. Only a target that does not exist is an error, so a stale row in the
// UI is reported as missing rather than as unsupported.
func (d *driver) attach(ctx context.Context, target console.TargetRef) (*corev1.Pod, console.Capabilities, error) {
	if target.Name == "" || target.Namespace == "" {
		return nil, console.Capabilities{}, console.ErrInvalidTarget("a Kubernetes console target requires both `namespace` and `name`")
	}

	resource := normalizeResource(target.Resource)
	if !isConsoleResource(resource) {
		return nil, console.Capabilities{
			Reason: "Terminal and log consoles attach to a pod; '" + target.Resource + "' does not run containers of its own.",
		}, nil
	}

	pod, err := d.podFor(ctx, resource, target)
	if err != nil {
		return nil, console.Capabilities{}, err
	}
	if pod == nil {
		return nil, console.Capabilities{
			Reason: "The " + resource + " '" + target.Name + "' has no pods to attach to.",
		}, nil
	}

	caps := console.Capabilities{
		// Logs stay available for a pod in any phase: a Pending pod may have a
		// running init container worth tailing, and a Failed pod's logs are
		// precisely what a user comes for. When there is genuinely nothing to
		// read the API server says so, and that error is more informative than
		// a capability flag would be.
		Logs:       true,
		Containers: containerNames(pod),
	}
	caps.DefaultContainer = defaultContainer(pod, caps.Containers)

	// Exec needs a process to attach to, which only a running pod has.
	caps.Terminal = pod.Status.Phase == corev1.PodRunning
	if !caps.Terminal {
		caps.Reason = "The pod is in the " + string(pod.Status.Phase) + " phase; an interactive terminal requires a running pod."
	}

	return pod, caps, nil
}

// podFor resolves a console target onto one pod. A pod target is used as named;
// an owner workload is found through the label selector it stamps on its pods.
//
// A nil pod with no error means the resource exists but runs nothing, which the
// caller reports as an unsupported console rather than a missing target.
func (d *driver) podFor(ctx context.Context, resource string, target console.TargetRef) (*corev1.Pod, error) {
	if resource == ResourcePod {
		pod, err := d.clientset.CoreV1().Pods(target.Namespace).Get(ctx, target.Name, metav1.GetOptions{})
		if err != nil {
			return nil, resourceError(err, target)
		}
		return pod, nil
	}

	selector, err := d.workloadSelector(ctx, resource, target)
	if err != nil {
		return nil, err
	}
	if selector == "" {
		// An owner resource with no selector names no pods. Listing everything in
		// the namespace would attach the console to somebody else's workload.
		return nil, nil
	}

	pods, err := d.clientset.CoreV1().Pods(target.Namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, resourceError(err, target)
	}
	if len(pods.Items) == 0 {
		return nil, nil
	}
	return selectPod(pods.Items), nil
}

// workloadSelector returns the label selector the owner workload uses to name its
// pods. Fetching the owner is not incidental: it is what separates "this
// deployment has no pods" from "this deployment does not exist".
func (d *driver) workloadSelector(ctx context.Context, resource string, target console.TargetRef) (string, error) {
	opts := metav1.GetOptions{}
	switch resource {
	case ResourceDeployment:
		deployment, err := d.clientset.AppsV1().Deployments(target.Namespace).Get(ctx, target.Name, opts)
		if err != nil {
			return "", resourceError(err, target)
		}
		return selectorString(deployment.Spec.Selector)
	case ResourceStatefulSet:
		statefulSet, err := d.clientset.AppsV1().StatefulSets(target.Namespace).Get(ctx, target.Name, opts)
		if err != nil {
			return "", resourceError(err, target)
		}
		return selectorString(statefulSet.Spec.Selector)
	case ResourceDaemonSet:
		daemonSet, err := d.clientset.AppsV1().DaemonSets(target.Namespace).Get(ctx, target.Name, opts)
		if err != nil {
			return "", resourceError(err, target)
		}
		return selectorString(daemonSet.Spec.Selector)
	case ResourceReplicaSet:
		replicaSet, err := d.clientset.AppsV1().ReplicaSets(target.Namespace).Get(ctx, target.Name, opts)
		if err != nil {
			return "", resourceError(err, target)
		}
		return selectorString(replicaSet.Spec.Selector)
	case ResourceJob:
		job, err := d.clientset.BatchV1().Jobs(target.Namespace).Get(ctx, target.Name, opts)
		if err != nil {
			return "", resourceError(err, target)
		}
		return selectorString(job.Spec.Selector)
	case ResourceReplicationController:
		controller, err := d.clientset.CoreV1().ReplicationControllers(target.Namespace).Get(ctx, target.Name, opts)
		if err != nil {
			return "", resourceError(err, target)
		}
		return selectorString(&metav1.LabelSelector{MatchLabels: controller.Spec.Selector})
	default:
		return "", console.ErrInvalidTarget("unsupported console resource '" + target.Resource + "'")
	}
}

// selectPod picks the pod a console attaches to when an owner workload runs
// several. Running pods win, so a terminal finds a process to attach to; among
// them the most recently created wins, because after a rollout the pods created
// last carry the current spec while the older ones are draining. Ties are broken
// by name, so two consoles opened against one workload land on the same pod
// rather than on whichever the API server happened to return first.
func selectPod(pods []corev1.Pod) *corev1.Pod {
	best := &pods[0]
	bestRunning := best.Status.Phase == corev1.PodRunning
	for i := 1; i < len(pods); i++ {
		candidate := &pods[i]
		running := candidate.Status.Phase == corev1.PodRunning
		switch {
		case running && !bestRunning:
			best, bestRunning = candidate, true
		case running == bestRunning && isNewerPod(candidate, best):
			best = candidate
		}
	}
	return best
}

func isNewerPod(candidate, existing *corev1.Pod) bool {
	if !candidate.CreationTimestamp.Equal(&existing.CreationTimestamp) {
		return candidate.CreationTimestamp.After(existing.CreationTimestamp.Time)
	}
	return candidate.Name < existing.Name
}

// selectorString renders a label selector as the query string a pod list expects.
// A nil or empty selector yields "", which podFor reads as "this workload names
// no pods": listing the whole namespace would attach the console to somebody
// else's workload.
func selectorString(selector *metav1.LabelSelector) (string, error) {
	if selector == nil {
		return "", nil
	}
	// LabelSelectorAsSelector normalises matchExpressions into the `key in (a,b)`
	// form the selector grammar accepts. LabelSelector.String() is protobuf's debug
	// rendering of the struct, not a selector.
	parsed, err := metav1.LabelSelectorAsSelector(selector)
	if err != nil {
		return "", console.ErrOpenConsole(err, "capability")
	}
	return parsed.String(), nil
}

// resourceError maps a Kubernetes API failure onto a coded console error. A
// missing resource is a 404 the user needs to hear about; anything else is a
// failure reaching the target behind the console.
func resourceError(err error, target console.TargetRef) error {
	if apierrors.IsNotFound(err) {
		return console.ErrTargetNotFound(err, target.Namespace+"/"+target.Name)
	}
	return console.ErrOpenConsole(err, "capability")
}

// OpenTerminal execs an interactive shell in the target container.
func (d *driver) OpenTerminal(ctx context.Context, target console.TargetRef, opts console.TerminalOptions, stream console.TerminalIO) error {
	pod, caps, err := d.attach(ctx, target)
	if err != nil {
		return err
	}
	if !console.Supports(caps, console.KindTerminal) {
		return console.ErrUnsupportedConsole(string(console.KindTerminal), target.Resource, caps.Reason)
	}
	container, err := resolveContainer(target, caps)
	if err != nil {
		return err
	}

	command := opts.Command
	if len(command) == 0 {
		command = defaultShellCommand
	}

	req := d.clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Namespace(target.Namespace).
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdin:     true,
			Stdout:    true,
			// With a TTY the kernel line discipline merges the two output
			// streams onto the pty, and the API server rejects a request that
			// asks for a separate stderr alongside TTY.
			Stderr: false,
			TTY:    true,
		}, scheme.ParameterCodec)

	executor, err := d.newExecutor(req.URL())
	if err != nil {
		return console.ErrOpenConsole(err, string(console.KindTerminal))
	}

	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             stream.Stdin,
		Stdout:            stream.Stdout,
		Tty:               true,
		TerminalSizeQueue: &resizeQueue{ctx: ctx, sizes: stream.Resize},
	})
	if err != nil {
		// A shell that exits non-zero (a failed command, or Ctrl-D after one)
		// ended the console normally. Report the status rather than dressing a
		// routine exit up as a transport failure.
		var codeErr k8sexec.CodeExitError
		if errors.As(err, &codeErr) {
			return &console.ExitError{Code: codeErr.Code, Err: err}
		}
		// The client hanging up cancels ctx, which surfaces here as a context
		// error. That is the expected end of a console, not a failure.
		if ctx.Err() != nil {
			return nil
		}
		return console.ErrOpenConsole(err, string(console.KindTerminal))
	}
	return nil
}

// StreamLogs opens the target container's log stream.
func (d *driver) StreamLogs(ctx context.Context, target console.TargetRef, opts console.LogOptions) (io.ReadCloser, error) {
	pod, caps, err := d.attach(ctx, target)
	if err != nil {
		return nil, err
	}
	if !console.Supports(caps, console.KindLogs) {
		return nil, console.ErrUnsupportedConsole(string(console.KindLogs), target.Resource, caps.Reason)
	}
	container, err := resolveContainer(target, caps)
	if err != nil {
		return nil, err
	}

	req := d.clientset.CoreV1().Pods(target.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
		Container:    container,
		Follow:       opts.Follow,
		Previous:     opts.Previous,
		Timestamps:   opts.Timestamps,
		TailLines:    opts.TailLines,
		SinceSeconds: opts.SinceSeconds,
	})

	stream, err := req.Stream(ctx)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, console.ErrTargetNotFound(err, target.Namespace+"/"+target.Name)
		}
		return nil, console.ErrOpenConsole(err, string(console.KindLogs))
	}
	return stream, nil
}

// newExecutor prefers the WebSocket exec protocol and falls back to SPDY.
//
// The API server has spoken WebSocket for exec since v1.29, but SPDY is what
// older clusters — and intermediate proxies that only learned SPDY — still
// require. The fallback executor tries WebSocket and retries over SPDY when the
// upgrade is refused, which is what kubectl does.
func (d *driver) newExecutor(u *url.URL) (remotecommand.Executor, error) {
	// A copy, because the streaming config must not disturb the shared client's
	// config. Timeout is zeroed: it bounds a whole request, and an interactive
	// terminal is a request that is meant to last as long as the user keeps it
	// open.
	cfg := rest.CopyConfig(&d.restConfig)
	cfg.Timeout = 0

	websocketExec, err := remotecommand.NewWebSocketExecutor(cfg, "GET", u.String())
	if err != nil {
		return nil, err
	}
	spdyExec, err := remotecommand.NewSPDYExecutor(cfg, "POST", u)
	if err != nil {
		return nil, err
	}
	return remotecommand.NewFallbackExecutor(websocketExec, spdyExec, func(err error) bool {
		return httpstream.IsUpgradeFailure(err) || httpstream.IsHTTPSProxyError(err)
	})
}

// normalizeResource canonicalises the resource kind a target names. MeshSync
// reports Kubernetes kinds in PascalCase ("Deployment"); the driver's constants
// are lowercase.
func normalizeResource(resource string) string {
	return strings.ToLower(strings.TrimSpace(resource))
}

// isConsoleResource reports whether this driver can address a resource kind at
// all, before any of it is looked up in the cluster.
func isConsoleResource(resource string) bool {
	switch resource {
	case ResourcePod, ResourceDeployment, ResourceStatefulSet, ResourceDaemonSet, ResourceJob, ResourceReplicaSet, ResourceReplicationController:
		return true
	default:
		return false
	}
}

// resolveContainer picks the container a console attaches to, and refuses a
// name the pod does not have rather than letting the API server return an
// opaque error mid-handshake.
func resolveContainer(target console.TargetRef, caps console.Capabilities) (string, error) {
	if target.Container == "" {
		if caps.DefaultContainer == "" {
			return "", console.ErrInvalidTarget("the pod exposes no containers to attach to")
		}
		return caps.DefaultContainer, nil
	}
	for _, name := range caps.Containers {
		if name == target.Container {
			return name, nil
		}
	}
	return "", console.ErrInvalidTarget("the pod has no container named '" + target.Container + "'; it has: " + strings.Join(caps.Containers, ", "))
}

// containerNames lists every addressable container in the pod. Init and
// ephemeral containers are included: init container logs are often the reason a
// pod is being inspected, and an ephemeral debug container is exactly the thing
// a user wants a terminal into.
func containerNames(pod *corev1.Pod) []string {
	names := make([]string, 0, len(pod.Spec.Containers)+len(pod.Spec.InitContainers)+len(pod.Spec.EphemeralContainers))
	for _, c := range pod.Spec.Containers {
		names = append(names, c.Name)
	}
	for _, c := range pod.Spec.InitContainers {
		names = append(names, c.Name)
	}
	for _, c := range pod.Spec.EphemeralContainers {
		names = append(names, c.Name)
	}
	return names
}

// defaultContainer honours the kubectl default-container annotation when it
// names a container that exists, and otherwise falls back to the first regular
// container — the same precedence kubectl applies.
func defaultContainer(pod *corev1.Pod, names []string) string {
	if annotated := pod.Annotations[defaultContainerAnnotation]; annotated != "" {
		for _, name := range names {
			if name == annotated {
				return annotated
			}
		}
	}
	if len(pod.Spec.Containers) > 0 {
		return pod.Spec.Containers[0].Name
	}
	if len(names) > 0 {
		return names[0]
	}
	return ""
}

// resizeQueue adapts the transport's resize channel to the shape client-go
// wants. Returning nil tells client-go to stop polling, which it does when the
// channel closes at console end or the context is cancelled.
type resizeQueue struct {
	ctx   context.Context
	sizes <-chan console.TerminalSize
}

func (q *resizeQueue) Next() *remotecommand.TerminalSize {
	select {
	case size, ok := <-q.sizes:
		if !ok {
			return nil
		}
		return &remotecommand.TerminalSize{Width: size.Width, Height: size.Height}
	case <-q.ctx.Done():
		return nil
	}
}
