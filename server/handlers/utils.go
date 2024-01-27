package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	defaultPageSize = 25
)

func getPaginationParams(req *http.Request) (page, offset, limit int, search, order, sortOnCol, status string) {

	urlValues := req.URL.Query()
	page, _ = strconv.Atoi(urlValues.Get("page"))
	limitstr := urlValues.Get("pagesize")
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 {
			limit = defaultPageSize
		}
	}

	search = urlValues.Get("search")
	order = urlValues.Get("order")
	sortOnCol = urlValues.Get("sort")
	status = urlValues.Get("status")

	if page <= 0 {
		page = 1
	}
	offset = (page - 1) * limit

	if sortOnCol == "" {
		sortOnCol = "updated_at"
	}
	return
}
func (handler *Handler) ViewHandler(responseWriter http.ResponseWriter, request *http.Request) {
	filePath, err := url.QueryUnescape(request.URL.Query().Get("file"))
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var jsonData interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonData); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(responseWriter)
	if err := encoder.Encode(jsonData); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DownloadEntity(responseWriter http.ResponseWriter, request *http.Request) {
	filePath, err := url.QueryUnescape(request.URL.Query().Get("file"))
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Header().Set("Content-Disposition", "attachment; filename=entity.json")

	_, err = io.Copy(responseWriter, file)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}
