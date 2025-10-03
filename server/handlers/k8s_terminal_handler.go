// Package handlers : collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/meshery/meshery/server/models"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Configure proper CORS in production
		return true
	},
}

// TerminalMessage represents WebSocket messages for terminal
type TerminalMessage struct {
	Op   string          `json:"op"`   // stdin, stdout, stderr, resize, ping
	Data string          `json:"data"` // actual data
	Rows uint16          `json:"rows,omitempty"`
	Cols uint16          `json:"cols,omitempty"`
}

// TerminalSession handles WebSocket terminal I/O
type TerminalSession struct {
	conn     *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

// NewTerminalSession creates a new terminal session
func NewTerminalSession(conn *websocket.Conn) *TerminalSession {
	return &TerminalSession{
		conn:     conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
}

// Read implements io.Reader interface
func (t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.conn.ReadMessage()
	if err != nil {
		return 0, err
	}

	var msg TerminalMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return 0, err
	}

	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{
			Width:  msg.Cols,
			Height: msg.Rows,
		}
		return 0, nil
	}

	return 0, nil
}

// Write implements io.Writer interface
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg := TerminalMessage{
		Op:   "stdout",
		Data: string(p),
	}

	err := t.conn.WriteJSON(msg)
	return len(p), err
}

// Next implements remotecommand.TerminalSizeQueue
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

// Close closes the terminal session
func (t *TerminalSession) Close() {
	close(t.doneChan)
}

// swagger:route GET /api/system/kubernetes/exec KubernetesAPI idK8sExec
// Handle WebSocket connection for pod exec
//
// Creates an interactive terminal session to a Kubernetes pod
//
// responses:
//	101: Switching Protocols

// K8sExecHandler handles WebSocket connections for pod exec
func (h *Handler) K8sExecHandler(w http.ResponseWriter, r *http.Request,
	prefObj *models.Preference, user *models.User, provider models.Provider) {

	// Parse query parameters
	namespace := r.URL.Query().Get("namespace")
	podName := r.URL.Query().Get("pod")
	containerName := r.URL.Query().Get("container")
	contextID := r.URL.Query().Get("context")
	shell := r.URL.Query().Get("shell")

	if namespace == "" || podName == "" || contextID == "" {
		h.log.Error(ErrInvalidK8sExecParams())
		http.Error(w, "Missing required parameters: namespace, pod, context", http.StatusBadRequest)
		return
	}

	// Default shell command
	if shell == "" {
		shell = "/bin/sh"
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error(ErrUpgradingWebSocket(err))
		return
	}
	defer conn.Close()

	// Get Kubernetes client for the specified context
	k8sClient, err := h.getK8sClientFromContextID(r.Context(), contextID, provider, user)
	if err != nil {
		h.log.Error(ErrGettingK8sClient(err))
		_ = conn.WriteJSON(TerminalMessage{
			Op:   "stdout",
			Data: fmt.Sprintf("\r\nError: Failed to connect to Kubernetes cluster: %v\r\n", err),
		})
		return
	}

	// Construct the shell command with fallback
	cmd := []string{
		shell, "-c",
		"TERM=xterm-256color; export TERM; " +
			"[ -x /bin/bash ] && " +
			"([ -x /usr/bin/script ] && /usr/bin/script -q -c \"/bin/bash\" /dev/null || exec /bin/bash) || " +
			"exec /bin/sh",
	}

	// Create pod exec request
	req := k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	execOptions := &corev1.PodExecOptions{
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
		Command:   cmd,
	}

	if containerName != "" {
		execOptions.Container = containerName
	}

	req.VersionedParams(execOptions, scheme.ParameterCodec)

	// Create SPDY executor
	executor, err := remotecommand.NewSPDYExecutor(k8sClient.Config, "POST", req.URL())
	if err != nil {
		h.log.Error(ErrCreatingSPDYExecutor(err))
		_ = conn.WriteJSON(TerminalMessage{
			Op:   "stdout",
			Data: fmt.Sprintf("\r\nError: Failed to create executor: %v\r\n", err),
		})
		return
	}

	// Create terminal session
	terminalSession := NewTerminalSession(conn)
	defer terminalSession.Close()

	// Send welcome message
	_ = conn.WriteJSON(TerminalMessage{
		Op:   "stdout",
		Data: fmt.Sprintf("\r\nConnected to pod: %s/%s\r\n", namespace, podName),
	})

	// Start streaming with context
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Handle ping/pong for connection keep-alive
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteJSON(TerminalMessage{Op: "ping"}); err != nil {
					cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Stream the connection
	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             terminalSession,
		Stdout:            terminalSession,
		Stderr:            terminalSession,
		Tty:               true,
		TerminalSizeQueue: terminalSession,
	})

	if err != nil {
		h.log.Error(ErrStreamingExec(err))
		_ = conn.WriteJSON(TerminalMessage{
			Op:   "stdout",
			Data: fmt.Sprintf("\r\nConnection closed: %v\r\n", err),
		})
	}
}

