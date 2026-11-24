package models

import (
	"encoding/json"
	stderrors "errors"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/httputil"
	"github.com/meshery/schemas/models/v1beta2/organization"
)

// CurrentOrgIDHeader is the organization-selection header Meshery Cloud uses
// to name the caller's acting organization on org-scoped APIs, including
// the Cloud credential endpoints (list, get by id, save, update, delete).
const CurrentOrgIDHeader = "X-Current-Orgid"

// OrganizationSelectionRequiredCode is the machine-readable reason in the
// v1beta2/core OrganizationSelectionRequiredError body (schemas v1.3.126).
// It is a contract value, not a MeshKit component-scoped code.
const OrganizationSelectionRequiredCode = "organization_selection_required"

const organizationSelectionRequiredMessage = "select an organization"

// OrganizationSelectionChoice is one organization the caller may name as
// X-Current-Orgid. The shape matches schemas v1beta2/core
// OrganizationSelectionChoice (id, name, description).
type OrganizationSelectionChoice struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// OrganizationSelectionRequiredError is the 409 body defined by
// schemas v1beta2/core OrganizationSelectionRequiredError. It is not a
// MeshKit error: MeshKit codes are unstable component-scoped identifiers.
type OrganizationSelectionRequiredError struct {
	Code    string                        `json:"code"`
	Message string                        `json:"message"`
	Choices []OrganizationSelectionChoice `json:"choices"`
}

// OrganizationSelectionRequired is the Go error carried for a 409
// organization_selection_required, including the choices the UI should show.
type OrganizationSelectionRequired struct {
	Code    string
	Message string
	Choices []OrganizationSelectionChoice
}

func (e *OrganizationSelectionRequired) Error() string {
	if e == nil || strings.TrimSpace(e.Message) == "" {
		return organizationSelectionRequiredMessage
	}
	return e.Message
}

// Body returns the contract JSON the client is shown.
func (e *OrganizationSelectionRequired) Body() OrganizationSelectionRequiredError {
	msg := organizationSelectionRequiredMessage
	var choices []OrganizationSelectionChoice
	if e != nil {
		if strings.TrimSpace(e.Message) != "" {
			msg = e.Message
		}
		choices = e.Choices
	}
	if choices == nil {
		choices = []OrganizationSelectionChoice{}
	}
	return OrganizationSelectionRequiredError{
		Code:    OrganizationSelectionRequiredCode,
		Message: msg,
		Choices: choices,
	}
}

// ErrOrganizationSelectionRequired returns a 409 asking the UI to have the
// user pick one of choices as X-Current-Orgid.
func ErrOrganizationSelectionRequired(choices []OrganizationSelectionChoice) error {
	return httputil.WithProviderStatus(
		&OrganizationSelectionRequired{
			Code:    OrganizationSelectionRequiredCode,
			Message: organizationSelectionRequiredMessage,
			Choices: choices,
		},
		http.StatusConflict,
	)
}

// AsOrganizationSelectionRequired unwraps a 409 organization-selection error,
// including one joined with a provider status.
func AsOrganizationSelectionRequired(err error) (*OrganizationSelectionRequired, bool) {
	var sel *OrganizationSelectionRequired
	if stderrors.As(err, &sel) && sel != nil {
		return sel, true
	}
	return nil, false
}

// parseOrganizationSelectionRequired reads a Cloud 409 body. The contract
// body is OrganizationSelectionRequiredError, not MeshKit.
func parseOrganizationSelectionRequired(body []byte) *OrganizationSelectionRequired {
	var payload OrganizationSelectionRequiredError
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil
	}
	if payload.Code != "" && payload.Code != OrganizationSelectionRequiredCode {
		return nil
	}
	if payload.Code == "" && len(payload.Choices) == 0 {
		return nil
	}
	msg := payload.Message
	if strings.TrimSpace(msg) == "" {
		msg = organizationSelectionRequiredMessage
	}
	return &OrganizationSelectionRequired{Code: payload.Code, Message: msg, Choices: payload.Choices}
}

