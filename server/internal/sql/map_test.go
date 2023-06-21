package sql

import (
	"database/sql/driver"
	"testing"
)

// Test cases for sql.Map
func TestSMap(t *testing.T) {
	// Test sql.Map.Scan
	m := Map{}
	err := m.Scan([]byte(`{"foo": "bar"}`))
	if err != nil {
		t.Errorf("Scan error: %v", err)
	}
	if m["foo"] != "bar" {
		t.Errorf("Scan error: expected %v, got %v", "bar", m["foo"])
	}

	// Test sql.Map.Value
	v, err := m.Value()
	if err != nil {
		t.Errorf("Value error: %v", err)
	}

	if v != driver.Value(`{"foo":"bar"}`) {
		t.Errorf("Value error: expected %v, got %v", `{"foo":"bar"}`, string(v.([]byte)))
	}

	// Test sql.Map.UnmarshalJSON
	err = m.UnmarshalJSON([]byte(`{"foo": "bar"}`))
	if err != nil {
		t.Errorf("UnmarshalJSON error: %v", err)
	}
	if m["foo"] != "bar" {
		t.Errorf("UnmarshalJSON error: expected %v, got %v", "bar", m["foo"])
	}
}
