// main is an example web app using Login with Twitter.
package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"github.com/layer5io/meshery/appoptics"
	"github.com/layer5io/meshery/istio"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dghubble/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionUserKey       = "twitterID"
	sessionUserName      = "twitterUserName"
	sessionTwitterToken  = "token"
	sessionTwitterSecret = "secret"
	cookieSuffix         = "_referrer"
)

// sessionStore encodes and decodes session data stored in signed cookies
var (
	sessionName   = os.Getenv("EVENT") + "_pl"
	sessionSecret = base64.StdEncoding.EncodeToString([]byte(sessionName))
	sessionStore  = sessions.NewCookieStore([]byte(sessionSecret), nil)
	// cookieDuration = time.Now().Add(5 * time.Minute)
	cookieName = os.Getenv("EVENT") + cookieSuffix + "_pl"

	dashTempl = template.Must(template.ParseFiles("./dashboard.html"))
)

// Config configures the main ServeMux.
type Config struct {
	IstioClient *istio.IstioClient

	CallbackURL string
}

// New returns a new ServeMux with app routes.
func New(config *Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/play/static/", http.StripPrefix("/play/static/", http.FileServer(http.Dir("./static/"))))
	mux.Handle("/play/dashboard", authMiddleware(http.HandlerFunc(dashboardHandler)))
	mux.HandleFunc("/play/", indexHandler)

	mux.Handle("/play/load-test", authMiddleware(http.HandlerFunc(loadTestHandler)))
	mux.Handle("/play/istio", authMiddleware(http.HandlerFunc(istioHandler(config.IstioClient))))
	mux.HandleFunc("/play/logout", logoutHandler)
	mux.Handle("/play/tweet", authMiddleware(http.HandlerFunc(tweetHandler)))
	mux.HandleFunc("/play/login", loginHandler)

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/play/dashboard", http.StatusPermanentRedirect)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tu := "http://" + r.Host + r.RequestURI
	username := r.URL.Query().Get("username")
	if username == "" {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    "/play/dashboard",
			Expires:  time.Now().Add(5 * time.Minute),
			Path:     "/play/",
			HttpOnly: true,
		})

		http.Redirect(w, r, os.Getenv("TWITTER_APP_HOST")+"/twitter/login?source="+base64.URLEncoding.EncodeToString([]byte(tu)), http.StatusFound)
		return
	} else {
		issueSession(w, r)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		page, err := ioutil.ReadFile("./get-ao-token.html")
		if err != nil {
			logrus.Errorf("error reading get-ao-token.html: %v", err)
			http.Error(w, "unable to find the requested file", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, string(page))
	} else if r.Method == http.MethodPost {
		aoDashRenderer(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func aoDashRenderer(w http.ResponseWriter, req *http.Request) {
	token := req.FormValue("token")
	if token == "" {
		// http.Error(w, "Token not found", http.StatusNotFound)
		token = "ad6c84be90e7e16ef9150e0c0d809644956d5df6897b73d2340b3238fda40d9d"
	}
	spaceName := req.FormValue("dashboard")
	if spaceName == "" {
		spaceName = "istio"
	}
	logrus.Infof("retrieved token from query: %s", token)
	ao, err := appoptics.NewAOClient(token, spaceName)
	if err != nil {
		logrus.Errorf("error getting AO data: %v", err)
		http.Error(w, "unable to get data for the token", http.StatusNotFound)
		return
	}
	ad := ao.GenerateDataForTemplate()
	logrus.Infof("Retrieved AO data: %+#v", ad)
	err = dashTempl.Execute(w, ad)
	if err != nil {
		logrus.Errorf("error rendering the template for the dashboard: %v", err)
		http.Error(w, "unable to get data for the token", http.StatusInternalServerError)
		return
	}
}

// issueSession issues a cookie session after successful Twitter login
func issueSession(w http.ResponseWriter, req *http.Request) {
	var reffURL string
	reffCk, _ := req.Cookie(cookieName)
	if reffCk != nil {
		reffURL = reffCk.Value
	}
	logrus.Infof("preparing to issue session. retrieved reff url: %s", reffURL)
	if reffURL == "" {
		reffURL = "/play/"
	}
	sessionStore.Config.Path = "/play"
	session := sessionStore.New(sessionName)
	for k, va := range req.URL.Query() {
		for _, v := range va {
			if k == "username" {
				// logrus.Infof("setting user in session: %s", v)
				session.Values["user"] = v
			} else {
				// logrus.Infof("setting tweet cookie in session: %s", v)
				session.Values["tweet-cookie-name"] = k
				session.Values["tweet-cookie-value"] = v
			}
		}
	}
	if reffCk.Name != "" {
		reffCk.Expires = time.Now().Add(-2 * time.Second)
		// http.SetCookie(w, reffCk)
	}

	err := session.Save(w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
	}
	http.Redirect(w, req, reffURL, http.StatusFound)
}

func loadTestHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		logrus.Errorf("Error: unable to parse form: %v", err)
		http.Error(w, "unable to process the received data", http.StatusForbidden)
		return
	}
	q := req.Form
	q.Set("t", q.Get("t")+"m") // following fortio time indication
	q.Set("load", "Start")
	q.Set("runner", "http")

	cc, _ := strconv.Atoi(q.Get("c"))
	if cc < 1 || cc > 5 {
		q.Set("c", "5")
	}

	q.Set("url", os.Getenv("PRODUCT_PAGE_URL"))

	client := http.DefaultClient
	fortioURL, err := url.Parse(os.Getenv("FORTIO_URL"))
	if err != nil {
		logrus.Errorf("unable to parse the provided fortio url: %v", err)
		http.Error(w, "error while running load test", http.StatusInternalServerError)
		return
	}
	fortioURL.RawQuery = q.Encode()
	logrus.Infof("load test constructed url: %s", fortioURL.String())
	_, err = client.Get(fortioURL.String())
	if err != nil {
		logrus.Errorf("Error: unable to call fortio: %v", err)
		http.Error(w, "error while running load test", http.StatusInternalServerError)
		return
	}
}

func istioHandler(iClient *istio.IstioClient) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		query := req.PostFormValue("query")

		session, err := sessionStore.Get(req, sessionName)
		if err != nil {
			logrus.Error("unable to get session data")
			http.Error(w, "unable to get user data", http.StatusUnauthorized)
			return
		}
		userName, _ := session.Values["user"].(string)

		yamlTemplate := "virtual-service-all-v1"
		switch query {
		case "v1all":
			yamlTemplate = "virtual-service-all-v1"
		case "reviewsv2user":
			yamlTemplate = "virtual-service-reviews-test-v2"
		case "userTestDelay":
			yamlTemplate = "virtual-service-ratings-test-delay"
		case "userTestAbort":
			yamlTemplate = "virtual-service-ratings-test-abort"
		case "50v3":
			yamlTemplate = "virtual-service-reviews-50-v3"
		case "100v3":
			yamlTemplate = "virtual-service-reviews-v3"
		case "reset":
			istio.DeleteAllCreatedResources(iClient, "default")
			return
		}

		err = istio.ApplyYamlChange(iClient, fmt.Sprintf("templates/%s.yaml", yamlTemplate), userName, "default")
		if err != nil {
			logrus.Error(err)
			http.Error(w, "there was an error creating the services", http.StatusInternalServerError)
			return
		}
	}
}

