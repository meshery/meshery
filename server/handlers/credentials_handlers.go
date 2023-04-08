package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(fmt.Errorf("error reading request body: %v", err))
		http.Error(w, "unable to read result data", http.StatusInternalServerError)
		return
	}

	userUUID := uuid.FromStringOrNil(user.ID)
	credential := models.Credential{
		UserID: &userUUID,
		Secret: map[string]interface{}{},
	}

	err = json.Unmarshal(bd, &credential)
	if err != nil {
		h.log.Error(fmt.Errorf("error unmarshal request body: %v", err))
		http.Error(w, "unable to parse credential data", http.StatusInternalServerError)
		return
	}

	result := provider.GetGenericPersister().Table("credentials").Create(&credential) // CredentialsDAO.SaveUserCredentials(&credential)
	if result.Error != nil {
		h.log.Error(fmt.Errorf("error saving user credentials: %v", result.Error))
		http.Error(w, "unable to save user credentials", http.StatusInternalServerError)
		return
	}

	h.log.Info("credential saved successfully")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ReadUserCredentials(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	order := q.Get("order")
	search := q.Get("search")
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	if pageSize > 25 {
		pageSize = 25
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if page < 0 {
		page = 0
	}
	if order == "" {
		order = "created_at desc"
	}

	h.log.Debug(fmt.Sprintf("page: %d, page size: %d, search: %s, order: %s", page+1, pageSize, search, order))

	result := provider.GetGenericPersister().Select("*").Where("user_id=? and deleted_at is NULL", user.UserID) // CredentialsDAO.GetUserCredentials(user.ID, page, pageSize, order, search)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		result = result.Where("(lower(name) like ?)", like)
	}

	result = result.Order(order)

	count, err := qc.Count(&models.Credential{})
	if err != nil {
		logrus.Errorf("error retrieving count of credentials for user id: %s - %v", userID, err)
		return nil, err
	}
	logrus.Debugf("retrieved total count: %d", count)

	credentialsList := []*models.Credential{}
	if count > 0 {
		if err := qc.Paginate(page+1, pageSize).All(&credentialsList); err != nil {
			if err.Error() != sql.ErrNoRows.Error() {
				logrus.Errorf("error retrieving credentials for user id: %s - %v", userID, err)
				return nil, err
			}
		}
	}
	logrus.Debugf("retrieved credentials: %+v", credentials)

	credentials := models.CredentialsPage{
		Credentials: credentialsList,
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  count,
	}

	h.log.Debug("credentials: ", credentials)
	if result.Error != nil {
		h.log.Error(fmt.Errorf("error getting user credentials: %v", result.Error))
		http.Error(w, "unable to get user credentials", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(credentials); err != nil {
		h.log.Error(fmt.Errorf("error encoding user credentials: %v", err))
		http.Error(w, "unable to encode user credentials", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateUserCredential(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
}

func (h *Handler) DeleteUserCredential(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
}
