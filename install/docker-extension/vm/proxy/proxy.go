package proxy

import (
	"errors"
	"flag"
	"fmt"

	// "fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
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

// redactURL returns u as a string with the token query parameter masked, so
// the auth token used by /token/store never ends up in the logs. A query that
// does not parse cleanly is dropped entirely: url.Query silently skips
// malformed pairs, so a token inside one would otherwise be logged as is.
func redactURL(u *url.URL) string {
	redacted := *u
	q, err := url.ParseQuery(u.RawQuery)
	switch {
	case err != nil:
		redacted.RawQuery = "REDACTED"
	case q.Has("token"):
		q.Set("token", "REDACTED")
		redacted.RawQuery = q.Encode()
	}
	return redacted.String()
}

// wsConn serializes writes to a websocket. Gorilla allows only one concurrent
// writer and panics otherwise, and three goroutines write to each connection:
// the handshake, the auth-state listener, and the echo loop.
type wsConn struct {
	*websocket.Conn
	writeMu sync.Mutex
}

func (c *wsConn) WriteMessage(messageType int, data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.Conn.WriteMessage(messageType, data)
}

// urlQuery matches the query part of a URL quoted in a string, such as an
// error message. Go quotes URLs in errors with %q, so the query runs until an
// unescaped double quote; escaped characters, apostrophes and spaces inside
// it are all consumed. Outside a quoted string it runs to the end, which only
// redacts more.
var urlQuery = regexp.MustCompile(`\?(?:\\.|[^"\\])*`)

// redactTokens masks the whole query of every URL quoted in s. It does not
// look for a "token" name, because a name can be percent-encoded (%74oken)
// and still be decoded by the receiver.
func redactTokens(s string) string {
	return urlQuery.ReplaceAllString(s, "?REDACTED")
}

func handleWsMessage(conn *wsConn) {
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
	log.Println(req.RemoteAddr, " ", req.Method, " ", redactURL(req.URL))

	client := &http.Client{}

	enableCors(&wr)

	switch req.URL.Path {
	case "/ws":
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		upgraded, err := upgrader.Upgrade(wr, req, nil)
		if err != nil {
			log.Println(err)
			return
		}
		ws := &wsConn{Conn: upgraded}

		// Each connection gets its own channel, and only the newest connection
		// receives auth notifications. Closing the previous channel stops that
		// connection's listener, so two connections never split notifications
		// between them.
		authWsChan := make(chan bool, 1)
		p.mu.Lock()
		if p.authWsChan != nil {
			close(p.authWsChan)
		}
		p.authWsChan = authWsChan
		// Queue the current state for the new connection. A login or logout
		// that happened after the UI last read /token but before this channel
		// was installed would otherwise never reach it.
		authWsChan <- p.token != ""
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
		// Send the handshake before starting the listener, so the client sees
		// it before the queued auth state.
		err = ws.WriteMessage(1, []byte(InitiationSuccessful))
		log.Println("Sent message to ws client: ", InitiationSuccessful)
		if err != nil {
			log.Println(err)
		}
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

	case "/token/store":
		if req.Method == http.MethodGet {
			values := req.URL.Query()
			var token string
			if values["token"] != nil {
				token = values["token"][0]
			}
			if token != "" {
				p.setToken(token)
				log.Println("Stored a new auth token")
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
		log.Println("proxy: forwarding", req.Method, redactURL(req.URL))
		resp, err := client.Do(req)
		if err != nil {
			http.Error(wr, "Server Error", http.StatusInternalServerError)
			// err is a *url.Error that embeds the full request URL, which can
			// carry a token query parameter, so log the redacted URL and only
			// the underlying cause.
			var uerr *url.Error
			if errors.As(err, &uerr) {
				err = uerr.Err
			}
			// The underlying error can still quote a URL, for example a
			// malformed redirect Location, so mask tokens in its text as well.
			log.Println("proxy: backend request failed:", req.Method, redactURL(req.URL), redactTokens(err.Error()))
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
