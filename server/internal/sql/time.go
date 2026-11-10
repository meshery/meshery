package sql

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Time is a wrapper for time.Time and implements
// Scanner and Valuer interfaces to support SQL driver
//
// This wrapper adds support for unconventional data types likes string
// []byte as date types but also supports time.Time parsing
type Time struct {
	Time time.Time
}

// Scan implements the SQL driver scanner interface
func (t *Time) Scan(v interface{}) error {
	var vt time.Time
	var err error
	layout := "2006-01-02 15:04:05"

	switch v := v.(type) {
	case string:
		vt, err = time.Parse(layout, v)
		if err != nil {
			return err
		}
	case []byte:
		vt, err = time.Parse(layout, string(v))
		if err != nil {
			return err
		}
	case time.Time:
		vt = v
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", v, t)
	}

	*t = Time{
		Time: vt,
	}

	return nil
}

// Value implements the SQL driver valuer interface
func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface. The time is expected to be
// a quoted string in RFC 3339 format
func (t *Time) UnmarshalJSON(b []byte) error {
	return t.Time.UnmarshalJSON(b)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in RFC 3339 format.
func (t *Time) UnmarshalText(text []byte) error {
	return t.Time.UnmarshalText(text)
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (t *Time) MarshalJSON() ([]byte, error) {
	return t.Time.MarshalJSON()
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in RFC 3339 format, with sub-second precision added if present.
func (t *Time) MarshalText() ([]byte, error) {
	return t.Time.MarshalText()
}
