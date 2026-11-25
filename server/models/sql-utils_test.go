package models

import "testing"

// The `order` query parameter is camelCase on the wire (identifier-naming
// guide: "?order=updatedAt desc") while the ORDER BY runs against snake_case
// database columns. SanitizeOrderInput is the translation point, and it is also
// the injection guard for a clause that gets interpolated into SQL, so both
// properties are covered here.
func TestSanitizeOrderInput(t *testing.T) {
	validColumns := []string{"created_at", "updated_at", "name"}

	tests := []struct {
		name     string
		order    string
		columns  []string
		expected string
	}{
		{
			name:     "canonical camelCase wire key folds onto the db column",
			order:    "createdAt desc",
			columns:  validColumns,
			expected: "created_at desc",
		},
		{
			name:     "legacy snake_case input still accepted",
			order:    "created_at desc",
			columns:  validColumns,
			expected: "created_at desc",
		},
		{
			name:     "single-word column is unaffected by folding",
			order:    "name asc",
			columns:  validColumns,
			expected: "name asc",
		},
		{
			name:     "direction is case-insensitive",
			order:    "createdAt DESC",
			columns:  validColumns,
			expected: "created_at desc",
		},
		{
			name:     "unknown direction falls back to asc",
			order:    "createdAt sideways",
			columns:  validColumns,
			expected: "created_at asc",
		},
		{
			name:     "column outside the whitelist is ignored",
			order:    "password desc",
			columns:  validColumns,
			expected: "",
		},
		{
			name:     "more than two words is rejected",
			order:    "created_at desc; DROP TABLE users",
			columns:  validColumns,
			expected: "",
		},
		{
			name:     "sql injection in the column position is rejected",
			order:    "created_at;DROP desc",
			columns:  validColumns,
			expected: "",
		},
		{
			name:     "empty order is rejected",
			order:    "",
			columns:  validColumns,
			expected: "",
		},
		{
			name:     "single word is rejected",
			order:    "created_at",
			columns:  validColumns,
			expected: "",
		},
		{
			name:     "extra spaces",
			order:    "createdAt   desc",
			columns:  validColumns,
			expected: "",
		},
		{
			name:     "tab separator",
			order:    "createdAt\tdesc",
			columns:  validColumns,
			expected: "",
		},
		{
			name:     "uppercase column",
			order:    "CREATED_AT desc",
			columns:  validColumns,
			expected: "",
		},
		{
			name:     "SQL comment in column",
			order:    "created_at-- desc",
			columns:  validColumns,
			expected: "",
		},
		{
			name:     "quoted column",
			order:    "'name' desc",
			columns:  validColumns,
			expected: "",
		},
		{
			name:     "empty columns slice",
			order:    "createdAt desc",
			columns:  []string{},
			expected: "",
		},
		{
			name:     "nil columns slice",
			order:    "createdAt desc",
			columns:  nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeOrderInput(tt.order, tt.columns); got != tt.expected {
				t.Errorf("SanitizeOrderInput(%q) = %q, want %q", tt.order, got, tt.expected)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"createdAt", "created_at"},
		{"userID", "user_id"},
		{"created_at", "created_at"},
		{"name", "name"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := toSnakeCase(tt.input); got != tt.expected {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