// resolveActingOrganizationID returns the organization id that must be sent as
// X-Current-Orgid on a remote-provider call scoped to the caller's acting
// organization.
//
// The saved selection (Preference.SelectedOrganizationID) is the source of
// truth. When none is saved and the user belongs to exactly one organization,
// that organization is auto-selected and persisted through the same
// WriteToPersister + preference-sync path UserPrefsHandler uses. When several
// organizations qualify and none is selected, the credentials endpoint is not
// called; a 409 organization_selection_required with choices is returned so
// the UI can prompt. When no organization qualifies, a 403 is returned.
func (l *RemoteProvider) resolveActingOrganizationID(req *http.Request, userID, token string) (string, error) {
	if id := l.selectedOrganizationID(req, userID); id != "" {
		return id, nil
	}
	if userID == "" {
		// Token-only callers (for example K8sContextFromConnection) cannot
		// consult a saved selection. Do not list organizations: that would
		// 409 every multi-org user and persist nothing useful. Leave the
		// header unset and let Cloud apply its default (auto-select a sole
		// org, or 403 until those callers pass a user id).
		return "", nil
	}

	token, err := l.requestToken(req, token)
	if err != nil {
		return "", err
	}
	body, err := l.GetOrganizations(token, "0", "all", "", "", "")
	if err != nil {
		return "", err
	}

	orgs, total, err := parseActingOrganizations(body)
	if err != nil {
		return "", ErrDataRead(err, "Organization")
	}
	switch {
	case total == 1 && len(orgs) == 1:
		orgID := orgs[0].ID.String()
		l.persistActingOrganization(req, userID, orgID, token)
		return orgID, nil
	case total <= 0 || len(orgs) == 0:
		return "", ErrNoQualifyingOrganization()
	default:
		return "", ErrOrganizationSelectionRequired(organizationSelectionChoices(orgs))
	}
}

func (l *RemoteProvider) selectedOrganizationID(req *http.Request, userID string) string {
	if req != nil {
		if pref, _ := req.Context().Value(PerfObjCtxKey).(*Preference); pref != nil {
			if id := strings.TrimSpace(pref.SelectedOrganizationID); id != "" {
				return id
			}
		}
	}
	if l == nil || l.SessionPreferencePersister == nil || userID == "" {
		return ""
	}
	pref, err := l.ReadFromPersister(userID)
	if err != nil || pref == nil {
		return ""
	}
	return strings.TrimSpace(pref.SelectedOrganizationID)
}

func parseActingOrganizations(body []byte) ([]*organization.Organization, int, error) {
	var page OrganizationsPage
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, 0, err
	}
	var orgs []*organization.Organization
	for _, org := range page.Organizations {
		if org == nil || org.ID == uuid.Nil {
			continue
		}
		orgs = append(orgs, org)
	}
	total := page.TotalCount
	if total == 0 {
		total = len(orgs)
	}
	return orgs, total, nil
}

func organizationSelectionChoices(orgs []*organization.Organization) []OrganizationSelectionChoice {
	choices := make([]OrganizationSelectionChoice, 0, len(orgs))
	for _, org := range orgs {
		if org == nil {
			continue
		}
		choices = append(choices, OrganizationSelectionChoice{
			ID:          org.ID,
			Name:        org.Name,
			Description: org.Description,
		})
	}
	return choices
}

func (l *RemoteProvider) organizationChoicesFromRequest(req *http.Request) []OrganizationSelectionChoice {
	if l == nil || req == nil {
		return nil
	}
	token, err := l.GetToken(req)
	if err != nil {
		return nil
	}
	body, err := l.GetOrganizations(token, "0", "all", "", "", "")
	if err != nil {
		return nil
	}
	orgs, _, err := parseActingOrganizations(body)
	if err != nil {
		return nil
	}
	return organizationSelectionChoices(orgs)
}

