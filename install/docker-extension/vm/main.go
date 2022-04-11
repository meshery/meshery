package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/docker/meshery-extension/vm/pkg/socket"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func main() {
	var socketPath = flag.String("socket", "/run/guest/volumes-service.sock", "Unix domain socket to listen on")
	var testPort = flag.Int("simplePort", 0, "Test port to expose instead of socket")
	flag.Parse()
	unixSocket := "unix:" + *socketPath
	logrus.New().Infof("Starting listening on %s\n", unixSocket)

	router := echo.New()
	router.HideBanner = true

	startURL := ""

	if *testPort != 0 {
		startURL = fmt.Sprintf(":%d", *testPort)
	} else {
		ln, err := socket.ListenOn(unixSocket)
		if err != nil {
			log.Fatal(err)
		}
		router.Listener = ln
	}

	router.GET("/ping", func(c echo.Context) error {
		resp, err := http.Get("http://host.docker.internal:9081/api/system/version")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		return c.String(http.StatusOK, string(body))
	})

	log.Fatal(router.Start(startURL))
}
