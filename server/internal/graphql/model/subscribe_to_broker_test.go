package model_test

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshkit/broker/nats"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	"k8s.io/client-go/rest"
)

// fakeNATS speaks just enough of the NATS client protocol for nats.go to
// connect, and records every SUB and PUB subject it sees.
type fakeNATS struct {
	ln   net.Listener
	mu   sync.Mutex
	subs []string
	pubs []string
}

func newFakeNATS(t *testing.T) *fakeNATS {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	f := &fakeNATS{ln: ln}
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go f.serve(conn)
		}
	}()
	return f
}

func (f *fakeNATS) addr() string { return f.ln.Addr().String() }

func (f *fakeNATS) serve(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	port := f.ln.Addr().(*net.TCPAddr).Port
	_, _ = fmt.Fprintf(conn, "INFO {\"server_id\":\"fake\",\"version\":\"2.10.0\",\"host\":\"127.0.0.1\",\"port\":%d,\"max_payload\":1048576,\"proto\":1}\r\n", port)
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch strings.ToUpper(fields[0]) {
		case "PING":
			_, _ = conn.Write([]byte("PONG\r\n"))
		case "SUB":
			f.mu.Lock()
			f.subs = append(f.subs, fields[1])
			f.mu.Unlock()
		case "PUB":
			size, _ := strconv.Atoi(fields[len(fields)-1])
			payload := make([]byte, size+2) // payload plus trailing \r\n
			if _, err := io.ReadFull(r, payload); err != nil {
				return
			}
			f.mu.Lock()
			f.pubs = append(f.pubs, fields[1])
			f.mu.Unlock()
		}
	}
}

func (f *fakeNATS) seen() (subs, pubs []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.subs...), append([]string(nil), f.pubs...)
}

// fakeBrokerAPI serves the meshery-broker custom resource with both
// endpoints pointing at the fake NATS server.
func fakeBrokerAPI(t *testing.T, endpoint string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apis/meshery.io/v1alpha1/namespaces/meshery/brokers/meshery-broker" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"apiVersion":"meshery.io/v1alpha1","kind":"Broker","metadata":{"name":"meshery-broker","namespace":"meshery"},"status":{"endpoint":{"internal":%q,"external":%q}}}`, endpoint, endpoint)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// SubscribeToBroker used to queue-subscribe meshery.meshsync.core into a
// channel that was never created or read, so the first MeshSync message
// blocked that subscription forever. It should only connect and ask MeshSync
// to resync.
func TestSubscribeToBroker_RequestsResyncWithoutSubscribing(t *testing.T) {
	natsSrv := newFakeNATS(t)
	api := fakeBrokerAPI(t, natsSrv.addr())

	kubeClient := &mesherykube.Client{RestConfig: rest.Config{Host: api.URL}}
	brokerConn := &nats.Nats{}
	t.Cleanup(brokerConn.CloseConnection)

	endpoint, err := model.SubscribeToBroker(nil, kubeClient, brokerConn, model.NewK8sConnctionTracker())
	if err != nil {
		t.Fatalf("SubscribeToBroker returned error: %v", err)
	}
	if endpoint != natsSrv.addr() {
		t.Fatalf("endpoint = %q, want %q", endpoint, natsSrv.addr())
	}

	deadline := time.Now().Add(5 * time.Second)
	var subs, pubs []string
	for time.Now().Before(deadline) {
		subs, pubs = natsSrv.seen()
		if len(pubs) > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if len(pubs) != 1 || pubs[0] != model.RequestSubject {
		t.Errorf("published subjects = %v, want [%s]", pubs, model.RequestSubject)
	}
	if len(subs) != 0 {
		t.Errorf("expected no subscriptions, got %v", subs)
	}
}
