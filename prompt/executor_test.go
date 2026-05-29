package prompt

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/c-bata/go-prompt"
)

func TestFetchMetadataSuggestions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock table names query
	tableRows := sqlmock.NewRows([]string{"TABLE_NAME"}).
		AddRow("USERS").
		AddRow("ORDERS")
	mock.ExpectQuery("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES").WillReturnRows(tableRows)

	// Mock column names query
	columnRows := sqlmock.NewRows([]string{"COLUMN_NAME"}).
		AddRow("ID").
		AddRow("NAME")
	mock.ExpectQuery("SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS").WillReturnRows(columnRows)

	suggestions := fetchMetadataSuggestions(db)

	if len(suggestions) != 4 {
		t.Errorf("Expected 4 suggestions, got %d", len(suggestions))
	}

	expected := map[string]string{
		"USERS":  "Table Name",
		"ORDERS": "Table Name",
		"ID":     "Column Name",
		"NAME":   "Column Name",
	}

	for _, s := range suggestions {
		desc, exists := expected[s.Text]
		if !exists {
			t.Errorf("Unexpected suggestion: %s", s.Text)
		} else if desc != s.Description {
			t.Errorf("For %s, expected description %q, got %q", s.Text, desc, s.Description)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPromptSessionCompleter(t *testing.T) {
	session := &PromptSession{
		suggestions: []prompt.Suggest{
			{Text: "SELECT", Description: "SQL Keyword"},
			{Text: "USERS", Description: "Table Name"},
			{Text: "USER_ID", Description: "Column Name"},
		},
	}

	tests := []struct {
		input string
		want  []string
	}{
		{
			input: "SEL",
			want:  []string{"SELECT"},
		},
		{
			input: "US",
			want:  []string{"USERS", "USER_ID"},
		},
		{
			input: "XYZ",
			want:  []string{},
		},
		{
			input: "",
			want:  []string{},
		},
	}

	for _, tt := range tests {
		doc := prompt.Document{
			Text: tt.input,
		}
		// Set unexported cursorPosition field using reflection and unsafe
		field := reflect.ValueOf(&doc).Elem().FieldByName("cursorPosition")
		ptr := unsafe.Pointer(field.UnsafeAddr())
		*(*int)(ptr) = len(tt.input)

		got := session.completer(doc)
		if len(got) != len(tt.want) {
			t.Errorf("For input %q, expected %d suggestions, got %d", tt.input, len(tt.want), len(got))
			continue
		}
		for i, s := range got {
			if s.Text != tt.want[i] {
				t.Errorf("For input %q at index %d, expected %q, got %q", tt.input, i, tt.want[i], s.Text)
			}
		}
	}
}
