// Package kubernetes implements the sessions driver for connections of kind
// "kubernetes": interactive terminals backed by the pods/exec subresource, and
// log tails backed by pods/log.
//
// Importing this package registers the driver with the sessions registry. The
// HTTP layer never imports it directly; it is pulled in for its side effect
// wherever sessions are served.
package kubernetes

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strings"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/sessions"
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

	// ResourcePod is the only resource type that supports sessions today.
	// Terminals and logs are properties of a running container, and a pod is
	// the smallest thing that has one. Owner workloads (Deployment,
	// StatefulSet, ...) would be served by resolving them to pods; that
	// resolution belongs here, behind Capabilities, when it is added.
	ResourcePod = "pod"

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
	sessions.Register(factory{})
}

// factory builds a driver per Kubernetes connection.
type factory struct{}

func (factory) Kind() string { return ConnectionKind }

// NewDriver resolves the connection's credential into a Kubernetes client.
//
// The credential is fetched with the requesting user's token, so the session
// runs as that user against the cluster and Kubernetes RBAC governs what it may
// do. Meshery adds no authorization of its own beyond deciding who may reach
// the connection: a user who cannot `create pods/exec` gets a 403 from the API
// server, exactly as they would with kubectl.
func (factory) NewDriver(_ context.Context, cc sessions.ConnectionContext) (sessions.Driver, error) {
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

// driver serves sessions against one Kubernetes cluster.
//
// It holds kubernetes.Interface rather than the concrete clientset so that the
// resource lookups behind Capabilities can be exercised against a fake.
type driver struct {
	clientset  k8sclient.Interface
	restConfig rest.Config
}

var (
	_ sessions.TerminalDriver = (*driver)(nil)
	_ sessions.LogStreamer    = (*driver)(nil)
)

// Capabilities inspects the live pod to decide what it admits.
func (d *driver) Capabilities(ctx context.Context, target sessions.TargetRef) (sessions.Capabilities, error) {
	if !strings.EqualFold(target.Resource, ResourcePod) {
		return sessions.Capabilities{
			Reason: "Terminal and log sessions are available for pods; '" + target.Resource + "' resources do not run containers directly.",
		}, nil
	}
	if target.Name == "" || target.Namespace == "" {
		return sessions.Capabilities{}, sessions.ErrInvalidTarget("a pod session requires both `namespace` and `name`")
	}

	pod, err := d.clientset.CoreV1().Pods(target.Namespace).Get(ctx, target.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return sessions.Capabilities{}, sessions.ErrTargetNotFound(err, target.Namespace+"/"+target.Name)
		}
		return sessions.Capabilities{}, sessions.ErrOpenSession(err, "capability")
	}

	caps := sessions.Capabilities{
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

	return caps, nil
}

// OpenTerminal execs an interactive shell in the target container.
func (d *driver) OpenTerminal(ctx context.Context, target sessions.TargetRef, opts sessions.TerminalOptions, stream sessions.TerminalIO) error {
	caps, err := d.Capabilities(ctx, target)
	if err != nil {
		return err
	}
	if !sessions.Supports(caps, sessions.KindTerminal) {
		return sessions.ErrUnsupportedSession(string(sessions.KindTerminal), target.Resource, caps.Reason)
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
		Name(target.Name).
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
		return sessions.ErrOpenSession(err, string(sessions.KindTerminal))
	}

	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             stream.Stdin,
		Stdout:            stream.Stdout,
		Tty:               true,
		TerminalSizeQueue: &resizeQueue{ctx: ctx, sizes: stream.Resize},
	})
	if err != nil {
		// A shell that exits non-zero (a failed command, or Ctrl-D after one)
		// ended the session normally. Report the status rather than dressing a
		// routine exit up as a transport failure.
		var codeErr k8sexec.CodeExitError
		if errors.As(err, &codeErr) {
			return &sessions.ExitError{Code: codeErr.Code, Err: err}
		}
		// The client hanging up cancels ctx, which surfaces here as a context
		// error. That is the expected end of a session, not a failure.
		if ctx.Err() != nil {
			return nil
		}
		return sessions.ErrOpenSession(err, string(sessions.KindTerminal))
	}
	return nil
}

// StreamLogs opens the target container's log stream.
func (d *driver) StreamLogs(ctx context.Context, target sessions.TargetRef, opts sessions.LogOptions) (io.ReadCloser, error) {
	caps, err := d.Capabilities(ctx, target)
	if err != nil {
		return nil, err
	}
	if !sessions.Supports(caps, sessions.KindLogs) {
		return nil, sessions.ErrUnsupportedSession(string(sessions.KindLogs), target.Resource, caps.Reason)
	}
	container, err := resolveContainer(target, caps)
	if err != nil {
		return nil, err
	}

	req := d.clientset.CoreV1().Pods(target.Namespace).GetLogs(target.Name, &corev1.PodLogOptions{
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
			return nil, sessions.ErrTargetNotFound(err, target.Namespace+"/"+target.Name)
		}
		return nil, sessions.ErrOpenSession(err, string(sessions.KindLogs))
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

// resolveContainer picks the container a session attaches to, and refuses a
// name the pod does not have rather than letting the API server return an
// opaque error mid-handshake.
func resolveContainer(target sessions.TargetRef, caps sessions.Capabilities) (string, error) {
	if target.Container == "" {
		if caps.DefaultContainer == "" {
			return "", sessions.ErrInvalidTarget("the pod exposes no containers to attach to")
		}
		return caps.DefaultContainer, nil
	}
	for _, name := range caps.Containers {
		if name == target.Container {
			return name, nil
		}
	}
	return "", sessions.ErrInvalidTarget("the pod has no container named '" + target.Container + "'; it has: " + strings.Join(caps.Containers, ", "))
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
// channel closes at session end or the context is cancelled.
type resizeQueue struct {
	ctx   context.Context
	sizes <-chan sessions.TerminalSize
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
