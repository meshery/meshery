package sql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Map type is an alias for map[string]interface{}
//
// It implements native SQL driver interfaces and hence can
// be used for SQL json or jsonb types as a drop in replacement
// of golang native maps
type Map map[string]interface{}

// Interface implements the nulls.nullable interface.
func (m Map) Interface() interface{} {
	return map[string]interface{}(m)
}

// Scan implements the sql.Scanner interface.
// It allows to read the map from the database value.
func (m *Map) Scan(src interface{}) error {
	var b []byte

	switch t := src.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		return fmt.Errorf("scan source was not []byte nor string but %T", src)
	}

	return json.Unmarshal(b, m)
}

// Value implements the driver.Valuer interface.
// It allows to convert the map to a driver.value.
func (m Map) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// UnmarshalJSON will unmarshall JSON value into
// the map representation of this value.
func (m *Map) UnmarshalJSON(b []byte) error {
	var stuff map[string]interface{}

	if err := json.Unmarshal(b, &stuff); err != nil {
		return err
	}

	if *m == nil {
		*m = Map{}
	}

	for key, value := range stuff {
		(*m)[key] = value
	}

	return nil
}

// UnmarshalText will unmarshall text value into
// the map representation of this value.
func (m Map) UnmarshalText(text []byte) error {
	return json.Unmarshal(text, &m)
}