// swagger:route GET /api/system/kubernetes/logs KubernetesAPI idK8sLogs
// Handle WebSocket connection for pod logs
//
// Streams logs from a Kubernetes pod
//
// responses:
//	101: Switching Protocols

// K8sLogsHandler handles WebSocket connections for pod logs streaming
func (h *Handler) K8sLogsHandler(w http.ResponseWriter, r *http.Request,
	prefObj *models.Preference, user *models.User, provider models.Provider) {

	// Parse query parameters
	namespace := r.URL.Query().Get("namespace")
	podName := r.URL.Query().Get("pod")
	containerName := r.URL.Query().Get("container")
	contextID := r.URL.Query().Get("context")
	follow := r.URL.Query().Get("follow") == "true"
	previous := r.URL.Query().Get("previous") == "true"
	tailLines := int64(100) // Default to last 100 lines

	if namespace == "" || podName == "" || contextID == "" {
		h.log.Error(ErrInvalidK8sLogsParams())
		http.Error(w, "Missing required parameters: namespace, pod, context", http.StatusBadRequest)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error(ErrUpgradingWebSocket(err))
		return
	}
	defer conn.Close()

	// Get Kubernetes client
	k8sClient, err := h.getK8sClientFromContextID(r.Context(), contextID, provider, user)
	if err != nil {
		h.log.Error(ErrGettingK8sClient(err))
		_ = conn.WriteJSON(TerminalMessage{
			Op:   "stdout",
			Data: fmt.Sprintf("Error: Failed to connect to Kubernetes cluster: %v\n", err),
		})
		return
	}

	// Create pod log options
	logOptions := &corev1.PodLogOptions{
		Follow:    follow,
		Previous:  previous,
		TailLines: &tailLines,
	}

	if containerName != "" {
		logOptions.Container = containerName
	}

	// Get log stream
	req := k8sClient.CoreV1().Pods(namespace).GetLogs(podName, logOptions)
	stream, err := req.Stream(r.Context())
	if err != nil {
		h.log.Error(ErrGettingPodLogs(err))
		_ = conn.WriteJSON(TerminalMessage{
			Op:   "stdout",
			Data: fmt.Sprintf("Error: Failed to get logs: %v\n", err),
		})
		return
	}
	defer stream.Close()

	// Send initial message
	_ = conn.WriteJSON(TerminalMessage{
		Op:   "stdout",
		Data: fmt.Sprintf("=== Logs from %s/%s ===\n", namespace, podName),
	})

	// Stream logs to WebSocket
	buf := make([]byte, 2048)
	for {
		n, err := stream.Read(buf)
		if n > 0 {
			if writeErr := conn.WriteJSON(TerminalMessage{
				Op:   "stdout",
				Data: string(buf[:n]),
			}); writeErr != nil {
				h.log.Error(ErrWritingToWebSocket(writeErr))
				return
			}
		}
		if err != nil {
			if err != io.EOF {
				h.log.Error(ErrReadingLogs(err))
			}
			break
		}
	}
}

// Helper function to get Kubernetes client from context ID
func (h *Handler) getK8sClientFromContextID(ctx context.Context, contextID string,
	provider models.Provider, user *models.User) (*kubernetes.Clientset, error) {

	// Get the K8s context from database
	k8sContext, err := provider.GetK8sContext(contextID, user.ID)
	if err != nil {
		return nil, ErrGettingK8sContext(err)
	}

	// Get Kubernetes client config
	config, err := k8sContext.GenerateKubeConfig()
	if err != nil {
		return nil, ErrGeneratingKubeConfig(err)
	}

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, ErrCreatingK8sClient(err)
	}

	return clientset, nil
}

