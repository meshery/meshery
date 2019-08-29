package handlers

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"net/url"

	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

// func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "/play/dashboard", http.StatusPermanentRedirect)
// }

// LoginHandler redirects user for auth or issues session
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	tu := "http://" + r.Host + r.RequestURI
	token := r.URL.Query().Get(h.config.SaaSTokenName)

	sess, err := h.config.SessionStore.Get(r, h.config.SessionName)
	if err == nil {
		sess.Options.MaxAge = -1
		sess.Save(r, w)
	}

	if token == "" {
		http.SetCookie(w, &http.Cookie{
			Name:     h.config.RefCookieName,
			Value:    "/",
			Expires:  time.Now().Add(loginCookieDuration),
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, h.config.SaaSBaseURL+"?source="+base64.URLEncoding.EncodeToString([]byte(tu)), http.StatusFound)
		return
	}
	h.issueSession(w, r)
}

// issueSession issues a cookie session after successful Twitter login
func (h *Handler) issueSession(w http.ResponseWriter, req *http.Request) {
	var reffURL string
	reffCk, _ := req.Cookie(h.config.RefCookieName)
	if reffCk != nil {
		reffURL = reffCk.Value
	}
	logrus.Infof("preparing to issue session. retrieved reff url: %s", reffURL)
	if reffURL == "" {
		reffURL = "/"
	}
	// session, err := h.config.SessionStore.New(req, h.config.SessionName)
	session, _ := h.config.SessionStore.New(req, h.config.SessionName)
	// if err != nil {
	// 	logrus.Errorf("unable to create session: %v", err)
	// 	http.Error(w, "unable to create session", http.StatusInternalServerError)
	// 	return
	// }
	session.Options.Path = "/"
	token := ""
	for k, va := range req.URL.Query() {
		for _, v := range va {
			if k == h.config.SaaSTokenName {
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
	session.Values[h.config.SaaSTokenName] = token
	user, err := h.getUserDetails(token)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)

	}
	session.Values["user"] = user
	err = session.Save(req, w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
	}
	http.Redirect(w, req, reffURL, http.StatusFound)
}

func (h *Handler) getUserDetails(tokenVal string) (*models.User, error) {
	saasURL, _ := url.Parse(h.config.SaaSBaseURL + "/user")
	req, _ := http.NewRequest(http.MethodGet, saasURL.String(), nil)
	req.AddCookie(&http.Cookie{
		Name:     h.config.SaaSTokenName,
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

// LogoutHandler destroys the session on POSTs and redirects to home.
func (h *Handler) LogoutHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	client := http.Client{}
	cReq, err := http.NewRequest(http.MethodGet, h.config.SaaSBaseURL+"/logout", req.Body)
	if err != nil {
		logrus.Errorf("Error creating a client to logout from tweet app: %v", err)
		http.Error(w, "unable to logout at the moment", http.StatusInternalServerError)
		return
	}
	client.Do(cReq)
	// sessionStore.Destroy(w, sessionName)

	sess, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err == nil {
		sess.Options.MaxAge = -1
		sess.Save(req, w)
	}

	http.Redirect(w, req, "/login", http.StatusFound)
}

func (h *Handler) setupSession(userName string, req *http.Request, w http.ResponseWriter) *models.User {
	// sessionStore.Config.Path = "/play"
	session, err := h.config.SessionStore.New(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error creating a session: %v", err)
		http.Error(w, "unable to create a session", http.StatusInternalServerError)
		return nil
	}

	user := &models.User{
		UserID: userName,
	}
	session.Values["user"] = user
	err = session.Save(req, w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
	}
	return user
}
