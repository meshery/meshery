package handlers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func loadTestHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
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

	q.Set("url", os.Getenv("PRODUCT_PAGE_URL"))

	q.Set("json", "on")

	client := http.DefaultClient
	fortioURL, err := url.Parse(os.Getenv("FORTIO_URL"))
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
	w.Write(bd)
}
