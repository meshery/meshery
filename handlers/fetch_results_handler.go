package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

func (h *Handler) FetchResultsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// ensuring session is intact before running load test
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("Error: unable to get session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}
	tokenVal, _ := session.Values[h.config.SaaSTokenName].(string)

	err = req.ParseForm()
	if err != nil {
		logrus.Errorf("Error: unable to parse form: %v", err)
		http.Error(w, "unable to process the received data", http.StatusForbidden)
		return
	}
	q := req.Form

	bdr, err := h.getResultsFromSaaS(h.config.SaaSTokenName, tokenVal, q.Get("page"), q.Get("pageSize"), q.Get("search"), q.Get("order"))
	if err != nil {
		http.Error(w, "error while getting load test results", http.StatusInternalServerError)
		return
	}
	w.Write(bdr)
}

func (h *Handler) getResultsFromSaaS(tokenKey, tokenVal, page, pageSize, search, order string) ([]byte, error) {
	logrus.Infof("attempting to fetch results from SaaS")
	saasURL, _ := url.Parse(h.config.SaaSBaseURL + "/results")
	q := saasURL.Query()
	if page != "" {
		q.Set("page", page)
	}
	if pageSize != "" {
		q.Set("page_size", pageSize)
	}
	if search != "" {
		q.Set("search", search)
	}
	if order != "" {
		q.Set("order", order)
	}
	saasURL.RawQuery = q.Encode()
	logrus.Debugf("constructed results url: %s", saasURL.String())
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
		logrus.Errorf("unable to get results: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		logrus.Infof("results successfully retrieved from SaaS")
		return bdr, nil
	}
	logrus.Errorf("error while fetching results: %s", bdr)
	return nil, fmt.Errorf("error while sending results - Status code: %d, Body: %s", resp.StatusCode, bdr)
}
