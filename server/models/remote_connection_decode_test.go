package models

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid"
)

func TestDecodeRemoteProviderConnection(t *testing.T) {
	const (
		firstID   = "11111111-1111-1111-1111-111111111111"
		secondID  = "22222222-2222-2222-2222-222222222222"
		firstJSON = `{"id":"` + firstID + `","name":"prod","kind":"kubernetes","type":"platform","subType":"orchestrator","status":"discovered"}`
	)

	cases := []struct {
		name        string
		body        string
		wantMissing bool
		wantJSONErr bool
		wantNilConn bool
		wantID      string
		wantName    string
	}{
		{
			name:     "page uses the first connection",
			body:     `{"page":0,"page_size":2,"total_count":2,"connections":[` + firstJSON + `,{"id":"` + secondID + `","name":"other","kind":"kubernetes","status":"discovered"}]}`,
			wantID:   firstID,
			wantName: "prod",
		},
		{
			name:     "camelCase page metadata",
			body:     `{"page":0,"pageSize":1,"totalCount":1,"connections":[` + firstJSON + `]}`,
			wantID:   firstID,
			wantName: "prod",
		},
		{
			name:     "bare connection",
			body:     firstJSON,
			wantID:   firstID,
			wantName: "prod",
		},
		{
			name:     "bare connection without an id",
			body:     `{"name":"prod","kind":"kubernetes"}`,
			wantName: "prod",
		},
		{
			name:        "empty page",
			body:        `{"page":0,"page_size":0,"total_count":0,"connections":[]}`,
			wantMissing: true,
		},
		{
			name:        "null connections array",
			body:        `{"page":0,"page_size":0,"total_count":0,"connections":null}`,
			wantMissing: true,
		},
		{
			name:        "empty object",
			body:        `{}`,
			wantMissing: true,
		},
		{
			name:        "json null",
			body:        `null`,
			wantMissing: true,
		},
		{
			name:        "null page entry",
			body:        `{"connections":[null]}`,
			wantNilConn: true,
		},
		{
			name:        "empty body",
			body:        `   `,
			wantJSONErr: true,
		},
		{
			name:        "invalid json",
			body:        `{`,
			wantJSONErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn, err := decodeRemoteProviderConnection([]byte(tc.body))
			if tc.wantMissing {
				if !errors.Is(err, errRemoteProviderConnectionMissing) {
					t.Fatalf("err = %v, want missing connection", err)
				}
				if conn != nil {
					t.Fatalf("connection = %+v, want nil", conn)
				}
				return
			}
			if tc.wantJSONErr {
				if err == nil || errors.Is(err, errRemoteProviderConnectionMissing) {
					t.Fatalf("err = %v, want a JSON error", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if tc.wantNilConn {
				if conn != nil {
					t.Fatalf("connection = %+v, want nil", conn)
				}
				return
			}
			if conn == nil {
				t.Fatal("connection is nil")
			}
			if tc.wantID != "" && conn.ID.String() != tc.wantID {
				t.Errorf("id = %s, want %s", conn.ID, tc.wantID)
			}
			if tc.wantID == "" && conn.ID != uuid.Nil {
				t.Errorf("id = %s, want nil", conn.ID)
			}
			if conn.Name != tc.wantName {
				t.Errorf("name = %q, want %q", conn.Name, tc.wantName)
			}
		})
	}
}
