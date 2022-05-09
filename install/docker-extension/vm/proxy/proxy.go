package proxy

import (
	"flag"
	"fmt"

	// "fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

var (
	MesheryServerHost = "host.docker.internal:9081"
	TestingServer     = "localhost:9081"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a coma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

type Proxy struct {
	token      string
	authWsChan chan bool
}

func handleWsMessage(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr, " ", req.Method, " ", req.URL)
	if p.authWsChan == nil {
		p.authWsChan = make(chan bool)
	}

	client := &http.Client{}

	wr.Header().Set("Access-Control-Allow-Origin", "*")

	switch req.URL.Path {
	case "/ws":
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(wr, req, nil)
		if err != nil {
			log.Println(err)
		}
		go handleWsMessage(ws)

		for res := range p.authWsChan {
			if res == true {
				err = ws.WriteMessage(1, []byte("Authenticated"))
				if err != nil {
					log.Println(err)
				}
			}
		}

	case "/token/store":
		if req.Method == "GET" {
			values := req.URL.Query()
			var token string
			if values["token"] != nil {
				token = values["token"][0]
			}
			if token != "" {
				p.token = token
				p.authWsChan <- true
				log.Println("Setting the value of token to be: ", token)
			}
			if p.token == "" {
				wr.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				htmlTemplate := `
        <html>
<head>
<title>Page Title</title>
</head>
<body>
<script type="text/javascript">
  window.open('docker-desktop://dashboard/open','_self')
</script>
<p>You have been authenticated succesfully, you can close this window now.</p>

</body>
</html>
        `
				fmt.Fprint(wr, htmlTemplate)
			}
		}
	case "/token":
		if req.Method == "GET" {
			if p.token != "" {
				fmt.Fprintf(wr, p.token)
			}
		}
	default:
		//http: Request.RequestURI can't be set in client requests.
		//http://golang.org/src/pkg/net/http/client.go
		req.RequestURI = ""
		req.URL.Scheme = "http"

		delHopHeaders(req.Header)

		req.URL.Host = MesheryServerHost

		if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			appendHostToXForwardHeader(req.Header, clientIP)
		}

		if p.token == "" {
			log.Println("Sending request without auth token")
		} else {
			log.Println("Setting token cookie to the request")
			req.AddCookie(&http.Cookie{
				Name:     "token",
				Value:    string(p.token),
				Path:     "/",
				HttpOnly: true,
			})
			req.AddCookie(&http.Cookie{
				Name:     "meshery-provider",
				Value:    "Meshery",
				Path:     "/",
				HttpOnly: true,
			})
			req.AddCookie(&http.Cookie{
				Name:     "meshery.layer5.io_ref",
				Value:    "/",
				Path:     "/",
				HttpOnly: true,
			})
		}
		log.Println(*req)
		resp, err := client.Do(req)
		if err != nil {
			http.Error(wr, "Server Error", http.StatusInternalServerError)
			log.Fatal("ServeHTTP:", err)
		}
		defer resp.Body.Close()

		log.Println(req.RemoteAddr, " ", resp.Status)

		delHopHeaders(resp.Header)

		copyHeader(wr.Header(), resp.Header)
		wr.WriteHeader(resp.StatusCode)
		io.Copy(wr, resp.Body)
	}
}

func main() {
	var addr = flag.String("addr", "127.0.0.1:8080", "The addr of the application.")
	flag.Parse()

	handler := &Proxy{}

	log.Println("Starting proxy server on", *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
