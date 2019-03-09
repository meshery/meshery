package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

func (h *Handler) LoadTestHandler(w http.ResponseWriter, req *http.Request) {
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

	tt, _ := strconv.Atoi(q.Get("t"))
	if tt < 1 {
		q.Set("t", "1m")
	} else {
		q.Set("t", q.Get("t")+"m") // following fortio time indication
	}

	q.Set("load", "Start")
	q.Set("runner", "http")

	cc, _ := strconv.Atoi(q.Get("c"))
	if cc < 1 {
		q.Set("c", "1")
	}

	loadTestURL := q.Get("url")
	ltURL, err := url.Parse(loadTestURL)
	if err != nil || !ltURL.IsAbs() {
		logrus.Errorf("unable to parse the provided load test url: %v", err)
		http.Error(w, "invalid load test URL", http.StatusBadRequest)
		return
	}

	qps, _ := strconv.Atoi(q.Get("qps"))
	if qps < 0 {
		q.Set("qps", "0")
	}

	q.Set("json", "on")

	client := &http.Client{}
	fortioURL, err := url.Parse(h.config.FortioURL)
	if err != nil {
		logrus.Errorf("unable to parse the provided fortio url: %v", err)
		http.Error(w, "error while running load test", http.StatusInternalServerError)
		return
	}
	fortioURL.RawQuery = q.Encode()
	logrus.Infof("load test constructed url: %s", fortioURL.String())
	fortioResp, err := client.Get(fortioURL.String())
	if err != nil {
		logrus.Errorf("Error: unable to call fortio: %v", err)
		http.Error(w, "error while running load test", http.StatusInternalServerError)
		return
	}
	defer fortioResp.Body.Close()
	bd, err := ioutil.ReadAll(fortioResp.Body)
	if err != nil {
		logrus.Errorf("Error: unable to parse response from fortio: %v", err)
		http.Error(w, "error while running load test", http.StatusInternalServerError)
		return
	}

	err = h.publishResultsToSaaS(h.config.SaaSTokenName, tokenVal, bd)
	if err != nil {
		// logrus.Errorf("Error: unable to parse response from fortio: %v", err)
		http.Error(w, "error while getting load test results", http.StatusInternalServerError)
		return
	}
	w.Write(bd)
}

func (h *Handler) publishResultsToSaaS(tokenKey, tokenVal string, bd []byte) error {
	logrus.Infof("attempting to publish results to SaaS")
	bf := bytes.NewBuffer(bd)
	saasURL, _ := url.Parse(h.config.SaaSBaseURL + "/result")
	req, _ := http.NewRequest(http.MethodPost, saasURL.String(), bf)
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
		logrus.Errorf("unable to send results: %v", err)
		return err
	}
	if resp.StatusCode == http.StatusCreated {
		logrus.Infof("results successfully pushlished to SaaS")
		return nil
	}
	defer resp.Body.Close()
	bdr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("unable to read response body: %v", err)
		return err
	}
	logrus.Errorf("error while sending results: %s", bdr)
	return fmt.Errorf("error while sending results - Status code: %d, Body: %s", resp.StatusCode, bdr)
}
