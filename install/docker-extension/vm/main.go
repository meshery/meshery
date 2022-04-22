package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/docker/meshery-extension/vm/pkg/socket"
	"github.com/docker/meshery-extension/vm/proxy"
	"github.com/sirupsen/logrus"
)

var (
	MesheryServerHost = "http://host.docker.internal:9081"
)

// func main() {
// 	var socketPath = flag.String("socket", "/run/guest/volumes-service.sock", "Unix domain socket to listen on")
// 	var testPort = flag.Int("simplePort", 0, "Test port to expose instead of socket")
// 	flag.Parse()
// 	unixSocket := "unix:" + *socketPath
// 	logrus.New().Infof("Starting listening on %s\n", unixSocket)

// 	router := echo.New()
// 	router.HideBanner = true

// 	startURL := ""

// 	if *testPort != 0 {
// 		startURL = fmt.Sprintf(":%d", *testPort)
// 	} else {
// 		ln, err := socket.ListenOn(unixSocket)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		router.Listener = ln
// 	}

// 	router.GET("/ping", func(c echo.Context) error {
// 		resp, err := http.Get(MesheryServerHost + "/api/system/version")
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return err
// 		}
// 		defer resp.Body.Close()
// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		return c.String(http.StatusOK, string(body))
// 	})

// 	router.POST("/mesh/deploy", func(c echo.Context) error {
// 		meshLocationURL := c.FormValue("meshLocationURL")

// 		data := url.Values{}
// 		data.Set("meshLocationURL", meshLocationURL)

// 		u, _ := url.ParseRequestURI(MesheryServerHost)
// 		u.Path = "/api/system/adapter/manage"
// 		urlStr := u.String()

// 		client := &http.Client{}
// 		r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
// 		r.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

// 		resp, err := client.Do(r)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return err
// 		}
// 		defer resp.Body.Close()
// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		fmt.Println(resp.Status)
// 		return c.String(http.StatusOK, string(body))
// 	})

// 	log.Fatal(router.Start(startURL))
// }

func main() {

	var socketPath = flag.String("socket", "/run/guest/volumes-service.sock", "Unix domain socket to listen on")
	flag.Parse()
	unixSocket := "unix:" + *socketPath
	ln, err := socket.ListenOn(unixSocket)
	if err != nil {
		log.Fatal(err)
	}

	handler := &proxy.Proxy{}
	logrus.New().Infof("Starting listening on %s\n", unixSocket)
	if err := http.Serve(ln, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