// persistActingOrganization records orgID as the user's selected organization
// locally and enqueues a remote preference sync when that path is available.
// Failure to persist is logged and does not block the in-flight call, which
// already has the organization id to send as X-Current-Orgid.
func (l *RemoteProvider) persistActingOrganization(req *http.Request, userID, orgID, token string) {
	if l == nil || orgID == "" {
		return
	}

	var pref *Preference
	if req != nil {
		if existing, ok := req.Context().Value(PerfObjCtxKey).(*Preference); ok && existing != nil {
			pref = existing
		}
	}
	if pref == nil && l.SessionPreferencePersister != nil && userID != "" {
		stored, err := l.ReadFromPersister(userID)
		if err != nil {
			// A read failure must not persist NewDefaultPreference, which
			// would overwrite the user's saved preferences locally and on
			// the remote provider. The in-flight call already has orgID.
			if l.Log != nil {
				l.Log.Error(err)
			}
			return
		} else if stored != nil {
			pref = stored
		}
	}
	if pref == nil {
		pref = NewDefaultPreference()
	}
	pref.SelectedOrganizationID = orgID
	if l.Log != nil {
		l.Log.Info("auto-selected sole organization ", orgID)
	}

	if userID == "" || l.SessionPreferencePersister == nil {
		return
	}
	if err := l.WriteToPersister(userID, pref); err != nil {
		if l.Log != nil {
			l.Log.Error(err)
		}
		return
	}
	if l.syncChan == nil {
		return
	}
	token, err := l.requestToken(req, token)
	if err != nil {
		if l.Log != nil {
			l.Log.Error(err)
		}
		return
	}
	l.syncChan <- &userSession{
		token:   token,
		session: pref,
	}
}

func requestUserID(req *http.Request) string {
	if req == nil {
		return ""
	}
	user, _ := req.Context().Value(UserCtxKey).(*User)
	if user == nil || user.ID == uuid.Nil {
		return ""
	}
	return user.ID.String()
}

func credentialUserID(credential *Credential) string {
	if credential == nil || credential.UserId == uuid.Nil {
		return ""
	}
	return credential.UserId.String()
}

func (l *RemoteProvider) requestToken(req *http.Request, token string) (string, error) {
	token = strings.TrimSpace(token)
	if token != "" {
		return token, nil
	}
	if l == nil || req == nil {
		return "", ErrGetToken(stderrors.New("missing token"))
	}
	return l.GetToken(req)
}

// requestWithToken returns a request GetToken can read the token from. When
// req is already present it is cloned so the caller's cookie jar is not
// mutated.
func requestWithToken(req *http.Request, token string) *http.Request {
	token = strings.TrimSpace(token)
	if token == "" {
		return req
	}
	if req == nil {
		req, _ = http.NewRequest(http.MethodGet, "/", nil)
	} else {
		req = req.Clone(req.Context())
	}
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: token})
	return req
}

// setActingOrganizationHeader resolves the acting organization and sets
// X-Current-Orgid on the Cloud request. On a 409 or 403 from resolution it
// returns without calling the credentials API.
func (l *RemoteProvider) setActingOrganizationHeader(req *http.Request, userID, token string, cReq *http.Request) (string, error) {
	if userID == "" {
		userID = requestUserID(req)
	}
	orgID, err := l.resolveActingOrganizationID(req, userID, token)
	if err != nil {
		return "", err
	}
	if orgID != "" && cReq != nil {
		cReq.Header.Set(CurrentOrgIDHeader, orgID)
	}
	return orgID, nil
}

// organizationScopePolicy controls how Cloud 409/403/404 responses are
// interpreted for a particular credential operation.
type organizationScopePolicy struct {
	// treatConflictAsSelection reports whether any 409 (including an
	// unparseable body) is organization_selection_required. That is true
	// for the credentials list, whose only 409 is the org-selection
	// contract, and false for save/update/delete/get-by-id, where a 409
	// may also be an ordinary resource conflict such as a duplicate name.
	treatConflictAsSelection bool
	// treatForbiddenAsNoOrg reports whether a plain 403 (not the
	// selection contract) is "no qualifying organization". That is true
	// for list and save. On get/update/delete an org has already been
	// sent, so a plain 403 is left as a resource permission error.
	treatForbiddenAsNoOrg bool
	// treatNotFoundAsMembership reports whether a 404 is always "not a
	// member of selected org". True for list and save, which have no
	// missing-credential 404. Get/update/delete only map 404 to
	// membership when the sent org is not one the caller belongs to.
	treatNotFoundAsMembership bool
}

