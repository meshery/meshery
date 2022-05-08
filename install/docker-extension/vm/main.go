package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	// "github.com/docker/meshery-extension/vm/pkg/socket"
	"github.com/docker/meshery-extension/vm/proxy"
	"github.com/sirupsen/logrus"
)

var (
	MesheryServerHost = "http://host.docker.internal:9081"
)

func main() {

	var socketPath = flag.String("socket", "/run/guest/volumes-service.sock", "Unix domain socket to listen on")
	flag.Parse()
	unixSocket := "unix:" + *socketPath
	// ln, err := socket.ListenOn(unixSocket)
	ln, err := net.Listen("tcp", "localhost:7877")
	if err != nil {
		log.Fatal(err)
	}

	handler := &proxy.Proxy{}
	logrus.New().Infof("Starting listening on %s\n", unixSocket)
	if err := http.Serve(ln, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
