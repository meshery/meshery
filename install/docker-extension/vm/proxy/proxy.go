package proxy

import (
	"flag"
	"fmt"

	// "fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	MesheryServerHost = "host.docker.internal:9081"
	TestingServer     = "localhost:9081"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const AuthenticatedMsg = "Authenticated"
const UnauthenticatedMsg = "Unauthenticated"
const InitiationSuccessful = "Initiated ws connection"

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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language,Content-Type, Access-Control-Request-Method")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
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
	// mu guards token and authWsChan, which are shared by concurrent requests.
	mu         sync.Mutex
	token      string
	authWsChan chan bool
}

func (p *Proxy) getToken() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.token
}

// setToken stores the token and tells the connected /ws client, if any, about
// the new auth state. It never blocks: without a client there is nobody to
// notify, and the UI reads the current state from /token when it connects.
func (p *Proxy) setToken(token string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.token = token
	if p.authWsChan == nil {
		return
	}
	// Replace any pending notification with the latest auth state, so a
	// logout is never lost behind a login notification still waiting to be
	// read. The channel is buffered (size 1) and setToken is its only
	// sender, so once drained the send below cannot block.
	select {
	case <-p.authWsChan:
	default:
	}
	p.authWsChan <- token != ""
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

	client := &http.Client{}

	enableCors(&wr)

	switch req.URL.Path {
	case "/ws":
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(wr, req, nil)
		if err != nil {
			log.Println(err)
			return
		}

		p.mu.Lock()
		if p.authWsChan == nil {
			p.authWsChan = make(chan bool, 1)
		}
		authWsChan := p.authWsChan
		p.mu.Unlock()

		ws.SetCloseHandler(func(code int, text string) error {
			p.mu.Lock()
			if p.authWsChan == authWsChan {
				close(p.authWsChan)
				p.authWsChan = nil
			}
			p.mu.Unlock()
			err := ws.Close()
			if err != nil {
				return err
			}
			return nil
		})
		go handleWsMessage(ws)
		go func() {
			for res := range authWsChan {
				log.Println("Got msg from authChan: ", res)
				if res == true {
					if werr := ws.WriteMessage(1, []byte(AuthenticatedMsg)); werr != nil {
						log.Println(werr)
					}
					log.Println("Sent message to ws client: ", AuthenticatedMsg)
				} else {
					if werr := ws.WriteMessage(1, []byte(UnauthenticatedMsg)); werr != nil {
						log.Println(werr)
					}
					log.Println("Sent message to ws client: ", UnauthenticatedMsg)
				}
			}
		}()
		err = ws.WriteMessage(1, []byte(InitiationSuccessful))
		log.Println("Sent message to ws client: ", InitiationSuccessful)
		if err != nil {
			log.Println(err)
		}

	case "/token/store":
		if req.Method == http.MethodGet {
			values := req.URL.Query()
			var token string
			if values["token"] != nil {
				token = values["token"][0]
			}
			if token != "" {
				p.setToken(token)
				log.Println("Setting the value of token to be: ", token)
			}
			if p.getToken() == "" {
				wr.WriteHeader(http.StatusUnauthorized)
				return
			} else {
				wr.Header().Set("Content-Type", "text/html; charset=utf-8")
				htmlTemplate := `<html>

<head>
  <title>Meshery | Docker Desktop</title>
</head>

<body>
  <script type="text/javascript">
    window.open('docker-desktop://dashboard/open', '_self')
  </script>
  <p>You have been authenticated successfully, you can safely close this window.</p>

</body>

</html>`
				fmt.Fprint(wr, htmlTemplate)
				// http.ServeFile(wr, req, "../assets/auth.html")
				return
			}
		}

	case "/token":
		if req.Method == http.MethodGet {
			if token := p.getToken(); token != "" {
				fmt.Fprint(wr, token)
			} else {
				fmt.Fprint(wr, "null")
			}
			return
		}
		if req.Method == http.MethodDelete {
			p.setToken("")
			log.Println("Deleted the existing token")
			wr.Write([]byte(""))
			return
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

		token := p.getToken()
		if token == "" {
			log.Println("Sending request without auth token")
		} else {
			log.Println("Setting token cookie to the request")
			req.AddCookie(&http.Cookie{
				Name:     "token",
				Value:    token,
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
				Name:     "cloud.meshery.io_ref",
				Value:    "/",
				Path:     "/",
				HttpOnly: true,
			})
		}
		log.Println(*req)
		resp, err := client.Do(req)
		if err != nil {
			http.Error(wr, "Server Error", http.StatusInternalServerError)
			log.Println("proxy: backend request failed:", err)
			return
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