// logoutHandler destroys the session on POSTs and redirects to home.
func logoutHandler(w http.ResponseWriter, req *http.Request) {
	client := http.DefaultClient
	cReq, err := http.NewRequest(http.MethodGet, os.Getenv("TWITTER_APP_HOST")+"/twitter/logout", req.Body)
	if err != nil {
		logrus.Errorf("Error creating a client to logout from tweet app: %v", err)
		http.Error(w, "unable to logout at the moment", http.StatusInternalServerError)
		return
	}
	sess, err := sessionStore.Get(req, sessionName)
	if err != nil {
		logrus.Errorf("Error retrieving tweet cookie: %v", err)
		http.Error(w, "unable to logout at the moment", http.StatusInternalServerError)
		return
	}
	ckNamei, ok := sess.Values["tweet-cookie-name"]
	if ok {
		ckVali, ok1 := sess.Values["tweet-cookie-value"]
		if ok1 {
			ckName, _ := ckNamei.(string)
			ckVal, _ := ckVali.(string)
			cReq.AddCookie(&http.Cookie{
				Name:  ckName,
				Value: ckVal,
			})
			client.Do(cReq)
		}
	}
	sessionStore.Destroy(w, sessionName)
	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func tweetHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		defer req.Body.Close()

		client := http.DefaultClient
		cReq, err := http.NewRequest(http.MethodPost, os.Getenv("TWITTER_APP_HOST")+"/twitter/tweet", req.Body)
		if err != nil {
			logrus.Errorf("Error creating a client to post a tweet: %v", err)
			http.Error(w, "unable to post a tweet at the moment", http.StatusInternalServerError)
			return
		}
		sess, err := sessionStore.Get(req, sessionName)
		if err != nil {
			logrus.Errorf("Error retrieving tweet cookie: %v", err)
			http.Error(w, "unable to logout at the moment", http.StatusInternalServerError)
			return
		}
		ckNamei, ok := sess.Values["tweet-cookie-name"]
		if ok {
			ckVali, ok1 := sess.Values["tweet-cookie-value"]
			if ok1 {
				ckName, _ := ckNamei.(string)
				ckVal, _ := ckVali.(string)
				cReq.AddCookie(&http.Cookie{
					Name:  ckName,
					Value: ckVal,
				})
				cReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				_, err = client.Do(cReq)
				if err != nil {
					logrus.Errorf("Error while posting the tweet: %v", err)
					http.Error(w, "unable to post a tweet at the moment", http.StatusInternalServerError)
					return
				}
			}
		}
	}
}

func authMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if !validateAuth(req) {
			http.Redirect(w, req, "/play/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func validateAuth(req *http.Request) bool {
	_, err := sessionStore.Get(req, sessionName)
	if err == nil {
		return true
	}
	logrus.Errorf("session invalid, error: %v", err)
	return false
}

func main() {
	twitterHost := os.Getenv("TWITTER_APP_HOST")
	if twitterHost == "" {
		logrus.Fatalf("TWITTER_APP_HOST environment variable not set.")
	}

	fortio := os.Getenv("FORTIO_URL")
	if fortio == "" {
		logrus.Fatalf("FORTIO_URL environment variable not set.")
	}

	productPageURL := os.Getenv("PRODUCT_PAGE_URL")
	if productPageURL == "" {
		logrus.Fatalf("PRODUCT_PAGE_URL environment variable not set.")
	}

	istioClient, err := istio.NewClient()
	if err != nil {
		logrus.Fatalf("Error creating an istio client: %v", err)
	}
	config := &Config{
		IstioClient: istioClient,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("Starting Server listening on %s", (":" + port))
	err = http.ListenAndServe(":"+port, New(config))
	if err != nil {
		logrus.Fatalf("ListenAndServe Error: %v", err)
	}
}
