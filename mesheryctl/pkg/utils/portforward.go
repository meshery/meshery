package utils

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// PortForwarder is a type that implements port forwarding to a pod
type PortForwarder struct {
	forwarder *portforward.PortForwarder
	stopChan  chan struct{}
	readyChan <-chan struct{}
	sigChan   chan os.Signal
	done      chan struct{}
}

// NewPortForwarder creates a new port forwarder to a pod
func NewPortForwarder(dialer httpstream.Dialer, portSpec string) (*PortForwarder, error) {
	stopChan := make(chan struct{})
	readyChan := make(chan struct{})
	forwarder, err := portforward.New(dialer, []string{portSpec}, stopChan, readyChan, ioutil.Discard, os.Stderr)
	if err != nil {
		return nil, errors.Errorf("Error setting up port forwarding: %s", err)
	}

	return &PortForwarder{
		forwarder: forwarder,
		stopChan:  stopChan,
		readyChan: readyChan,
		done:      make(chan struct{}),
	}, nil
}

// Start starts the port forwarding and calls the readyFunc callback function when port forwarding is ready
func (pf *PortForwarder) Start(readyFunc func(pf *PortForwarder) error) error {
	// Set up a channel to process OS signals
	pf.sigChan = make(chan os.Signal, 1)
	signal.Notify(pf.sigChan, os.Interrupt)

	// Set up a channel to process forwarding errors
	errChan := make(chan error, 1)

	// Set up forwarding, and relay errors to errChan so that caller is able to handle it
	go func() {
		err := pf.forwarder.ForwardPorts()
		errChan <- err
	}()

	// Stop forwarding in case OS signals are received
	go func() {
		<-pf.sigChan
		pf.Stop()
	}()

	// Process signals and call readyFunc
	select {
	case <-pf.readyChan:
		return readyFunc(pf)

	case err := <-errChan:
		return errors.Errorf("Error during port forwarding: %s", err)
	}
}

// Stop stops the port forwarding if not stopped already
func (pf *PortForwarder) Stop() {
	defer close(pf.done)
	signal.Stop(pf.sigChan)
	if pf.stopChan != nil {
		close(pf.stopChan)
	}
}

// Done returns a channel that is closed after Stop has been called.
func (pf *PortForwarder) Done() <-chan struct{} {
	return pf.done
}

// DialerToPod constructs a new httpstream.Dialer to connect to a pod for use
// with a PortForwarder
func DialerToPod(conf *rest.Config, clientSet kubernetes.Interface, podName string, namespace string) (httpstream.Dialer, error) {
	roundTripper, upgrader, err := spdy.RoundTripperFor(conf)
	if err != nil {
		return nil, errors.Errorf("Error setting up round tripper for port forwarding: %s", err)
	}

	serverURL := clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("portforward").URL()

	return spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, serverURL), nil
}
