package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"net/url"

	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

const saasTokenName = "meshery_saas"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/play/dashboard", http.StatusPermanentRedirect)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tu := "http://" + r.Host + r.RequestURI
	username := r.URL.Query().Get(saasTokenName)

	sess, err := sessionStore.Get(r, sessionName)
	if err == nil {
		sess.Config.MaxAge = -1
		sess.Save(w)
	}

	if username == "" {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    "/play/dashboard",
			Expires:  time.Now().Add(5 * time.Minute),
			Path:     "/play/",
			HttpOnly: true,
		})

		http.Redirect(w, r, os.Getenv("TWITTER_APP_HOST")+"?source="+base64.URLEncoding.EncodeToString([]byte(tu)), http.StatusFound)
		return
	}
	issueSession(w, r)
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
			data := map[string]interface{}{
				"ByPassAuth": byPassAuth,
			}

			if !byPassAuth {
				session, err := sessionStore.Get(r, sessionName)
				if err != nil {
					logrus.Errorf("error getting session: %v", err)
					http.Error(w, "unable to get session", http.StatusUnauthorized)
					return
				}
				user, _ := session.Values["user"].(*models.User)
				data["User"] = user
			}

			err := getAOTokenTempl.Execute(w, data)
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

	token := ""
	for k, va := range req.URL.Query() {
		for _, v := range va {
			if k == saasTokenName {
				// logrus.Infof("setting user in session: %s", v)
				token = v
				break
			}
		}
	}
	if reffCk != nil && reffCk.Name != "" {
		reffCk.Expires = time.Now().Add(-2 * time.Second)
		http.SetCookie(w, reffCk)
	}
	session.Values[saasTokenName] = token
	user, err := getUserDetails(saasTokenName, token)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)

	}
	session.Values["user"] = user
	err = session.Save(w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
	}
	http.Redirect(w, req, reffURL, http.StatusFound)
}

func getUserDetails(tokenKey, tokenVal string) (*models.User, error) {
	saasURL, _ := url.Parse(os.Getenv("TWITTER_APP_HOST") + "/user")
	req, _ := http.NewRequest(http.MethodGet, saasURL.String(), nil)
	req.AddCookie(&http.Cookie{
		Name:     tokenKey,
		Value:    tokenVal,
		Path:     "/",
		HttpOnly: true,
		Domain:   saasURL.Hostname(),
	})
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		logrus.Errorf("unable to fetch user data: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read body: %v", err)
		return nil, err
	}
	u := &models.User{}
	err = json.Unmarshal(bd, u)
	if err != nil {
		logrus.Errorf("unable to unmarshal user: %v", err)
		return nil, err
	}
	logrus.Infof("retrieved user: %v", u)
	return u, nil
}

// logoutHandler destroys the session on POSTs and redirects to home.
func logoutHandler(w http.ResponseWriter, req *http.Request) {
	client := http.Client{}
	cReq, err := http.NewRequest(http.MethodGet, os.Getenv("TWITTER_APP_HOST")+"/logout", req.Body)
	if err != nil {
		logrus.Errorf("Error creating a client to logout from tweet app: %v", err)
		http.Error(w, "unable to logout at the moment", http.StatusInternalServerError)
		return
	}
	client.Do(cReq)
	// sessionStore.Destroy(w, sessionName)

	sess, err := sessionStore.Get(req, sessionName)
	if err == nil {
		sess.Config.MaxAge = -1
		sess.Save(w)
	}

	if byPassAuth {
		http.Redirect(w, req, req.Referer(), http.StatusFound)
	} else {
		http.Redirect(w, req, "/play/login", http.StatusFound)
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

func setupSession(userName string, w http.ResponseWriter) *models.User {
	sessionStore.Config.Path = "/play"
	session := sessionStore.New(sessionName)
	user := &models.User{
		UserId: userName,
	}
	session.Values["user"] = user
	err := session.Save(w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
	}
	return user
}
