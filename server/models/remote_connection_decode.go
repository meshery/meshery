package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/schemas/models/core"
)

// errRemoteProviderConnectionMissing means the body parsed as JSON but was
// neither a ConnectionPage with an entry nor a bare Connection.
var errRemoteProviderConnectionMissing = errors.New("remote provider response did not include a connection")

// decodeRemoteProviderConnection reads a remote-provider register or delete
// body. registerConnection is moving from a ConnectionPage to a bare
// Connection, so both shapes are accepted:
//
//   - an object with a connections array is a ConnectionPage and yields its
//     first entry. A null first entry is returned as (nil, nil), which is
//     what SaveConnection did when the page contained [null].
//   - any other non-empty object is a bare Connection (a non-nil id, or any
//     other fields and no connections array).
//
// An empty page, JSON null, or an empty object returns
// errRemoteProviderConnectionMissing. An empty or invalid payload returns
// the encoding/json error so callers can keep treating it as a malformed body.
func decodeRemoteProviderConnection(body []byte) (*connections.Connection, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		var discard struct{}
		return nil, json.Unmarshal(trimmed, &discard)
	}

	var probe map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &probe); err != nil {
		return nil, err
	}
	if probe == nil {
		return nil, errRemoteProviderConnectionMissing
	}

	if _, isPage := probe["connections"]; isPage {
		var page connections.ConnectionPage
		if err := json.Unmarshal(trimmed, &page); err != nil {
			return nil, err
		}
		if len(page.Connections) > 0 {
			return page.Connections[0], nil
		}
		return nil, errRemoteProviderConnectionMissing
	}

	if len(probe) == 0 {
		return nil, errRemoteProviderConnectionMissing
	}

	var conn connections.Connection
	if err := json.Unmarshal(trimmed, &conn); err != nil {
		return nil, err
	}
	return &conn, nil
}

// connectionFromRemoteDelete turns a successful delete response into a
// connection callers can dereference. 204 and an empty body carry only the
// id that was deleted. A JSON body is a bare Connection or a ConnectionPage;
// when that body does not include an id, the deleted id is filled in.
func connectionFromRemoteDelete(body []byte, status int, connectionID core.Uuid) (*connections.Connection, error) {
	if status == http.StatusNoContent || len(bytes.TrimSpace(body)) == 0 {
		return &connections.Connection{ID: connectionID}, nil
	}

	decoded, err := decodeRemoteProviderConnection(body)
	if err != nil {
		if errors.Is(err, errRemoteProviderConnectionMissing) {
			return &connections.Connection{ID: connectionID}, nil
		}
		return nil, err
	}
	if decoded == nil {
		decoded = &connections.Connection{}
	}
	if decoded.ID == uuid.Nil {
		decoded.ID = connectionID
	}
	return decoded, nil
}
