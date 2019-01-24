package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"github.com/layer5io/meshery/meshes"
	"github.com/sirupsen/logrus"
)

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

func dashboardHandler(ctx context.Context, meshClient meshes.MeshClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// page, err := ioutil.ReadFile("../public/get-ao-token.html")
			// if err != nil {
			// 	logrus.Errorf("error reading get-ao-token.html: %v", err)
			// 	http.Error(w, "unable to find the requested file", http.StatusNotFound)
			// 	return
			// }
			// fmt.Fprintf(w, string(page))
			err := getAOTokenTempl.Execute(w, nil)
			if err != nil {
				logrus.Errorf("error rendering the template for the page: %v", err)
				http.Error(w, "unable to serve the requested file", http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			aoDashRenderer(ctx, meshClient, w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
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

func validateAuth(req *http.Request) bool {
	_, err := sessionStore.Get(req, sessionName)
	if err == nil {
		return true
	}
	logrus.Errorf("session invalid, error: %v", err)
	return false
}