var (
	credentialCollectionScope = organizationScopePolicy{
		treatConflictAsSelection:  true,
		treatForbiddenAsNoOrg:     true,
		treatNotFoundAsMembership: true,
	}
	credentialResourceScope = organizationScopePolicy{
		treatConflictAsSelection:  false,
		treatForbiddenAsNoOrg:     false,
		treatNotFoundAsMembership: false,
	}
	// Save can 409 on a duplicate name, so only a contract 409 is
	// treated as organization_selection_required. A 403/404 is still
	// organization-scope: there is no missing-credential 404 on POST.
	credentialSaveScope = organizationScopePolicy{
		treatConflictAsSelection:  false,
		treatForbiddenAsNoOrg:     true,
		treatNotFoundAsMembership: true,
	}
)

// organizationScopeError maps a Cloud organization-scope response. 409 and a
// contract 403 become organization_selection_required. A plain 403 is "no
// qualifying organization". A 404 is "not a member of selected org" when
// treatNotFoundAsMembership is true, or when the org we sent is no longer
// one the caller belongs to. It returns nil when the status is not an
// organization-scope failure (including a resource 404 on get/update/delete).
func (l *RemoteProvider) organizationScopeError(req *http.Request, status int, body []byte, orgID string, policy organizationScopePolicy) error {
	switch status {
	case http.StatusConflict:
		sel := parseOrganizationSelectionRequired(body)
		if sel != nil && sel.Code == OrganizationSelectionRequiredCode {
			if len(sel.Choices) < 2 {
				if choices := l.organizationChoicesFromRequest(req); len(choices) >= 2 {
					sel.Choices = choices
				}
			}
			return httputil.WithProviderStatus(sel, http.StatusConflict)
		}
		if !policy.treatConflictAsSelection {
			return nil
		}
		if sel == nil {
			sel = &OrganizationSelectionRequired{
				Code:    OrganizationSelectionRequiredCode,
				Message: organizationSelectionRequiredMessage,
			}
		}
		if len(sel.Choices) < 2 {
			if choices := l.organizationChoicesFromRequest(req); len(choices) >= 2 {
				sel.Choices = choices
			}
		}
		return httputil.WithProviderStatus(sel, http.StatusConflict)
	case http.StatusForbidden:
		// Live Cloud v1.0.384 may still answer 403 for a multi-org caller
		// with no named org. Promote to 409 only when the body is already
		// the selection-required contract (the code and at least two
		// choices). Otherwise this is "no qualifying organization" only
		// when the operation is collection-scoped (list/save).
		if sel := parseOrganizationSelectionRequired(body); sel != nil && sel.Code == OrganizationSelectionRequiredCode && len(sel.Choices) >= 2 {
			return httputil.WithProviderStatus(sel, http.StatusConflict)
		}
		if !policy.treatForbiddenAsNoOrg {
			return nil
		}
		err := ErrNoQualifyingOrganization()
		if l != nil && l.Log != nil {
			l.Log.Error(err)
		}
		return err
	case http.StatusNotFound:
		if policy.treatNotFoundAsMembership || l.actingOrganizationIsNotMember(req, orgID) {
			err := ErrOrganizationMembership()
			if l != nil && l.Log != nil {
				l.Log.Error(err)
			}
			return err
		}
		return nil
	default:
		return nil
	}
}

func (l *RemoteProvider) actingOrganizationIsNotMember(req *http.Request, orgID string) bool {
	orgID = strings.TrimSpace(orgID)
	if l == nil || orgID == "" {
		return false
	}
	choices := l.organizationChoicesFromRequest(req)
	if len(choices) == 0 {
		// Could not list organizations; do not claim a membership failure.
		return false
	}
	for _, choice := range choices {
		if choice.ID.String() == orgID {
			return false
		}
	}
	return true
}

func providerErrorStatus(err error, fallback int) int {
	if status, ok := httputil.ProviderStatusCode(err); ok {
		return status
	}
	return fallback
}
